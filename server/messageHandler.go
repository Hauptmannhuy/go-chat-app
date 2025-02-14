package main

import (
	"fmt"
	"go-chat-app/dbmanager/store"
	"log"
	"strings"
)

//go:generate jsonenums -type=Kind

type MessageHandler interface {
	// execute(messageType string, wsMsgType int, cl *Client)
	execute() (interface{}, error)
	assignID(cl *Client)
}

type wsMessage struct {
	owner            *Client
	payload          interface{}
	broadcastHandler broadcastHandler
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
	UserID    int    `json:"user_id"`
	Username  string `json:"username"`
	State     string `json:"state"`
	MessageID int    `json:"message_id"`
}

type Subscription struct {
	ChatID      int    `json:"chat_id"`
	UserID      int    `json:"user_id"`
	BodyMessage string `json:"body_message"`
	Username    string `json:"username"`
	CreatorID   int    `json:"creator_id"`
	ChatName    string `json:"chat_name"`
	msgID       int
}

type SearchQuery struct {
	Input         string `json:"input"`
	UserID        int    `json:"user_id"`
	SearchResults interface{}
	Status        map[string]string `json:"status"`
}

type NewGroupChat struct {
	Name      string `json:"chat_name"`
	ID        int    `json:"chat_id"`
	CreatorID int    `json:"creator_id"`
}
type NewPrivateChat struct {
	ChatName    string `json:"chat_name"`
	ChatID      int    `json:"chat_id"`
	ReceiverID  int    `json:"receiver_id"`
	InitiatorID int    `json:"initiator_id"`
	Username    string `json:"init_username"`
	Message     string `json:"body"`
	MessageID   int    `json:"message_id"`
	ChatType    string `json:"chat_type"`
}

type OfflineMessages struct {
	Messages [][]UserMessage
}

type WebSocketMessageStore struct {
	UserID int
	Data   interface{}
}

type WebSocketChatStore struct {
	UserID int
	Data   interface{}
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

func defineAlgo(data interface{}) broadcastHandler {
	switch data.(type) {
	case *UserMessage:
		fmt.Println("return new message algo")
		return &newMsgAlgo{}
	case *NewGroupChat:
		fmt.Println("return new grch algo")
		return &newGroupChatAlgo{}
	case *NewPrivateChat:
		fmt.Println("return new prch algo")
		return &newDialogueAlgo{}
	case *SearchQuery:
		fmt.Println("return  sq algo")
		return &searchQueryAlgo{}
	case *Subscription:
		fmt.Println("return  sub algo")
		return &newSubAlgo{}
	case *WebSocketChatStore:
		fmt.Println("return  chatstore algo")
		return &chatStoreAlgo{}
	case *WebSocketMessageStore:
		fmt.Println("return  messagestore algo")
		return &messageStoreAlgo{}
	default:
		fmt.Println("return  error algo")
		return &errorAlgo{}
	}
}

func (um *UserMessage) assignID(cl *Client) {
	um.UserID = cl.id
	um.Username = cl.username
	fmt.Println("client username", cl.username)
}

func (npc *NewPrivateChat) assignID(cl *Client) {
	npc.InitiatorID = cl.id
	npc.Username = cl.username
}

func (ngc *NewGroupChat) assignID(cl *Client) {
	ngc.CreatorID = cl.id
}

func (sq *SearchQuery) assignID(cl *Client) {
	sq.UserID = cl.id

}

func (sub *Subscription) assignID(cl *Client) {
	sub.UserID = cl.id
	sub.Username = cl.username
}
func (wsChatStore *WebSocketChatStore) assignID(cl *Client) {
	wsChatStore.UserID = cl.id
}

func (wsMessageStore *WebSocketMessageStore) assignID(cl *Client) {
	wsMessageStore.UserID = cl.id
}

func (wsChatStore *WebSocketChatStore) execute() (interface{}, error) {
	dbChatHandler := dbManager.initializeDBhandler("chat")
	groupChatData, err := dbChatHandler.LoadUserSubscribedChats(wsChatStore.UserID)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	privateChatData, err := dbChatHandler.LoadSubscribedPrivateChats(wsChatStore.UserID)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	chatContainer := make(map[string]interface{})
	chatContainer["private"] = privateChatData
	chatContainer["group"] = groupChatData
	wsChatStore.Data = chatContainer

	return chatContainer, nil
}

func (wsMessageStore *WebSocketMessageStore) execute() (interface{}, error) {
	db := getDB()
	dbMessageHandler := db.initializeDBhandler("message")
	dbSubsHandler := db.initializeDBhandler("subscription")
	subs, err := dbSubsHandler.LoadSubscriptions(wsMessageStore.UserID)
	if err != nil {
		return nil, err
	}
	data, err := dbMessageHandler.GetChatsMessages(subs)
	wsMessageStore.Data = data
	return data, err
}

type Error struct {
	Message string
}

func (m *UserMessage) execute() (interface{}, error) {

	messageHandler := getDB().initializeDBhandler("message")
	messageID, err := messageHandler.CreateMessageHandler(m.Body, m.ChatName, m.UserID)
	m.MessageID = messageID
	fmt.Println(m)
	return m, err
}

func (nc *NewGroupChat) execute() (interface{}, error) {
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

func (sub *Subscription) execute() (interface{}, error) {
	subHandler := dbManager.initializeDBhandler("subscription")
	chatHandler := dbManager.initializeDBhandler("chat")
	messgHandler := dbManager.initializeDBhandler("message")
	msgID, err := messgHandler.CreateMessageHandler(sub.BodyMessage, sub.ChatName, sub.UserID)

	if err != nil {
		log.Println(err)
		return nil, err
	}

	sub.msgID = msgID
	sub.CreatorID = chatHandler.RetrieveGroupChatCreatorID(sub.ChatID)
	subHandler.SaveSubHandler(sub.UserID, sub.ChatID)

	return sub, nil
}

func (sq *SearchQuery) execute() (interface{}, error) {
	chatHandler := dbManager.initializeDBhandler("chat")
	userHandler := dbManager.initializeDBhandler("user")
	result := make(map[string]interface{})
	sq.Status = make(map[string]string)
	resUsers, err := userHandler.SearchUser(sq.Input, sq.UserID)

	if err != nil {
		log.Println(err)
		return nil, err
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

func (newPrCh *NewPrivateChat) execute() (interface{}, error) {
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

func isTypeUnknown(err string) bool {
	decomposed := strings.Split(err, " ")
	if decomposed[0] == "invalid" && decomposed[1] == "Kind" {
		return true
	}
	return false
}
