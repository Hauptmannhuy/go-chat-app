package main

import (
	"fmt"
	"sync"
)

type Hub struct {
	Connections map[string]*Client
	Mutex       sync.Mutex
}

var connSockets Hub

func (hub *Hub) removeClient(cl *Client) {
	hub.Mutex.Lock()
	defer hub.Mutex.Unlock()
	fmt.Println("removing client..")
	delete(hub.Connections, cl.username)

}

func (h *Hub) AddHubMember(c *Client) {
	h.Mutex.Lock()
	defer h.Mutex.Unlock()
	h.Connections[c.username] = c

	fmt.Println("connected clients:", h.Connections)
}

func (h *Hub) initialize() {
	h.Connections = make(map[string]*Client)
}

func (h *Hub) isUserOnline(peerName string) bool {
	_, ok := h.Connections[peerName]
	return ok
}
