package handler

import (
	"fmt"
	"go-chat-app/dbmanager/service"
)

type argError struct {
	message string
}

func (e *argError) Error() string {
	return fmt.Sprintf("%s field cannot be blank", e.message)
}

type Handler struct {
	MessageService      service.Service
	ChatService         service.Service
	UserService         service.Service
	SubscriptionService service.Service
}
