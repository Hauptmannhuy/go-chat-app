package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
)

//go:generate jsonenums -type=Kind

type ActionOnType interface {
	perform(p []byte, msgT int, cl *Client)
	save() (interface{}, error)
}

type Kind int

const (
	NEW_MESSAGE Kind = iota
	NEW_CHAT
	SEARCH_QUERY
	NEW_PRIVATE_CHAT
	JOIN_CHAT
)

var kindHandlers = map[Kind]func(cl *Client) ActionOnType{

	NEW_MESSAGE:      func(cl *Client) ActionOnType { return &UserMessage{UserID: cl.index} },
	NEW_CHAT:         func(cl *Client) ActionOnType { return &NewGroupChat{CreatorID: cl.index} },
	SEARCH_QUERY:     func(cl *Client) ActionOnType { return &SearchQuery{UserID: cl.index} },
	NEW_PRIVATE_CHAT: func(cl *Client) ActionOnType { return &NewPrivateChat{InitiatorID: cl.index} },
	JOIN_CHAT:        func(cl *Client) ActionOnType { return &Subscription{UserID: cl.index} },
}

func (k *Kind) toValue() string {
	var keys = map[Kind]string{
		0: "NEW_MESSAGE",
		1: "NEW_CHAT",
		2: "SEARCH_QUERY",
		3: "NEW_PRIVATE_CHAT",
		4: "JOIN_CHAT",
	}
	return keys[*k]
}

type OutEnvelope struct {
	Type string
	Data interface{}
}

type InEnvelope struct {
	Type Kind
}

type UserMessage struct {
	Body     string `json:"body"`
	ChatName string `json:"chat_name"`
	UserID   string `json:"user_id"`
}

type Subscription struct {
	ChatID string `json:"chat_id"`
	UserID string `json:"user_id"`
}

type SearchQuery struct {
	Input         string `json:"input"`
	UserID        string `json:"user_id"`
	SearchResults interface{}
}

type NewGroupChat struct {
	Name      string `json:"chat_name"`
	CreatorID string `json:"creator_id"`
}

type NewPrivateChat struct {
	ChatName    string `json:"chat_id"`
	ReceiverID  string `json:"receiver_id"`
	InitiatorID string `json:"initiator_id"`
	Username    string `json:"init_username"`
	Message     string `json:"message"`
	ChatType    string `json:"chat_type"`
}

type Error struct {
	Message string
}

func (um *UserMessage) perform(jsonEnv []byte, msgT int, cl *Client) {
	chatID := um.ChatName
	chat := chatList.Chats[chatID]
	for _, cl := range chat.members {
		sendWsResponse(jsonEnv, cl, msgT)
	}
}

func (ngc *NewGroupChat) perform(jsonEnv []byte, msgT int, cl *Client) {
	sendWsResponse(jsonEnv, cl, msgT)
	chatList.CreateChat(ngc.Name)
	chat := chatList.Chats[ngc.Name]
	chat.AddMember(cl)
}

func (sc *SearchQuery) perform(jsonEnv []byte, msgT int, cl *Client) {
	sendWsResponse(jsonEnv, cl, msgT)
}

func (jn *Subscription) perform(jsonEnv []byte, msgT int, cl *Client) {
	chat := chatList.Chats[jn.ChatID]
	chat.AddMember(cl)
}

func (newPrCh *NewPrivateChat) perform(jsonEnv []byte, msgT int, cl *Client) {
	receiverSocket, ok := connSockets.Connections[newPrCh.ReceiverID]
	if ok {
		sendWsResponse(jsonEnv, receiverSocket, msgT)
	}
	sendWsResponse(jsonEnv, cl, msgT)
}

func (m *UserMessage) save() (interface{}, error) {
	messageHandler := dbManager.initializeDBhandler("message")
	err := messageHandler.CreateMessageHandler(m.Body, m.ChatName, m.UserID)
	fmt.Println("Result ENV:", m)

	return m, err
}

