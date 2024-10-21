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

type sqlDBwrap struct {
	db *sql.DB
}

var dbManager sqlDBwrap

func (dbm *sqlDBwrap) openAndMigrateDB() error {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
		return err
	}
	dotenv := os.Getenv("DATABASE_CREDS")
	fmt.Println(dotenv)
	dataSourceName := fmt.Sprintf("postgres://%s/dbmanager?sslmode=disable", dotenv)

	dbm.db, err = sql.Open("postgres", dataSourceName)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
		return err
	}
	fmt.Println("Connected to database")

	driver, err := postgres.WithInstance(dbm.db, &postgres.Config{})
	if err != nil {
		log.Fatalf("Could not start SQL driver: %v", err)
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"postgres", driver)
	if err != nil {
		log.Fatalf("Could not start migration: %v", err)
		return err
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Migration failed: %v", err)
		return err
	}
	return nil
}

func (dbm *sqlDBwrap) handleDataBase(env OutEnvelope) error {
	jsoned, _ := json.Marshal(env.Data)
	fmt.Println(env.Data)
	switch env.Type {
	case "NEW_MESSAGE":
		messageHandler := dbManager.initializeDBhandler("message")
		err := messageHandler.CreateMessageHandler(jsoned)
		return err
	case "NEW_CHAT":
		chatHandler := dbManager.initializeDBhandler("chat")
		err := chatHandler.CreateChatHandler(jsoned)

		return err
	case "JOIN_CHAT":
		subHandler := dbManager.initializeDBhandler("subscription")
		subHandler.SaveSubHandler(jsoned)
		return nil
	default:
		fmt.Println("No write to database")
		return nil
	}
}

func (dbm *sqlDBwrap) initializeDBhandler(handlerDeclaration string) handler.Handler {
	dbStore := store.SQLstore{DB: dbManager.db}
	service := service.Service{}
	handler := handler.Handler{}
	switch handlerDeclaration {
	case "user":
		service.UserStore = &dbStore
		handler.UserService = service
	case "chat":
		service.ChatStore = &dbStore
		handler.ChatService = service
	case "message":
		service.MessageStore = &dbStore
		handler.MessageService = service
	case "subscription":
		service.SubscriptionStore = &dbStore
		handler.SubscriptionService = service
	}
	return handler
}
