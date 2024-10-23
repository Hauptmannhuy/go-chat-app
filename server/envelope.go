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
	UserID string
}

type JoinNotification struct {
	ChatID string `json:"chatID"`
	UserID string `json:"userID"`
}

type Error struct {
	Message string
}

func handleResponseEnvelope(outEnv OutEnvelope, connSockets *Hub, msgT int, chats *ChatList, cl *Client) {
	fmt.Println(cl.index, "client index in handle response env")
	jsonEnv, err := json.Marshal(outEnv)
	fmt.Println("slice of sockets:", connSockets.Connections)
	if err != nil {
		fmt.Println(err)
		return
	}

	switch outEnv.Type {
	case "NEW_MESSAGE":
		msg := outEnv.Data.(UserMessage)

		chatID := msg.ChatID
		chat := chatList.Chats[chatID]

		for _, cl := range chat.members {
			sendWsResponse(jsonEnv, cl, msgT)
		}
	case "NEW_CHAT":
		for i := 0; i < len(connSockets.Connections); i++ {
			fmt.Println(len(connSockets.Connections), "connections")
			connCl := connSockets.Connections[i]
			sendWsResponse(jsonEnv, connCl, msgT)
		}
		msg := outEnv.Data.(Chat)
		chats.CreateChat(msg.ID)
	case "JOIN_CHAT":
		msg := outEnv.Data.(JoinNotification)
		chat := chatList.Chats[msg.ChatID]
		chat.AddMember(cl)
	case "ERROR":
		sendWsResponse(jsonEnv, cl, msgT)
	}

}

func sendWsResponse(p []byte, cl *Client, msgT int) {
	socket := cl.socket
	if err := socket.WriteMessage(msgT, p); err != nil {
		log.Println("Error writing to WebSocket:", err)
		return
	}
	fmt.Println("Message sent successfully to client", string(p))
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
		fmt.Println("res", outEnv.Data)
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
