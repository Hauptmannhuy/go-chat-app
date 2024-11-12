package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	_ "github.com/golang-migrate/migrate/v4/source/file" // Import the file source driver
	_ "github.com/lib/pq"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{

	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func main() {
	err := dbManager.openAndMigrateDB()
	if err != nil {
		fmt.Println(err)
		return
	}

	chatDBhandler := dbManager.initializeDBhandler("chat")
	chatList.initializeRooms(&chatDBhandler)

	connSockets.initialize()
	err = http.ListenAndServe(":8090", NewAuthMiddlewareHandler())

	if err != nil {
		fmt.Println(err)
	}
}

type Client struct {
	index     string
	username  string
	socket    *websocket.Conn
	connected bool
	mutex     sync.Mutex
	subs      []string
}

func initializeWSconn(w http.ResponseWriter, r *http.Request) *Client {

	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Println(err)
		return nil

	}
	usernameCookie, _ := r.Cookie("username")
	token, _ := r.Cookie("token")
	userIndex := fetchUserID(token.Value)
	subHandler := dbManager.initializeDBhandler("subscription")
	subs, err := subHandler.LoadSubscriptions(userIndex)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return &Client{
		connected: true,
		username:  usernameCookie.Value,
		socket:    conn,
		mutex:     sync.Mutex{},
		subs:      subs,
		index:     userIndex,
	}
}

func (c *Client) CloseConnection() {
	c.connected = false
}

func clientMessages(cl *Client) {
	defer func() {
		cl.socket.Close()
		connSockets.removeClient(cl)
	}()
	defer fmt.Println("Connection closed with", cl)
	for {
		cl.mutex.Lock()
		peer := cl.socket

		fmt.Println(cl.connected)

		messageType, p, err := peer.ReadMessage()
		cl.mutex.Unlock()
		if err != nil {
			log.Println(err)
			return
		}
		outEnv := processEnvelope(p, cl)
		outEnv.Data, err = dbManager.handleDatabase(outEnv)
		if err != nil {
			errorMessg := Error{err.Error()}
			outEnv = OutEnvelope{"ERROR", errorMessg}
			fmt.Println(err)
		}
		HandleWriteToWebSocket(outEnv, messageType, cl)

	}
}

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
	fmt.Println("connected clients:", h.Connections)
	h.Mutex.Lock()
	defer h.Mutex.Unlock()
	h.Connections[c.username] = c

}

func (cl *Client) sendSubscribedChats() {

	dbChatHandler := dbManager.initializeDBhandler("chat")
	groupChatData, err := dbChatHandler.LoadUserSubscribedChats(cl.index)
	if err != nil {
		fmt.Println(err)
		return
	}
	privateChatData, err := dbChatHandler.LoadSubscribedPrivateChats(cl.index)
	if err != nil {
		fmt.Println(err)
		return
	}
	chatContainer := make(map[string]interface{})
	chatContainer["private"] = privateChatData
	chatContainer["group"] = groupChatData
	var env OutEnvelope
	env.Type = "LOAD_SUBS"
	env.Data = chatContainer

	j, err := json.Marshal(env)
	if err != nil {
		fmt.Println(err)
		return
	}
	sendWsResponse(j, cl, websocket.TextMessage)
}

func (cl *Client) sendMessageHistory() {
	dbMessageHandler := dbManager.initializeDBhandler("message")
	data, err := dbMessageHandler.GetChatsMessages(cl.subs)
	if err != nil {
		fmt.Println(err)
		return
	}
	var env = OutEnvelope{
		Type: "LOAD_MESSAGES",
		Data: data,
	}
	json, err := json.Marshal(env)
	if err != nil {
		fmt.Println(err)
		return
	}
	sendWsResponse(json, cl, websocket.TextMessage)
}

func (h *Hub) initialize() {
	h.Connections = make(map[string]*Client)
}
