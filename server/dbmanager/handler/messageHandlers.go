package handler

import "fmt"

func (h *Handler) CreateMessageHandler(body, chatID, userID string) (int, error) {

	if body == "" || chatID == "" {
		var err = &argError{"Message body or chatID"}

		return -1, err
	}

	messageID, err := h.MessageService.CreateMessage(body, chatID, userID)
	if err != nil {
		fmt.Println(err, "Failed to create message")
		return -1, err
	}
	return messageID, nil
}

func (h *Handler) GetChatsMessages(subs []string) (interface{}, error) {
	return h.MessageService.RetrieveChatsMessages(subs)
}
