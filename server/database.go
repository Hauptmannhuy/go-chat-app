package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"go-chat-app/dbmanager/handler"
	"go-chat-app/dbmanager/service"
	"go-chat-app/dbmanager/store"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/joho/godotenv"
)

var db *sql.DB

func openAndMigrateDB() (*sql.DB, error) {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
		return nil, err
	}
	dotenv := os.Getenv("DATABASE_CREDS")
	fmt.Println(dotenv)
	dataSourceName := fmt.Sprintf("postgres://%s/dbmanager?sslmode=disable", dotenv)

	db, err = sql.Open("postgres", dataSourceName)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
		return nil, err
	}
	fmt.Println("Connected to database")

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatalf("Could not start SQL driver: %v", err)
		return nil, err
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"postgres", driver)
	if err != nil {
		log.Fatalf("Could not start migration: %v", err)
		return nil, err
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Migration failed: %v", err)
		return nil, err
	}
	return db, nil
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
