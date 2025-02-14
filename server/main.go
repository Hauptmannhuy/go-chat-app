package main

import (
	"log"
	"net/http"

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
	redisDB = redisManager{
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

	hub := newHub()
	go hub.run()

	mux := http.NewServeMux()
	mux.HandleFunc("/sign_in", signInHandler)
	mux.HandleFunc("/sign_up", signUpHandler)
	mux.HandleFunc("/sign_out", SignOutHandler)
	mux.HandleFunc("/chat", hub.handleClientConn)

	err = http.ListenAndServe(":8090", NewAuthMiddlewareHandler(mux))

	if err != nil {
		log.Println(err)
	}
}
