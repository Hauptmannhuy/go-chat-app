package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	_ "github.com/golang-migrate/migrate/v4/source/file" // Import the file source driver
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"

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
	redisManager = redisWrapper{
		redis: redis.NewClient(&redis.Options{
			Addr:     "localhost:6379",
			Password: "",
			DB:       0,
			Protocol: 3,
		}),
	}
	err := dbManager.openAndMigrateDB()
	if err != nil {
		log.Println(err)
		return
	}

	chatDBhandler := dbManager.initializeDBhandler("chat")
	chatList.initializeRooms(&chatDBhandler)

	connSockets.initialize()
	err = http.ListenAndServe(":8090", NewAuthMiddlewareHandler())

	if err != nil {
		log.Println(err)
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
		log.Println(err)
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
		chatList.removeClient(cl)
		broadcastUserStatus("offline", cl)
	}()
	defer fmt.Println("Connection closed with", cl)
	for {
		cl.mutex.Lock()
		peer := cl.socket

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
			log.Println(err)
		}
		dispatchAction(dataType, data, wsMessageType, cl)

	}
}
