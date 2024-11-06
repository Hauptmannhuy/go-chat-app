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

func (h *Handler) CreateMessageHandler(body, chatID, userID string) error {

	if body == "" || chatID == "" {
		var err = &argError{"Message body or chatID"}

		return err
	}

	if err := h.MessageService.CreateMessage(body, chatID, userID); err != nil {
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

func (h *Handler) SaveSubHandler(userID, chatID string) error {

	return h.SubscriptionService.SaveSubscription(userID, chatID)

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
