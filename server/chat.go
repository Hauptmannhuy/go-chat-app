package main

import (
	"fmt"
	"sync"
)

type Chat struct {
	members     map[string]*Client
	subscribers []string
	Name        string `json:"chat_id"`
	mutex       sync.Mutex
}

type ChatList struct {
	Chats map[string]*Chat
	mutex sync.Mutex
}

var chatList ChatList

func (chL *ChatList) CreateChat(chName string) *Chat {
	chL.mutex.Lock()
	defer chL.mutex.Unlock()

	chat := &Chat{
		Name:    chName,
		members: map[string]*Client{},
	}
	chL.Chats[chName] = chat
	return chat
}

func (chL *ChatList) initializeRooms(chatHandler ChatDBhandler) {
	chL.Chats = make(map[string]*Chat)
	list, _ := chatHandler.GetAllChats()

	for _, chat := range list {
		chat := &Chat{
			Name:        chat.Name,
			subscribers: chat.Subscribers,
			members:     make(map[string]*Client),
		}
		chL.Chats[chat.Name] = chat
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

func (chL *ChatList) removeClient(cl *Client) {
	for _, sub := range cl.subs {
		chat := chL.Chats[sub]
		chat.removeMember(cl.username)
	}
}

func (ch *Chat) AddMember(cl *Client) {
	ch.mutex.Lock()
	defer ch.mutex.Unlock()
	ch.members[cl.username] = cl
	fmt.Println(cl.username, "client added to", ch.Name)
}

func (ch *Chat) removeMember(username string) {
	for k := range ch.members {
		if k == username {
			delete(ch.members, k)
			break
		}
	}

}

func (c *Chat) checkOnline() map[string]bool {
	resultUsers := make(map[string]bool)
	for _, username := range c.subscribers {
		_, ok := c.members[username]
		if !ok {
			resultUsers[username] = false
		} else {
			resultUsers[username] = true
		}
	}
	return resultUsers
}
