package handler

import (
	"fmt"
	"go-chat-app/dbmanager/store"
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

func (h *Handler) GetAllChats() (store.Chats, error) {
	return h.ChatService.GetAllChats()
}

func (h *Handler) CreatePrivateChatHandler(initiatorID, receiverID string) (interface{}, error) {

	return h.ChatService.CreatePrivateChat(initiatorID, receiverID)
}

func (h *Handler) LoadUserSubscribedChats(id string) ([]interface{}, error) {
	return h.ChatService.LoadSubscribedChats(id)
}

func (h *Handler) LoadSubscribedPrivateChats(id string) (interface{}, error) {
	return h.ChatService.LoadSubscribedPrivateChats(id)
}

func (h *Handler) SearchChat(input, userID string) (interface{}, error) {
	if input == "" {
		return nil, &argError{"Input should not be empty"}
	}
	return h.ChatService.SearchChat(input, userID)
}

func (h *Handler) RetrieveGroupChatCreatorID(chatID string) string {
	return h.ChatService.RetrieveGroupChatCreatorID(chatID)
}
