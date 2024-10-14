package service

import "go-chat-app/dbmanager/store"

type Service struct {
	MessageStore store.MessageStore
	ChatStore    store.ChatStore
	UserStore    store.UserStore
}

func (s *Service) CreateMessage(body, chatID string) error {
	return s.MessageStore.SaveMessage(body, chatID)
}

func (s *Service) ListMessages() ([]store.Message, error) {
	return s.MessageStore.GetAllMessages()
}

func (s *Service) ListChats() {

}

func (s *Service) CreateChat(chatID string) error {
	return s.ChatStore.SaveChat(chatID)
}

func (s *Service) CreateAccount(name, email, pass string) error {
	return s.UserStore.SaveAccount(name, email, pass)
}
