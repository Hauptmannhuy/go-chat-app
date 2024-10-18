package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"go-chat-app/dbmanager/handler"
	"go-chat-app/dbmanager/service"
	"go-chat-app/dbmanager/store"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file" // Import the file source driver
	"github.com/joho/godotenv"
	_ "github.com/lib/pq" // PostgreSQL driver

	// PostgreSQL driver
	"github.com/gorilla/websocket"
)

type Client struct {
	Socket    *websocket.Conn
	Connected bool
	Mutex     sync.Mutex
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

type Error struct {
	Message string
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

func (c *Client) CloseConnection() {
	c.Connected = false
}

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

var upgrader = websocket.Upgrader{

	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type AuthorizationMiddleware struct {
	handler http.Handler
}

type AuthHandler struct{}

func (h AuthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

}

func NewAuthMiddlewareHandler(handler http.Handler) AuthorizationMiddleware {
	return AuthorizationMiddleware{
		handler: handler,
	}
}

func (am AuthorizationMiddleware) ServeHTTP(w http.ResponseWriter, req *http.Request) {

	if req.URL.Path == "/sign_in" {
		signInHandler(w, req)
		return
	}
	if req.URL.Path == "/sign_up" {
		signUpHandler(w, req)
		return
	}
	if req.URL.Path == "/sign_out" {
		SignOutHandler(w, req)
		return
	}

	cookie, err := req.Cookie("token")

	if err != nil {
		fmt.Println(err)
		fmt.Println("redirect")
		http.Redirect(w, req, "/sign_up", http.StatusSeeOther)
	} else {
		ok := verifyToken(cookie)
		if !ok {
			http.Redirect(w, req, "/sign_up", http.StatusSeeOther)
		} else {
			chatHandler(w, req)
		}
	}

}

func verifyToken(c *http.Cookie) bool {
	tokenS := c.Value

	token, err := jwt.Parse(tokenS, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("KEY")), nil
	})

	if err != nil {
		log.Fatal(err)
		return false
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		fmt.Println(claims["iss"])
	} else {
		fmt.Println(err)
	}
	return true
}

var chatList ChatList

var connSockets Hub

var db *sql.DB

func main() {

	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	dotenv := os.Getenv("DATABASE_CREDS")
	fmt.Println(dotenv)
	dataSourceName := fmt.Sprintf("postgres://%s/dbmanager?sslmode=disable", dotenv)
	db, err = sql.Open("postgres", dataSourceName)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	fmt.Println("Connected to database")

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatalf("Could not start SQL driver: %v", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations", // Correctly specify the file scheme
		"postgres", driver)
	if err != nil {
		log.Fatalf("Could not start migration: %v", err)
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Migration failed: %v", err)
	}

	err = http.ListenAndServe(":8090", NewAuthMiddlewareHandler(AuthHandler{}))

	if err != nil {
		fmt.Println(err)
	}
}

func signUpHandler(w http.ResponseWriter, r *http.Request) {

	fmt.Println("sign up")

	data, err := io.ReadAll(r.Body)

	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var message struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	json.Unmarshal(data, &message)

	userHandler := initializeDBhandler(db, "user")

	err = userHandler.CreateUserHandler(message.Username, message.Email, message.Password)

	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	s, err := generateToken()

	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	setAuthCookies(w, s, message.Username)

	w.WriteHeader(http.StatusOK)
}

func initializeDBhandler(db *sql.DB, handlerDeclaration string) handler.Handler {
	dbStore := store.SQLstore{DB: db}
	service := service.Service{}
	handler := handler.Handler{}
	if handlerDeclaration == "user" {
		service.UserStore = &dbStore
		handler.UserService = service
	} else if handlerDeclaration == "chat" {
		service.ChatStore = &dbStore
		handler.ChatService = service
	} else if handlerDeclaration == "message" {
		service.MessageStore = &dbStore
		handler.MessageService = service
	}
	return handler
}

func setAuthCookies(w http.ResponseWriter, s string, name string) {

	http.SetCookie(w, &http.Cookie{
		Name:     "username",
		Value:    name,
		Expires:  time.Now().AddDate(0, 0, 14),
		Path:     "/",
		SameSite: http.SameSiteStrictMode,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    s,
		Expires:  time.Now().AddDate(0, 0, 14),
		HttpOnly: true,
		Path:     "/",
		SameSite: http.SameSiteStrictMode,
	})
}

func SignOutHandler(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Expires:  time.Now(),
		HttpOnly: true,
		Path:     "/",
		SameSite: http.SameSiteStrictMode,
	})
	http.SetCookie(w, &http.Cookie{
		Name:     "username",
		Expires:  time.Now(),
		Path:     "/",
		SameSite: http.SameSiteStrictMode,
	})

	w.WriteHeader(http.StatusOK)
}

func signInHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("sign in")
	userHandler := initializeDBhandler(db, "user")
	var message struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	data, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = json.Unmarshal(data, &message)

	if err != nil {
		fmt.Println(err)
	}

	err = userHandler.LoginUserHandler(message.Username, message.Password)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	} else {
		token, _ := generateToken()
		setAuthCookies(w, token, message.Username)
		w.WriteHeader(http.StatusOK)
	}
}

func generateToken() (string, error) {
	key := []byte(os.Getenv("KEY"))
	t := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"iss": "go-chat-app",
		})
	s, err := t.SignedString(key)

	if err != nil {
		return "", err
	}
	return s, nil
}

func chatHandler(w http.ResponseWriter, r *http.Request) {
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

func handleDataBase(env OutEnvelope) error {
	jsoned, _ := json.Marshal(env.Data)
	switch env.Type {
	case "NEW_MESSAGE":
		messageHandler := initializeDBhandler(db, "message")
		err := messageHandler.CreateMessageHandler(jsoned)
		return err
	case "NEW_CHAT":
		chatHandler := initializeDBhandler(db, "chat")
		err := chatHandler.CreateChatHandler(jsoned)
		return err
	default:
		fmt.Println("No write to database")
		return nil
	}
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
		for _, cl := range chat.Members {
			sendResposeEnvelope(jsonEnv, cl, msgT)
		}
	case "NEW_CHAT":
		for _, cl := range connSockets.Connections {
			sendResposeEnvelope(jsonEnv, cl, msgT)
		}
		msg := outEnv.Data.(Chat)
		chats.CreateChat(msg.ID)
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
		fmt.Println("No type is matched while processing incoming envelope")
		return outEnv
	}
}
