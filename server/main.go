package main

import (
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
	var err error
	db, err = openAndMigrateDB()
	if err != nil {
		fmt.Println(err)
		return
	}
	err = http.ListenAndServe(":8090", NewAuthMiddlewareHandler(AuthHandler{}))

	if err != nil {
		fmt.Println(err)
	}
}

type Client struct {
	Socket    *websocket.Conn
	Connected bool
	Mutex     sync.Mutex
}

func (c *Client) CloseConnection() {
	c.Connected = false
}

func clientMessages(cl *Client) {
	defer func() {
		cl.Socket.Close()
		connSockets.removeClient(cl)
	}()
	defer fmt.Println("Connection closed with", cl)

	for {
		cl.Mutex.Lock()
		peerSoc := cl.Socket

		fmt.Println(cl.Connected)

		messageType, p, err := peerSoc.ReadMessage()
		cl.Mutex.Unlock()
		if err != nil {
			log.Println(err)
			return
		}
		outEnv := processEnvelope(p)
		err = handleDataBase(outEnv)
		if err != nil {
			errorMessg := Error{err.Error()}
			outEnv = OutEnvelope{"ERROR", errorMessg}
			fmt.Println(outEnv.Data, "outenv data")
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
		if cl.Socket == hub.Connections[i].Socket {
			fmt.Println("Removed")
			hub.Connections = append(hub.Connections[:i], hub.Connections[i+1:]...)
			return
		}
	}
}

func (h *Hub) AddConection(c *Client) {
	h.Mutex.Lock()
	defer h.Mutex.Unlock()
	h.Connections = append(h.Connections, c)
}
