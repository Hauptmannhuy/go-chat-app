package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

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

	if err != nil {
		log.Println(err)
		return
	}
	conn.WriteMessage(websocket.TextMessage, []byte("Bidirectional connection established"))
	defer conn.Close()
	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		if err := conn.WriteMessage(messageType, p); err != nil {
			log.Println(err)
			return
		}
		fmt.Println("Message data", messageType, p)
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
