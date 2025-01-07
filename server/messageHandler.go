package main

import (
	"encoding/json"
	"fmt"
	"go-chat-app/dbmanager/store"
	"log"
	"strings"
	"sync"

	"github.com/gorilla/websocket"
)

//go:generate jsonenums -type=Kind

type ActionOnType interface {
	perform(messageType string, wsMsgType int, cl *Client)
	saveDB() (interface{}, error)
	assignID(cl *Client)
}

type Cacheable interface {
	sendCache(cl *Client, messageType string, wsMessageT int)
}

type Kind int

const (
	NEW_MESSAGE Kind = iota
	NEW_CHAT
	SEARCH_QUERY
	NEW_PRIVATE_CHAT
	JOIN_CHAT
	LOAD_MESSAGES
	LOAD_SUBS
	NEW_GROUP_CHAT
)

func kindHandler(kind Kind) (interface{}, error) {
	var kindTypes = map[Kind]interface{}{
		NEW_MESSAGE:      &UserMessage{},
		NEW_CHAT:         &NewGroupChat{},
		SEARCH_QUERY:     &SearchQuery{},
		NEW_PRIVATE_CHAT: &NewPrivateChat{},
		JOIN_CHAT:        &Subscription{},
		LOAD_MESSAGES:    &WebSocketMessageStore{},
		LOAD_SUBS:        &WebSocketChatStore{},
		NEW_GROUP_CHAT:   &NewGroupChat{},
	}
	res, ok := kindTypes[kind]
	if !ok {
		log.Println("Unknown kind")
		return nil, fmt.Errorf("unknown kind")
	}
	return res, nil
}

func (um *UserMessage) assignID(cl *Client) {
	um.UserID = cl.index
	um.Username = cl.username

}

func (npc *NewPrivateChat) assignID(cl *Client) {
	npc.InitiatorID = cl.index
	npc.Username = cl.username
}

func (ngc *NewGroupChat) assignID(cl *Client) {
	ngc.CreatorID = cl.index
}

func (sq *SearchQuery) assignID(cl *Client) {
	sq.UserID = cl.index

}

func (sub *Subscription) assignID(cl *Client) {
	sub.UserID = cl.index
	sub.Username = cl.username
}
func (wsChatStore *WebSocketChatStore) assignID(cl *Client) {
	wsChatStore.UserID = cl.index
}

func (wsMessageStore *WebSocketMessageStore) assignID(cl *Client) {
	wsMessageStore.UserID = cl.index
}

func (k *Kind) toValue() string {
	var keys = map[Kind]string{
		0: "NEW_MESSAGE",
		1: "NEW_CHAT",
		2: "SEARCH_QUERY",
		3: "NEW_PRIVATE_CHAT",
		4: "JOIN_CHAT",
		5: "LOAD_MESSAGES",
		7: "NEW_GROUP_CHAT",
	}
	return keys[*k]
}

type OutEnvelope struct {
	Type string
	Data interface{}
}

type JSONenvelope struct {
	Type Kind
}

type UserMessage struct {
	Body      string `json:"body"`
	ChatName  string `json:"chat_name"`
	UserID    string `json:"user_id"`
	Username  string `json:"username"`
	State     string `json:"state"`
	MessageID int    `json:"message_id"`
}

type Subscription struct {
	ChatID      string `json:"chat_id"`
	UserID      string `json:"user_id"`
	BodyMessage string `json:"body_message"`
	Username    string `json:"username"`
	CreatorID   string `json:"creator_id"`
	ChatName    string `json:"chat_name"`
	msgID       int
}

type SearchQuery struct {
	Input         string `json:"input"`
	UserID        string `json:"user_id"`
	SearchResults interface{}
	Status        map[string]string `json:"status"`
}

type NewGroupChat struct {
	Name      string `json:"chat_name"`
	ID        string `json:"chat_id"`
	CreatorID string `json:"creator_id"`
}
type NewPrivateChat struct {
	ChatName    string `json:"chat_name"`
	ChatID      int    `json:"chat_id"`
	ReceiverID  string `json:"receiver_id"`
	InitiatorID string `json:"initiator_id"`
	Username    string `json:"init_username"`
	Message     string `json:"body"`
	MessageID   int    `json:"message_id"`
	ChatType    string `json:"chat_type"`
}

type OfflineMessages struct {
	Messages [][]UserMessage
}

type WebSocketMessageStore struct {
	UserID string
}

type WebSocketChatStore struct {
	UserID string
}

