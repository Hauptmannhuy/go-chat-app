package main

import (
	"encoding/json"
	"fmt"
	"log"
)

//go:generate jsonenums -type=Kind

type ActionOnType interface {
	perform(p []byte, msgT int, cl *Client)
}

type Kind int

const (
	NEW_MESSAGE Kind = iota
	NEW_CHAT
	SEARCH_QUERY
	JOIN_CHAT
	ERROR
)

var kindHandlers = map[Kind]func() ActionOnType{
	NEW_MESSAGE:  func() ActionOnType { return &UserMessage{} },
	NEW_CHAT:     func() ActionOnType { return &NewGroupChat{} },
	SEARCH_QUERY: func() ActionOnType { return &SearchQuery{} },
	JOIN_CHAT:    func() ActionOnType { return &JoinNotification{} },
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
	Body   string `json:"body"`
	ChatID string `json:"chat_id"`
	UserID string `json:"user_id"`
}

type JoinNotification struct {
	ChatID string `json:"chat_id"`
	UserID string `json:"user_id"`
}

type SearchQuery struct {
	Input         string `json:"input"`
	UserID        string `json:"user_id"`
	SearchResults interface{}
}

type NewGroupChat struct {
	ID       string `json:"chat_id"`
	UserID   string `json:"user_id"`
	ChatType string `json:"chat_type"`
}

type Error struct {
	Message string
}

func (um *UserMessage) perform(jsonEnv []byte, msgT int, cl *Client) {
	chatID := um.ChatID
	chat := chatList.Chats[chatID]
	for _, cl := range chat.members {
		sendWsResponse(jsonEnv, cl, msgT)
	}
}

func (ngc *NewGroupChat) perform(jsonEnv []byte, msgT int, cl *Client) {
	sendWsResponse(jsonEnv, cl, msgT)
	chatList.CreateChat(ngc.ID)
	chat := chatList.Chats[ngc.ID]
	chat.AddMember(cl)
}

func (sc *SearchQuery) perform(jsonEnv []byte, msgT int, cl *Client) {
	sendWsResponse(jsonEnv, cl, msgT)
}

func (jn *JoinNotification) perform(jsonEnv []byte, msgT int, cl *Client) {
	chat := chatList.Chats[jn.ChatID]
	chat.AddMember(cl)
}

func (e *Error) perform(jsonEnv []byte, msgT int, cl *Client) {
	sendWsResponse(jsonEnv, cl, msgT)
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
		outEnv.Data = ActionOnType(&errorMessage)
	}

	var action ActionOnType
	action = outEnv.Data.(ActionOnType)
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

func processEnvelope(p []byte) OutEnvelope {
	fmt.Println("raw json:", string(p))
	env := InEnvelope{}

	err := json.Unmarshal(p, &env)
	if err != nil {
		log.Fatal(err)
	}

	msg := kindHandlers[env.Type]()
	err = json.Unmarshal(p, msg)
	if err != nil {
		log.Fatal(err)
	}

	return OutEnvelope{
		Type: env.Type.toValue(),
		Data: msg,
	}
}
