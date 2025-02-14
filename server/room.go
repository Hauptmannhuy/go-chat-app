package main

import (
	"fmt"
	"strings"
)

type Room struct {
	read         chan OutEnvelope
	write        chan string
	newPeerReq   chan *Client
	leavePeerReq chan int
	done         chan struct{}

	connections map[int]*Client
	roomReg     *roomRegister

	subscribers []int
	name        string
}

func (room *Room) run() {
	defer func() {
		fmt.Println("closing room")
		close(room.leavePeerReq)
		close(room.read)
		close(room.newPeerReq)
		close(room.write)
	}()
	for {
		select {
		case peerID := <-room.leavePeerReq:
			if _, ok := room.connections[peerID]; ok {
				delete(room.connections, peerID)
				if len(room.connections) == 0 {
					fmt.Println(room.name)
					room.roomReg.closeRoomChan <- room.name
				}
			}
		case message := <-room.read:
			for _, id := range room.subscribers {
				peer, ok := room.connections[id]
				if ok {
					peer.messageBuffer <- message
				} else {
					if asserted, ok := message.Data.(redisBuffer); ok {
						fmt.Println("saving to redis")
						asserted.saveToBuff(id)
					}
				}
			}
		case newPeer := <-room.newPeerReq:
			room.connections[newPeer.id] = newPeer
		case <-room.done:
			return
		}

	}
}

func newRoom(name string, roomReg *roomRegister) *Room {
	handler := getDB().initializeDBhandler("subscription")
	var subscribers []int
	if len(strings.Split(name, "_")) > 1 {
		subscribers = handler.GetPrivateChatSubs(name)
	} else {
		subscribers = handler.GetGroupChatSubs(name)
	}

	return &Room{
		name:         name,
		newPeerReq:   make(chan *Client),
		leavePeerReq: make(chan int),
		read:         make(chan OutEnvelope),
		write:        make(chan string),
		connections:  map[int]*Client{},
		done:         make(chan struct{}),
		roomReg:      roomReg,
		subscribers:  subscribers,
	}
}
