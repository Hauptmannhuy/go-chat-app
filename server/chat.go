package main

import "sync"

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
	if chL.Chats == nil {
		chL.Chats = make(map[string]*Chat)
	}
	chat := &Chat{
		ID: chID,
	}
	chL.Chats[chID] = chat
}

func (c *Chat) AddMember(cl *Client) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.members = append(c.members, cl)
}