func (nc *NewGroupChat) save() (interface{}, error) {
	chatHandler := dbManager.initializeDBhandler("chat")
	subHandler := dbManager.initializeDBhandler("subscription")
	chatName, err := chatHandler.CreateChatHandler(nc.Name, nc.CreatorID)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	err = subHandler.SaveSubHandler(nc.CreatorID, chatName)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return nc, err
}

func (sub *Subscription) save() (interface{}, error) {
	subHandler := dbManager.initializeDBhandler("subscription")
	subHandler.SaveSubHandler(sub.UserID, sub.ChatID)
	return sub, nil
}

func (sq *SearchQuery) save() (interface{}, error) {
	chatHandler := dbManager.initializeDBhandler("chat")
	userHandler := dbManager.initializeDBhandler("user")
	fmt.Println("SEARCH QUERY USER ID", sq.UserID)
	result := make(map[string]interface{})
	resUsers, err := userHandler.SearchUser(sq.Input, sq.UserID)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	resChats, err := chatHandler.SearchChat(sq.Input, sq.UserID)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	result["users"] = resUsers
	result["chats"] = resChats
	sq.SearchResults = result
	fmt.Println("Result ENV:", sq)

	return sq, nil
}

func (newPrCh *NewPrivateChat) save() (interface{}, error) {
	chatHandler := dbManager.initializeDBhandler("chat")
	messageHandler := dbManager.initializeDBhandler("message")

	chatID, err := chatHandler.CreatePrivateChatHandler(newPrCh.InitiatorID, newPrCh.ReceiverID)
	if err != nil {
		log.Fatal(err)
	}

	messageData := UserMessage{
		Body:     newPrCh.Message,
		UserID:   newPrCh.InitiatorID,
		ChatName: chatID,
	}

	err = messageHandler.CreateMessageHandler(messageData.Body, messageData.ChatName, messageData.UserID)
	if err != nil {
		log.Fatal(err)
	}
	newPrCh.ChatName = chatID
	return newPrCh, nil
}

func HandleWriteToWebSocket(outEnv OutEnvelope, msgT int, cl *Client) {
	jsonEnv, err := json.Marshal(outEnv)
	fmt.Println(string(jsonEnv))
	fmt.Println(cl.index, "client index in handle response env")
	fmt.Println("slice of sockets:", connSockets.Connections)
	if err != nil {
		fmt.Println(err)
		return
	}

	if errorMessage, ok := outEnv.Data.(Error); ok {
		errorMessage.handleError(jsonEnv, cl, msgT)
		return
	}

	action := outEnv.Data.(ActionOnType)
	action.perform(jsonEnv, msgT, cl)

}

func sendWsResponse(p []byte, cl *Client, msgT int) {
	socket := cl.socket
	if err := socket.WriteMessage(msgT, p); err != nil {
		log.Println("Error writing to WebSocket:", err)
		return
	}
	fmt.Println("Message sent successfully to client", string(p))
}

func processEnvelope(p []byte, cl *Client) OutEnvelope {
	fmt.Println("raw json:", string(p))
	env := InEnvelope{}
	err := json.Unmarshal(p, &env)

	if err != nil {
		ok := isTypeUnknown(err.Error())
		if ok {
			msg := Error{
				Message: err.Error(),
			}
			return OutEnvelope{
				Type: "UNKNOWN_TYPE",
				Data: msg,
			}
		} else {
			log.Fatal(err)
		}
	}

	msg := kindHandlers[env.Type](cl)

	err = json.Unmarshal(p, msg)
	if err != nil {
		log.Fatal(err)
	}

	return OutEnvelope{
		Type: env.Type.toValue(),
		Data: msg,
	}
}

func (dbm *sqlDBwrap) handleDatabase(env OutEnvelope) (interface{}, error) {
	if action, ok := env.Data.(ActionOnType); ok {
		res, err := action.save()
		return res, err
	} else {
		fmt.Println("No write to database")
		return env.Data, nil
	}
}

func (e *Error) handleError(p []byte, cl *Client, msgT int) {
	sendWsResponse(p, cl, msgT)
}

func isTypeUnknown(err string) bool {
	decomposed := strings.Split(err, " ")
	if decomposed[0] == "invalid" && decomposed[1] == "Kind" {
		return true
	}
	return false
}
