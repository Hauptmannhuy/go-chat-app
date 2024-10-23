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
		Body   string
		ChatID string
		UserID string
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

func (h *Handler) GetAllChats() ([]string, error) {
	return h.ChatService.GetAllChats()
}

func (h *Handler) CreateChatHandler(j []byte) error {
	var chatObj struct {
		ID string
	}
	err := json.Unmarshal(j, &chatObj)

	if err != nil {
		fmt.Println(err)
	}

	if chatObj.ID == "" {
		return &argError{"chatID field cannot be blank"}
	}

	err = h.ChatService.CreateChat(chatObj.ID)

	if err != nil {
		fmt.Println(err, "Failed to create message")
		return err
	}
	return err
}

func (h *Handler) CreateUserHandler(username, email, pass string) error {

	if username == "" || pass == "" || email == "" {
		var err = &argError{"Username, email or password fields"}
		return err
	}
	return h.UserService.CreateAccount(username, email, pass)

}

func (h *Handler) LoginUserHandler(username, pass string) error {
	if username == "" || pass == "" {
		var err = &argError{"Username or password fields"}
		return err
	}
	return h.UserService.LoginUser(username, pass)
}

func (h *Handler) LoadUserSubscriptionsHandler(username string) ([]string, error) {
	return h.SubscriptionService.LoadSubscriptions(username)
}

func (h *Handler) SaveSubHandler(j []byte) error {
	var data struct {
		UserID string `json:"userID"`
		ChatID string `json:"chatID"`
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
