package service

import "go-chat-app/dbmanager/store"

type MessageService struct {
	MessageStore store.MessageStore
}

func (s *MessageService) CreateMessage(body, chatID string) error {
	return s.MessageStore.SaveMessage(body, chatID)
}

func (s *MessageService) ListMessages() ([]store.Message, error) {
	return s.MessageStore.GetAllMessages()
}
