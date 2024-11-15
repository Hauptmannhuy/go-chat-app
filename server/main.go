package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sync"

	_ "github.com/golang-migrate/migrate/v4/source/file" // Import the file source driver
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"

	"github.com/gorilla/websocket"
)

var redisDB *redis.Client

var upgrader = websocket.Upgrader{

	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func main() {

	redisDB = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
		Protocol: 3,
	})

	redisDB.Set(context.Background(), "test", "value", 0)

	s, err := redisDB.Get(context.Background(), "test").Result()

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(s)

	err = dbManager.openAndMigrateDB()
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

		wsMessageType, p, err := peer.ReadMessage()
		cl.mutex.Unlock()
		if err != nil {
			log.Println(err)
			return
		}
		data, dataType := processMessage(p, cl)
		data, err = dbManager.handleDatabase(data)
		if err != nil {
			data = Error{err.Error()}
			dataType = "ERROR"
			fmt.Println(err)
		}
		HandleWriteToWebSocket(dataType, data, wsMessageType, cl)

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
	h.Mutex.Lock()
	defer h.Mutex.Unlock()
	h.Connections[c.username] = c

	fmt.Println("connected clients:", h.Connections)
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

	sendWsResponse(chatContainer, "LOAD_SUBS", cl, websocket.TextMessage)
}

func (cl *Client) sendMessageHistory() {
	dbMessageHandler := dbManager.initializeDBhandler("message")
	data, err := dbMessageHandler.GetChatsMessages(cl.subs)
	if err != nil {
		fmt.Println(err)
		return
	}
	sendWsResponse(data, "LOAD_MESSAGES", cl, websocket.TextMessage)
}

func (h *Hub) initialize() {
	h.Connections = make(map[string]*Client)
}