func (wsChatStore *WebSocketChatStore) sendCache(cl *Client, msgT string, wsMsgT int) {
	dbChatHandler := dbManager.initializeDBhandler("chat")
	groupChatData, err := dbChatHandler.LoadUserSubscribedChats(cl.index)

	if err != nil {
		log.Println(err)
		return
	}

	privateChatData, err := dbChatHandler.LoadSubscribedPrivateChats(cl.index)

	if err != nil {
		log.Println(err)
		return
	}

	chatContainer := make(map[string]interface{})
	chatContainer["private"] = privateChatData
	chatContainer["group"] = groupChatData

	writeToSocket(chatContainer, "LOAD_SUBS", cl, websocket.TextMessage)
}

func (wsMessageStore *WebSocketMessageStore) sendCache(client *Client, msgType string, wsMsgType int) {
	dbMessageHandler := dbManager.initializeDBhandler("message")
	data, err := dbMessageHandler.GetChatsMessages(client.subs)

	if err != nil {
		log.Println(err)
		writeToSocket(Error{Message: err.Error()}, "ERROR", client, wsMsgType)
	}

	writeToSocket(data, msgType, client, wsMsgType)
}

type Error struct {
	Message string
}

func (um *UserMessage) perform(msgType string, wsMsgType int, sender *Client) {
	chatID := um.ChatName
	chat := chatList.Chats[chatID]
	onlineUsers := chat.checkOnline()

	for _, key := range chat.subscribers {
		if userOnline := onlineUsers[key]; userOnline {
			userSocket := chat.members[key]
			writeToSocket(um, msgType, userSocket, wsMsgType)
		} else {
			fmt.Println("Saving message to Redis...")
			redisManager.saveMessage(*um)
		}
	}
}

func (ngc *NewGroupChat) perform(messageType string, wsMsgType int, cl *Client) {
	chatList.CreateChat(ngc.Name)
	chat := chatList.Chats[ngc.Name]
	chat.AddMember(cl)
	chat.AppendSubs([]string{cl.username})

	writeToSocket(ngc, messageType, cl, wsMsgType)
}

func (sc *SearchQuery) perform(messageType string, wsMsgType int, cl *Client) {

	writeToSocket(sc, messageType, cl, wsMsgType)
}

func (jn *Subscription) perform(messageType string, wsMsgType int, cl *Client) {
	chat := chatList.Chats[jn.ChatName]
	chat.AddMember(cl)
	chat.AppendSubs([]string{cl.username})

	var group sync.WaitGroup
	group.Add(1)

	go func() {
		defer group.Done()
		writeToSocket(jn, messageType, cl, wsMsgType)
	}()

	group.Wait()
	msg := UserMessage{
		Body:      jn.BodyMessage,
		UserID:    jn.UserID,
		ChatName:  jn.ChatID,
		Username:  jn.Username,
		MessageID: jn.msgID,
	}
	msg.perform("NEW_MESSAGE", wsMsgType, cl)
}

func (newPrCh *NewPrivateChat) perform(messageType string, wsMsgType int, cl *Client) {
	newChat := chatList.CreateChat(newPrCh.ChatName)
	newChat.AddMember(cl)

	split := strings.Split(newPrCh.ChatName, "_")
	receiverName := split[1]
	names := strings.Split(newPrCh.ChatName, "_")
	newChat.AppendSubs([]string{names[0], names[1]})

	receiverSocket, ok := connSockets.Connections[receiverName]

	var group sync.WaitGroup
	group.Add(1)

	go func() {
		defer group.Done()

		if ok {
			newChat.AddMember(receiverSocket)
			writeToSocket(newPrCh, messageType, receiverSocket, wsMsgType)
		} else {
			fmt.Println("Saving new private chat's message to Redis...", newPrCh)
			redisManager.saveMessage(*newPrCh)
		}

		writeToSocket(newPrCh, messageType, cl, wsMsgType)
	}()

	group.Wait()

	s := UserMessage{
		Body:      newPrCh.Message,
		UserID:    newPrCh.InitiatorID,
		ChatName:  newPrCh.ChatName,
		Username:  newPrCh.Username,
		MessageID: 0,
	}

	s.perform("NEW_MESSAGE", wsMsgType, cl)
}

func (m *UserMessage) saveDB() (interface{}, error) {
	messageHandler := dbManager.initializeDBhandler("message")
	messageID, err := messageHandler.CreateMessageHandler(m.Body, m.ChatName, m.UserID)
	m.MessageID = messageID

	return m, err
}

func (nc *NewGroupChat) saveDB() (interface{}, error) {
	chatHandler := dbManager.initializeDBhandler("chat")
	subHandler := dbManager.initializeDBhandler("subscription")
	id, err := chatHandler.CreateChatHandler(nc.Name, nc.CreatorID)
	nc.ID = id
	if err != nil {
		log.Println(err)
		return nil, err
	}
	err = subHandler.SaveSubHandler(nc.CreatorID, id)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return nc, err
}

