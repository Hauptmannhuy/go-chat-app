package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type Client struct {
	Socket    *websocket.Conn
	Connected bool
	Mutex     sync.Mutex
	// Subscriptions []string
}

type OutEnvelope struct {
	Type string
	Data interface{}
}

type UserMessage struct {
	Body   string
	ChatID string
	// UserID string
}

type Envelope struct {
	Type string
}

type Chat struct {
	Members []*Client
	ID      string
	Mutex   sync.Mutex
}

func (c *Chat) AddMember(cl *Client) {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()
	c.Members = append(c.Members, cl)
}

type JoinNotification struct {
	ChatID string `json:"chatID"`
}

type Hub struct {
	Connections []*Client
	Mutex       sync.Mutex
}

type ChatList struct {
	Chats map[string]*Chat
	Mutex sync.Mutex
}

func (chL *ChatList) CreateChat(chID string) {
	chL.Mutex.Lock()
	defer chL.Mutex.Unlock()
	if chL.Chats == nil {
		chL.Chats = make(map[string]*Chat)
	}
	chat := &Chat{
		ID: chID,
	}
	chL.Chats[chID] = chat
}

var chatList ChatList

// func (c *Client) AddSub(cID string) {
// 	c.Subscriptions = append(c.Subscriptions, cID)
// }

func (c *Client) CloseConnection() {
	c.Connected = false
}

var connSockets Hub

func (cnSockets *Hub) removeClient(cl *Client) {
	cnSockets.Mutex.Lock()
	defer cnSockets.Mutex.Unlock()
	fmt.Println("removing client..")
	for i := range cnSockets.Connections {
		fmt.Println(i)
		if cl.Socket == cnSockets.Connections[i].Socket {
			fmt.Println("Removed")
			cnSockets.Connections = append(cnSockets.Connections[:i], cnSockets.Connections[i+1:]...)
			return
		}
	}
}

func (h *Hub) AddConection(c *Client) {
	h.Mutex.Lock()
	defer h.Mutex.Unlock()
	h.Connections = append(h.Connections, c)
}

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

	setCorsHeaders(w)
	conn, err := upgrader.Upgrade(w, r, nil)
	var newClient = &Client{
		Socket:    conn,
		Connected: true,
	}

	connSockets.AddConection(newClient)

	newClient.Socket.SetCloseHandler(func(code int, text string) error {
		newClient.CloseConnection()
		return nil
	})

	if err != nil {
		log.Println(err)
		return
	}

	fmt.Println(connSockets)

	go clientMessages(newClient)
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
		handleResponseEnvelope(p, &connSockets, messageType, &chatList, cl)

	}
}

func handleResponseEnvelope(p []byte, connSockets *Hub, msgT int, chats *ChatList, cl *Client) {
	outEnv := processEnvelope(p)
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
		for _, cl := range chat.Members {
			sendResposeEnvelope(jsonEnv, cl, msgT)
		}
	case "NEW_CHAT":
		for _, cl := range connSockets.Connections {
			sendResposeEnvelope(jsonEnv, cl, msgT)
		}
		msg := outEnv.Data.(*Chat)
		chats.CreateChat(msg.ID)
	case "JOIN_CHAT":
		msg := outEnv.Data.(JoinNotification)
		chat := chatList.Chats[msg.ChatID]
		chat.AddMember(cl)
		fmt.Println(chat, "Chat members")
	}

}

func sendResposeEnvelope(p []byte, cl *Client, msgT int) {
	cl.Mutex.Lock()
	defer cl.Mutex.Unlock()

	socket := cl.Socket
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
		return outEnv
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
