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

	chatList.initializeRooms()
	err = http.ListenAndServe(":8090", NewAuthMiddlewareHandler(AuthHandler{}))

	if err != nil {
		fmt.Println(err)
	}
}

type Client struct {
	index     int
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
	subHandler := dbManager.initializeDBhandler("subscription")
	subs, err := subHandler.LoadUserSubscriptionsHandler(usernameCookie.Value)
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
		peerSoc := cl.socket

		fmt.Println(cl.connected)

		messageType, p, err := peerSoc.ReadMessage()
		cl.mutex.Unlock()
		if err != nil {
			log.Println(err)
			return
		}
		outEnv := processEnvelope(p)
		outEnv.Data, err = dbManager.handleDataBase(outEnv)
		if err != nil {
			errorMessg := Error{err.Error()}
			outEnv = OutEnvelope{"ERROR", errorMessg}
			fmt.Println(err)
		}
		handleResponseEnvelope(outEnv, &connSockets, messageType, &chatList, cl)

	}
}

type Hub struct {
	Connections []*Client
	Mutex       sync.Mutex
}

var connSockets Hub

func (hub *Hub) removeClient(cl *Client) {
	hub.Mutex.Lock()
	defer hub.Mutex.Unlock()
	fmt.Println("removing client..")
	for i := range hub.Connections {
		fmt.Println(i)
		if cl.socket == hub.Connections[i].socket {
			fmt.Println("Removed")
			hub.Connections = append(hub.Connections[:i], hub.Connections[i+1:]...)
			return
		}
	}
}

func (h *Hub) AddConection(c *Client) {
	fmt.Println("connected sockets:", h.Connections)
	h.Mutex.Lock()
	defer h.Mutex.Unlock()
	h.Connections = append(h.Connections, c)

}

func (cl *Client) sendSubscriptions() {
	var env OutEnvelope
	env.Type = "LOAD_SUBS"
	env.Data = cl.subs

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
