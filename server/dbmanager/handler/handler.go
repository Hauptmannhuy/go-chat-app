package handler

import (
	"encoding/json"
	"fmt"
	"go-chat-app/dbmanager/service"
)

type MessageHandler struct {
	MessageService service.MessageService
}

func (h *MessageHandler) CreateMessageHandler(j []byte) {
	var message struct {
		Body   string
		ChatID string
		// IsGroup bool	`json: "isGroup"`
	}

	if err := json.Unmarshal(j, &message); err != nil {
		fmt.Println(err)
	}

	if message.Body == "" || message.ChatID == "" {
		fmt.Println("Missing message body or chatID")
		return
	}

	if err := h.MessageService.CreateMessage(message.Body, message.ChatID); err != nil {
		fmt.Println("Failed to create message")
		return
	}

}
