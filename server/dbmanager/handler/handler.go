package handler

import (
	"encoding/json"
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

func (h *Handler) CreateMessageHandler(j []byte) error {
	var message struct {
		Body   string `json:"body"`
		ChatID string `json:"chat_id"`
		UserID string `json:"user_id"`
	}

	if err := json.Unmarshal(j, &message); err != nil {
		fmt.Println(err)
		return err
	}

	if message.Body == "" || message.ChatID == "" {
		var err = &argError{"Message body or chatID"}

		return err
	}

	if err := h.MessageService.CreateMessage(message.Body, message.ChatID, message.UserID); err != nil {
		fmt.Println(err, "Failed to create message")
		return err
	}
	return nil
}

func (h *Handler) CreateUserHandler(username, email, pass string) (string, error) {

	if username == "" || pass == "" || email == "" {
		var err = &argError{"Username, email or password fields"}
		return "", err
	}
	return h.UserService.CreateAccount(username, email, pass)

}

func (h *Handler) LoginUserHandler(username, pass string) (string, error) {
	if username == "" || pass == "" {
		var err = &argError{"Username or password fields"}
		return "", err
	}
	return h.UserService.LoginUser(username, pass)
}

func (h *Handler) LoadSubscriptions(username string) ([]string, error) {
	return h.SubscriptionService.LoadSubscriptions(username)
}

func (h *Handler) SaveSubHandler(j []byte) error {
	var data struct {
		UserID string `json:"user_id"`
		ChatID string `json:"chat_id"`
	}
	err := json.Unmarshal(j, &data)

	if err != nil {
		fmt.Println(err)
		return err
	}
	return h.SubscriptionService.SaveSubscription(data.UserID, data.ChatID)

}

func (h *Handler) GetChatsMessages(subs []string) (interface{}, error) {
	return h.MessageService.RetrieveChatsMessages(subs)
}

func (h *Handler) SearchChat(input, userID string) ([]interface{}, error) {
	if input == "" {
		return nil, &argError{"Input should not be empty"}
	}
	return h.ChatService.SearchChat(input, userID)
}

func (h *Handler) LoadUserSubscribedChats(username string) ([]interface{}, error) {
	return h.ChatService.LoadSubscribedChats(username)
}
