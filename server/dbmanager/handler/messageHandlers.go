package handler

import "fmt"

func (h *Handler) CreateMessageHandler(body, chatName string, userID int) (int, error) {

	if body == "" || chatName == "" {
		var err = &argError{"Message body or chatID"}

		return -1, err
	}

	messageID, err := h.MessageService.CreateMessage(body, chatName, userID)
	if err != nil {
		fmt.Println(err, "Failed to create message")
		return -1, err
	}
	return messageID, nil
}

func (h *Handler) GetChatsMessages(subs []string) (interface{}, error) {
	return h.MessageService.RetrieveChatsMessages(subs)
}