func (sub *Subscription) saveDB() (interface{}, error) {
	subHandler := dbManager.initializeDBhandler("subscription")
	chatHandler := dbManager.initializeDBhandler("chat")
	messgHandler := dbManager.initializeDBhandler("message")
	msgID, err := messgHandler.CreateMessageHandler(sub.BodyMessage, sub.ChatID, sub.UserID)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	sub.msgID = msgID
	sub.CreatorID = chatHandler.RetrieveGroupChatCreatorID(sub.ChatID)
	subHandler.SaveSubHandler(sub.UserID, sub.ChatID)

	return sub, nil
}

func (sq *SearchQuery) saveDB() (interface{}, error) {
	chatHandler := dbManager.initializeDBhandler("chat")
	userHandler := dbManager.initializeDBhandler("user")
	result := make(map[string]interface{})
	sq.Status = make(map[string]string)
	resUsers, err := userHandler.SearchUser(sq.Input, sq.UserID)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	for key := range resUsers {

		ok := connSockets.isUserOnline(key)

		if ok {
			sq.Status[key] = "online"
		} else {
			sq.Status[key] = "offline"
		}
	}

	resChats, err := chatHandler.SearchChat(sq.Input, sq.UserID)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	result["users"] = resUsers
	result["chats"] = resChats
	sq.SearchResults = result

	return sq, nil
}

func (newPrCh *NewPrivateChat) saveDB() (interface{}, error) {
	chatHandler := dbManager.initializeDBhandler("chat")
	messageHandler := dbManager.initializeDBhandler("message")

	chatInfo, err := chatHandler.CreatePrivateChatHandler(newPrCh.InitiatorID, newPrCh.ReceiverID)
	if err != nil {
		log.Fatal(err)
	}

	data, _ := chatInfo.(store.PrivateChatInfo)

	chatID := data.ChatID
	chatName := data.ChatName

	messageData := UserMessage{
		Body:     newPrCh.Message,
		UserID:   newPrCh.InitiatorID,
		ChatName: chatName,
	}

	messageID, err := messageHandler.CreateMessageHandler(messageData.Body, messageData.ChatName, messageData.UserID)
	if err != nil {
		log.Fatal(err)
	}
	newPrCh.ChatName = chatName
	newPrCh.ChatID = chatID
	newPrCh.MessageID = messageID

	return newPrCh, nil
}

func dispatchAction(messageType string, data interface{}, wsMessageT int, cl *Client) {
	if errorMessage, ok := data.(Error); ok {
		errorMessage.handleError(data, cl, wsMessageT)
		return
	}
	action, ok := data.(ActionOnType)
	if ok {
		log.Println("performing action")
		action.perform(messageType, wsMessageT, cl)
		return
	} else {
		log.Println("No action to dispatch")
	}

	cacheAction, ok := data.(Cacheable)
	if ok {
		log.Println("sending cache")
		cacheAction.sendCache(cl, messageType, wsMessageT)
	} else {
		log.Println("No cacheable action")
	}

}

func processMessage(p []byte, cl *Client) (interface{}, string) {
	fmt.Println("raw json:", string(p))
	env := JSONenvelope{}
	err := json.Unmarshal(p, &env)

	if err != nil {
		ok := isTypeUnknown(err.Error())
		if ok {
			msg := Error{
				Message: err.Error(),
			}
			return msg, "UNKNOWN_TYPE"
		} else {
			log.Fatal(err)
		}
	}

	msg, err := kindHandler(env.Type)
	if err != nil {
		fmt.Println(err)
	}

	err = json.Unmarshal(p, msg)

	if err != nil {
		log.Fatal(err)
	}

	if assertedMsg, ok := msg.(ActionOnType); ok {
		assertedMsg.assignID(cl)
	}

	return msg, env.Type.toValue()
}

func (dbm *sqlDBwrap) handleDatabase(data interface{}) (interface{}, error) {
	if action, ok := data.(ActionOnType); ok {
		res, err := action.saveDB()
		return res, err
	} else {
		log.Println("No write to database")
		return data, nil
	}
}

func (e *Error) handleError(errorMessage interface{}, cl *Client, msgT int) {
	writeToSocket(errorMessage, "Error", cl, msgT)
}

func isTypeUnknown(err string) bool {
	decomposed := strings.Split(err, " ")
	if decomposed[0] == "invalid" && decomposed[1] == "Kind" {
		return true
	}
	return false
}
