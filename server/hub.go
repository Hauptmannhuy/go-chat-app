package main

import (
	"fmt"
	"net/http"
)

type hub struct {
	roomRegister *roomRegister

	wsMessageChan chan *wsMessage

	peerConnect    chan *Client
	peerDisconnect chan int
	newRoomChan    chan *Room
	closeRoomChan  chan string

	connections map[int]*Client
}

type roomRegister struct {
	rooms         map[string]*Room
	closeRoomChan chan string
}

type closeRoomReq struct {
	roomID string
}

func newHub() *hub {
	return &hub{
		roomRegister: &roomRegister{rooms: make(map[string]*Room), closeRoomChan: make(chan string)},

		peerDisconnect: make(chan int),
		peerConnect:    make(chan *Client),
		newRoomChan:    make(chan *Room),
		connections:    map[int]*Client{},
		wsMessageChan:  make(chan *wsMessage),
	}
}

func (hub *hub) handleClientConn(w http.ResponseWriter, r *http.Request) {
	client := initClient(w, r, hub)
	// notify other clients about new connection status
	// load other client client connection status and notify connected user
	fmt.Println("start client goroutine")
	go client.run()

	hub.peerConnect <- client

}

func (h *hub) run() {
	fmt.Println("start of hub goroutine")
	for {
		select {
		case peerReq := <-h.peerConnect:
			{
				handler := getDB().initializeDBhandler("subscription")
				subs, _ := handler.LoadSubscriptions(peerReq.id)
				peerReq.handleOfflineMessages(subs)
				h.broadcastUserStatus(peerReq, "online")
				h.connectPeer(peerReq, subs)
			}
		case id := <-h.peerDisconnect:
			{
				delete(h.connections, id)
				for _, room := range h.roomRegister.rooms {
					room.leavePeerReq <- id
				}
				fmt.Println("Connected peers", h.connections)
			}
		case closeRoomReq := <-h.roomRegister.closeRoomChan:
			{
				close(h.roomRegister.rooms[closeRoomReq].done)
				h.roomRegister.shutdownRoom(closeRoomReq)
			}
		case wsMsg := <-h.wsMessageChan:
			{
				wsMsg.broadcastHandler.exec(wsMsg.owner, wsMsg.payload)
			}
		}
	}
}

func (h *hub) connectPeer(client *Client, subs []string) {
	h.connections[client.id] = client

	for _, sub := range subs {
		if _, ok := h.roomRegister.rooms[sub]; ok {
			h.roomRegister.addClient(sub, client)

		} else {
			h.roomRegister.addRoom(sub)
			h.roomRegister.addClient(sub, client)
		}
	}
}

func (roomReg *roomRegister) addRoom(name string) {
	room := newRoom(name, roomReg)
	roomReg.rooms[name] = room
	go room.run()
}

func (roomReg *roomRegister) addClient(roomName string, client *Client) {
	roomReg.rooms[roomName].newPeerReq <- client
}

func (roomReg *roomRegister) shutdownRoom(name string) {
	delete(roomReg.rooms, name)
}

func (hub *hub) removeClient(cl *Client) {

}

func (hub *hub) broadcastUserStatus(client *Client, status string) {
	outEnv := OutEnvelope{
		Type: "USER_STATUS",
		Data: &UserStatus{
			Status: map[string]string{client.username: status},
		},
	}
	sended := map[int]int{}
	handler := getDB().initializeDBhandler("subscription")
	subscriptions, _ := handler.LoadSubscriptions(client.id)

	for _, sub := range subscriptions {
		if room, ok := hub.roomRegister.rooms[sub]; ok {
			for _, peer := range room.connections {

				if _, sent := sended[peer.id]; sent {
					continue
				}
				peer.socket.WriteJSON(outEnv)
				sended[peer.id] = peer.id
			}
		}
	}
}
