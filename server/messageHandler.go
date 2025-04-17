package main

import (
	"fmt"
	"go-chat-app/dbmanager/store"
	"log"
	"strings"
)

//go:generate jsonenums -type=Kind

type MessageHandler interface {
	Process(cl *Client)
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
	UserID    int    `json:"user_id"`
	Username  string `json:"username"`
	ChatName  string `json:"chat_name"`
	MessageID int    `json:"message_id"`
	Body      string `json:"body"`
	Image     Image  `json:"image"`
	State     string `json:"state"`
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

func kindHandler(kind Kind) (MessageHandler, error) {
	var kindTypes = map[Kind]MessageHandler{
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
		fmt.Println("return sub algo")
		return &newSubAlgo{}
	case *WebSocketChatStore:
		fmt.Println("return chatstore algo")
		return &chatStoreAlgo{}
	case *WebSocketMessageStore:
		fmt.Println("return  messagestore algo")
		return &messageStoreAlgo{}
	default:
		fmt.Println("return  error algo")
		return &errorAlgo{}
	}
}

func (msg *UserMessage) Process(cl *Client) {
	msg.UserID = cl.id
	msg.Username = cl.username
	msg.requestDB()
}

func (msg *NewPrivateChat) Process(cl *Client) {
	msg.InitiatorID = cl.id
	msg.Username = cl.username
	msg.requestDB()
}

func (msg *NewGroupChat) Process(cl *Client) {
	msg.CreatorID = cl.id
	msg.requestDB()
}

func (msg *SearchQuery) Process(cl *Client) {
	msg.UserID = cl.id
	msg.requestDB()
}

func (msg *Subscription) Process(cl *Client) {
	msg.UserID = cl.id
	msg.Username = cl.username
	msg.requestDB()
}
func (msg *WebSocketChatStore) Process(cl *Client) {
	msg.UserID = cl.id
	msg.requestDB()
}

func (msg *WebSocketMessageStore) Process(cl *Client) {
	msg.UserID = cl.id
	msg.requestDB()
}

func (msg *WebSocketChatStore) requestDB() error {
	dbChatHandler := dbManager.initializeDBhandler("chat")
	groupChatData, err := dbChatHandler.LoadUserSubscribedChats(msg.UserID)

	if err != nil {
		log.Println(err)
		return err
	}

	privateChatData, err := dbChatHandler.LoadSubscribedPrivateChats(msg.UserID)

	if err != nil {
		log.Println(err)
		return err
	}

	chatContainer := make(map[string]interface{})
	chatContainer["private"] = privateChatData
	chatContainer["group"] = groupChatData
	msg.Data = chatContainer

	return nil
}

func (msg *WebSocketMessageStore) requestDB() error {
	db := getDB()
	dbMessageHandler := db.initializeDBhandler("message")
	dbSubsHandler := db.initializeDBhandler("subscription")
	subs, err := dbSubsHandler.LoadSubscriptions(msg.UserID)
	if err != nil {
		return err
	}
	data, err := dbMessageHandler.GetChatsMessages(subs)
	msg.Data = data
	return nil
}

type Error struct {
	Message string
}

func (msg *UserMessage) requestDB() error {

	messageHandler := getDB().initializeDBhandler("message")
	messageID, err := messageHandler.CreateMessageHandler(msg.Body, msg.ChatName, msg.UserID)
	msg.MessageID = messageID
	fmt.Println(msg)
	return err
}

func (msg *NewGroupChat) requestDB() error {
	chatHandler := dbManager.initializeDBhandler("chat")
	subHandler := dbManager.initializeDBhandler("subscription")
	id, err := chatHandler.CreateChatHandler(msg.Name, msg.CreatorID)
	msg.ID = id
	if err != nil {
		log.Println(err)
		return err
	}
	err = subHandler.SaveSubHandler(msg.CreatorID, id)
	if err != nil {
		log.Println(err)
		return err
	}
	return err
}

func (msg *Subscription) requestDB() error {
	subHandler := dbManager.initializeDBhandler("subscription")
	chatHandler := dbManager.initializeDBhandler("chat")
	messgHandler := dbManager.initializeDBhandler("message")
	msgID, err := messgHandler.CreateMessageHandler(msg.BodyMessage, msg.ChatName, msg.UserID)

	if err != nil {
		log.Println(err)
		return err
	}

	msg.msgID = msgID
	msg.CreatorID = chatHandler.RetrieveGroupChatCreatorID(msg.ChatID)
	subHandler.SaveSubHandler(msg.UserID, msg.ChatID)

	return nil
}

func (msg *SearchQuery) requestDB() error {
	chatHandler := dbManager.initializeDBhandler("chat")
	userHandler := dbManager.initializeDBhandler("user")
	result := make(map[string]interface{})
	msg.Status = make(map[string]string)
	resUsers, err := userHandler.SearchUser(msg.Input, msg.UserID)

	if err != nil {
		log.Println(err)
		return err
	}

	resChats, err := chatHandler.SearchChat(msg.Input, msg.UserID)

	if err != nil {
		log.Println(err)
		return err
	}

	result["users"] = resUsers
	result["chats"] = resChats
	msg.SearchResults = result

	return nil
}

func (msg *NewPrivateChat) requestDB() error {
	chatHandler := dbManager.initializeDBhandler("chat")
	messageHandler := dbManager.initializeDBhandler("message")

	chatInfo, err := chatHandler.CreatePrivateChatHandler(msg.InitiatorID, msg.ReceiverID)
	if err != nil {
		log.Fatal(err)
	}

	data, _ := chatInfo.(store.PrivateChatInfo)

	chatID := data.ChatID
	chatName := data.ChatName

	messageData := UserMessage{
		Body:     msg.Message,
		UserID:   msg.InitiatorID,
		ChatName: chatName,
	}

	messageID, err := messageHandler.CreateMessageHandler(messageData.Body, messageData.ChatName, messageData.UserID)
	if err != nil {
		log.Fatal(err)
	}
	msg.ChatName = chatName
	msg.ChatID = chatID
	msg.MessageID = messageID

	return nil
}

func isTypeUnknown(err string) bool {
	decomposed := strings.Split(err, " ")
	if decomposed[0] == "invalid" && decomposed[1] == "Kind" {
		return true
	}
	return false
}
