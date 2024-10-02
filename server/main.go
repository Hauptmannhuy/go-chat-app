package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type Room struct {
	Members []*Client
}

type Client struct {
	Socket    *websocket.Conn
	Connected bool
}

func (c *Client) CloseConnection() {
	c.Connected = false
}

var hub Room

type ApiResponse struct {
	Message string `json:"message"`
}

var upgrader = websocket.Upgrader{

	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func main() {
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/chat", chatHandler)
	err := http.ListenAndServe(":8090", nil)
	if err != nil {
		fmt.Println(err)
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "OPTIONS" {
		setOptions(w)
	} else if r.Method == "GET" {
		getHome(w, r)
	}
}

func chatHandler(w http.ResponseWriter, r *http.Request) {
	roomID := r.URL.Query().Get("roomId")
	fmt.Println(roomID, "room id")
	if roomID == "" {
		http.Error(w, "Missing roomID", http.StatusBadRequest)
		return
	}

	setCorsHeaders(w)
	conn, err := upgrader.Upgrade(w, r, nil)
	var newClient = Client{
		Socket:    conn,
		Connected: true,
	}
	hub.Members = append(hub.Members, &newClient)

	newClient.Socket.SetCloseHandler(func(code int, text string) error {
		newClient.CloseConnection()
		return nil
	})

	if err != nil {
		log.Println(err)
		return
	}
	conn.WriteMessage(websocket.TextMessage, []byte("Bidirectional connection established"))

	fmt.Println(hub.Members)

	go clientMessages(newClient)
}

func clientMessages(cl Client) {

	defer func() {
		cl.Socket.Close()
		removeClient(cl)
	}()

	defer fmt.Println("Connection closed with", cl)
	for {
		peerSoc := cl.Socket
		fmt.Println(cl.Connected)
		messageType, p, err := peerSoc.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		for _, hubCl := range hub.Members {
			hubSoc := hubCl.Socket
			if err := hubSoc.WriteMessage(messageType, p); err != nil {
				log.Println(err)
				return
			}

		}
		fmt.Println("Message from", cl)

	}
}

func removeClient(cl Client) {
	fmt.Println("removing..")
	for i := range hub.Members {
		fmt.Println(i)
		if cl.Socket == hub.Members[i].Socket {
			fmt.Println("Found")
			hub.Members = append(hub.Members[:i], hub.Members[i+1:]...)
			fmt.Println(hub.Members)
			return
		}
	}
}

func getHome(w http.ResponseWriter, r *http.Request) {

	setCorsHeaders(w)

	if err := r; err != nil {
		fmt.Println(err)
	}

	response := ApiResponse{"Hello from backend!"}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func setOptions(w http.ResponseWriter) {
	setCorsHeaders(w)
	w.WriteHeader(http.StatusOK)
}

func setCorsHeaders(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
}
