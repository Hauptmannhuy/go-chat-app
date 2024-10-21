package main

import (
	"encoding/json"
	"fmt"
	"log"
)

type OutEnvelope struct {
	Type string
	Data interface{}
}

type Envelope struct {
	Type string
}

type UserMessage struct {
	Body   string
	ChatID string
	// UserID string
}

type JoinNotification struct {
	ChatID string `json:"chatID"`
}

type Error struct {
	Message string
}

func handleResponseEnvelope(outEnv OutEnvelope, connSockets *Hub, msgT int, chats *ChatList, cl *Client) {
	fmt.Println(outEnv)
	jsonEnv, err := json.Marshal(outEnv)

	if err != nil {
		fmt.Println(err)
		return
	}

	switch outEnv.Type {
	case "NEW_MESSAGE":
		msg, ok := outEnv.Data.(UserMessage)
		if ok {
			fmt.Println("Message body:", msg.Body)
			fmt.Println("Chat ID:", msg.ChatID)
		}
		chatID := msg.ChatID
		chat := chatList.Chats[chatID]
		for _, cl := range chat.members {
			sendResposeEnvelope(jsonEnv, cl, msgT)
		}
	case "NEW_CHAT":
		for _, cl := range connSockets.Connections {
			sendResposeEnvelope(jsonEnv, cl, msgT)
		}
		msg := outEnv.Data.(Chat)
		chats.CreateChat(msg.id)
	case "JOIN_CHAT":
		msg := outEnv.Data.(JoinNotification)
		chat := chatList.Chats[msg.ChatID]
		chat.AddMember(cl)
		fmt.Println(chat, "Chat members")
	case "ERROR":
		sendResposeEnvelope(jsonEnv, cl, msgT)
	}

}

func sendResposeEnvelope(p []byte, cl *Client, msgT int) {
	cl.mutex.Lock()
	defer cl.mutex.Unlock()

	socket := cl.socket
	if err := socket.WriteMessage(msgT, p); err != nil {
		log.Println(err)
		return
	}
}

func processEnvelope(p []byte) OutEnvelope {
	var outEnv OutEnvelope
	var env Envelope

	fmt.Println("Raw JSON", string(p))
	if err := json.Unmarshal(p, &env); err != nil {
		log.Fatal(err)
	}
	outEnv.Type = env.Type
	switch env.Type {
	case "NEW_MESSAGE":
		var s struct {
			Envelope
			UserMessage
		}
		if err := json.Unmarshal(p, &s); err != nil {
			log.Fatal(err)
		}
		outEnv.Data = s.UserMessage
		return outEnv
	case "NEW_CHAT":
		var s struct {
			Envelope
			Chat
		}
		if err := json.Unmarshal(p, &s); err != nil {
			log.Fatal(err)
		}
		outEnv.Data = s.Chat
		return outEnv
	case "JOIN_CHAT":
		var s struct {
			Envelope
			JoinNotification
		}
		if err := json.Unmarshal(p, &s); err != nil {
			log.Fatal(err)
		}
		outEnv.Data = s.JoinNotification
		return outEnv
	default:
		fmt.Println("No type is matched while processing incoming envelope")
		return outEnv
	}
}
