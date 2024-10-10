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
	MessageService service.Service
	ChatService    service.Service
}

func (h *Handler) CreateMessageHandler(j []byte) error {
	var message struct {
		Body   string
		ChatID string
		// IsGroup bool	`json: "isGroup"`
	}

	if err := json.Unmarshal(j, &message); err != nil {
		fmt.Println(err)
		return err
	}

	if message.Body == "" || message.ChatID == "" {
		var err = &argError{"Message body or chatID"}

		return err
	}

	if err := h.MessageService.CreateMessage(message.Body, message.ChatID); err != nil {
		fmt.Println(err, "Failed to create message")
		return err
	}
	return nil
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
