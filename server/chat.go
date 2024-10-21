package main

import "sync"

type Chat struct {
	Members []*Client
	ID      string
	Mutex   sync.Mutex
}

type ChatList struct {
	Chats map[string]*Chat
	Mutex sync.Mutex
}

var chatList ChatList

func (chL *ChatList) CreateChat(chID string) {
	chL.Mutex.Lock()
	defer chL.Mutex.Unlock()
	if chL.Chats == nil {
		chL.Chats = make(map[string]*Chat)
	}
	chat := &Chat{
		ID: chID,
	}
	chL.Chats[chID] = chat
}

func (c *Chat) AddMember(cl *Client) {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()
	c.Members = append(c.Members, cl)
}
