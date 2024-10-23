package main

import (
	"sync"
)

type Chat struct {
	members []*Client
	ID      string `json:"id"`
	mutex   sync.Mutex
}

type ChatList struct {
	Chats map[string]*Chat
	mutex sync.Mutex
}

var chatList ChatList

func (chL *ChatList) CreateChat(chID string) {
	chL.mutex.Lock()
	defer chL.mutex.Unlock()

	chat := &Chat{
		ID: chID,
	}
	chL.Chats[chID] = chat
}

func (chL *ChatList) initializeRooms() {
	dbChatHandler := dbManager.initializeDBhandler("chat")
	list, _ := dbChatHandler.GetAllChats()
	chL.Chats = make(map[string]*Chat)
	for _, chID := range list {
		chat := &Chat{
			ID: chID,
		}
		chL.Chats[chID] = chat
	}
}

func (chL *ChatList) addClientToSubRooms(cl *Client) {
	chL.mutex.Lock()
	defer chL.mutex.Unlock()
	for _, chID := range cl.subs {
		chat := chL.Chats[chID]
		chat.AddMember(cl)
	}
}

func (c *Chat) AddMember(cl *Client) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.members = append(c.members, cl)
}
