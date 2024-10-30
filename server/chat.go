package main

import (
	"fmt"
	"sync"
)

type Chat struct {
	members []*Client
	ID      string `json:"chat_id"`
	mutex   sync.Mutex
}

type PrivateChat struct {
	members      []*Client
	mutex        sync.Mutex
	ID           string `json:"chat_id"`
	FirstUserID  string `json:"user_a_id"`
	SecondUserID string `json:"user_b_id"`
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
	fmt.Println(list)
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
