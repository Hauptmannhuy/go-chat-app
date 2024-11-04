package handler

import (
	"encoding/json"
	"fmt"
)

func (h *Handler) CreateChatHandler(j []byte) error {
	var chatObj struct {
		ID   string `json:"chat_id"`
		Type string `json:"chat_type"`
	}
	err := json.Unmarshal(j, &chatObj)

	if err != nil {
		fmt.Println(err)
	}

	if chatObj.ID == "" || chatObj.Type == "" {
		return &argError{"chatID field cannot be blank"}
	}

	err = h.ChatService.CreateChat(chatObj.ID, chatObj.Type)

	if err != nil {
		fmt.Println(err, "Failed to create message")
		return err
	}
	return err
}

func (h *Handler) GetAllChats() ([]string, error) {
	return h.ChatService.GetAllChats()
}

func (h *Handler) SearchUser(input string) (interface{}, error) {

	if input == "" {
		return nil, &argError{"Input should not be empty"}
	}
	return h.UserService.SearchUser(input)
}

func (h *Handler) CreatePrivateChatHandler(p []byte) (string, error) {
	var data struct {
		User1ID string `json:"initiator_id"`
		User2ID string `json:"receiver_id"`
	}
	err := json.Unmarshal(p, &data)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	return h.ChatService.CreatePrivateChat(data.User1ID, data.User2ID)
}
