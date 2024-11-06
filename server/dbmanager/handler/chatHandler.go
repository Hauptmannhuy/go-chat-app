package handler

import (
	"fmt"
)

func (h *Handler) CreateChatHandler(name, creatorID string) (string, error) {

	if creatorID == "" {
		return "", &argError{"chatID field cannot be blank"}
	}

	str, err := h.ChatService.CreateChat(name, creatorID)

	if err != nil {
		fmt.Println(err, "Failed to create message")
		return "", err
	}
	return str, err
}

func (h *Handler) GetAllChats() ([]string, error) {
	return h.ChatService.GetAllChats()
}

func (h *Handler) SearchUser(input, userID string) (interface{}, error) {

	if input == "" {
		return nil, &argError{"Input should not be empty"}
	}
	return h.UserService.SearchUser(input, userID)
}

func (h *Handler) CreatePrivateChatHandler(initiatorID, receiverID string) (string, error) {

	return h.ChatService.CreatePrivateChat(initiatorID, receiverID)
}

func (h *Handler) LoadUserSubscribedChats(id string) ([]interface{}, error) {
	return h.ChatService.LoadSubscribedChats(id)
}

func (h *Handler) LoadSubscribedPrivateChats(id string) (interface{}, error) {
	return h.ChatService.LoadSubscribedPrivateChats(id)
}
