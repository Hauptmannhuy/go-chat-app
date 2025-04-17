package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type Client struct {
	id            int
	username      string
	socket        *websocket.Conn
	messageBuffer chan OutEnvelope
	hub           *hub
	done          chan struct{}
}

type UserStatus struct {
	Status map[string]string
}

func initClient(w http.ResponseWriter, r *http.Request, h *hub) *Client {
	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Println(err)
		return nil
	}

	token, _ := r.Cookie("token")
	username, _ := r.Cookie("username")
	userIndex := fetchUserID(token.Value)

	client := &Client{
		socket: conn,

		messageBuffer: make(chan OutEnvelope),
		done:          make(chan struct{}),

		id:       userIndex,
		hub:      h,
		username: username.Value,
	}

	conn.SetCloseHandler(func(code int, text string) error {
		client.hub.peerDisconnect <- client.id
		close(client.messageBuffer)
		close(client.done)
		return nil
	})

	return client
}

func (client *Client) run() {

	go func() {
		for {
			_, buffer, err := client.socket.ReadMessage()
			if err != nil {
				fmt.Println(err)
				return
			}
			msg, _ := client.processMessage(buffer)
			client.hub.wsMessageChan <- msg
		}
	}()

	for {
		select {
		case message := <-client.messageBuffer:
			err := client.socket.WriteJSON(message)
			fmt.Println("write to socket")
			if err != nil {
				fmt.Println("Error writing to socket", err)
			}
		case <-client.done:
			return
		}
	}

}

func (client *Client) processMessage(p []byte) (*wsMessage, error) {
	env := JSONenvelope{}
	err := json.Unmarshal(p, &env)

	if err != nil {
		log.Println(err)
		ok := isTypeUnknown(err.Error())
		if ok {
			return &wsMessage{
				payload:          "Unknown type",
				broadcastHandler: defineAlgo(""),
			}, fmt.Errorf("Unknown type")
		} else {
			log.Fatal(err)
		}
	}

	msg, err := kindHandler(env.Type)
	if err != nil {
		fmt.Println(err)
	}

	err = json.Unmarshal(p, msg)
	fmt.Println(string(p))

	if err != nil {
		log.Fatal(err)
	}
	broadcastHandler := defineAlgo(msg)

	msg.Process(client)

	return &wsMessage{
		payload:          msg,
		owner:            client,
		broadcastHandler: broadcastHandler,
	}, nil
}

func (cl *Client) handleOfflineMessages(clientSubs []string) {
	redisManager := getRedis()
	data := map[string][]interface{}{}
	fmt.Println(data)
	for _, sub := range clientSubs {
		key := fmt.Sprintf("offline:messages:%s:%d", sub, cl.id)
		if ok := redisManager.hasMessages(key); ok {
			data[sub] = redisManager.getOffMessages(key)
		}
	}
	if len(data) > 0 {
		cl.socket.WriteJSON(OutEnvelope{
			Type: "OFFLINE_MESSAGES",
			Data: data,
		})
	}
}
