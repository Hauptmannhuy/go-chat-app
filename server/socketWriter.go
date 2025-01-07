package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/gorilla/websocket"
)

type UserStatus struct {
	Status map[string]string
}

func writeToSocket(message interface{}, messageType string, cl *Client, wsMsgType int) {
	outEnv := OutEnvelope{
		Type: messageType,
		Data: message,
	}
	p, err := json.Marshal(outEnv)
	if err != nil {
		log.Fatal("error in send ws response", err)
	}
	socket := cl.socket
	if err := socket.WriteMessage(wsMsgType, p); err != nil {
		log.Println("Error writing to WebSocket:", err)
		return
	}
	fmt.Println("Message sent successfully to client", cl.username, "Message:", messageType)
}

func broadcastUserStatus(status string, peer *Client) {

	statusInfo := UserStatus{
		Status: map[string]string{peer.username: status},
	}
	for _, sub := range peer.subs {
		chat := chatList.Chats[sub]
		for _, peerToWrite := range chat.members {
			if peer.username == peerToWrite.username {
				continue
			}
			writeToSocket(statusInfo, "USER_STATUS", peerToWrite, websocket.TextMessage)
		}
	}
}

func loadPeersStatus(peerToNotify *Client) {
	statusInfo := UserStatus{
		Status: make(map[string]string),
	}

	for _, sub := range peerToNotify.subs {
		onlineUsers := chatList.Chats[sub].checkOnline()
		for name, isOnline := range onlineUsers {
			if isOnline {
				statusInfo.Status[name] = "online"
			} else {
				statusInfo.Status[name] = "offline"
			}
		}
	}
	writeToSocket(statusInfo, "USER_STATUS", peerToNotify, websocket.TextMessage)
}
