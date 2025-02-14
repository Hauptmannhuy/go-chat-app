package handler

import (
	"fmt"
	"go-chat-app/dbmanager/store"
)

func (h *Handler) CreateChatHandler(name string, creatorID int) (int, error) {

	if creatorID == 0 {
		return 0, &argError{"creatorID field cannot be blank"}
	}

	n, err := h.ChatService.CreateChat(name, creatorID)

	if err != nil {
		fmt.Println(err, "Failed to create message")
		return 0, err
	}
	return n, err
}

func (h *Handler) GetAllChats() (store.Chats, error) {
	return h.ChatService.GetAllChats()
}

func (h *Handler) CreatePrivateChatHandler(initiatorID, receiverID int) (interface{}, error) {

	return h.ChatService.CreatePrivateChat(initiatorID, receiverID)
}

func (h *Handler) LoadUserSubscribedChats(id int) ([]interface{}, error) {
	return h.ChatService.LoadSubscribedChats(id)
}

func (h *Handler) LoadSubscribedPrivateChats(id int) (interface{}, error) {
	return h.ChatService.LoadSubscribedPrivateChats(id)
}

func (h *Handler) SearchChat(input string, userID int) (interface{}, error) {
	if input == "" {
		return nil, &argError{"Input should not be empty"}
	}
	return h.ChatService.SearchChat(input, userID)
}

func (h *Handler) RetrieveGroupChatCreatorID(chatID int) int {
	return h.ChatService.RetrieveGroupChatCreatorID(chatID)
}
