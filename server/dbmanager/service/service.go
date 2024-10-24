package service

import "go-chat-app/dbmanager/store"

type Service struct {
	MessageStore      store.MessageStore
	ChatStore         store.ChatStore
	UserStore         store.UserStore
	SubscriptionStore store.SubscriptionStore
}

func (s *Service) CreateMessage(body, chatID, userID string) error {
	return s.MessageStore.SaveMessage(body, chatID, userID)
}

func (s *Service) GetAllChats() ([]string, error) {
	return s.ChatStore.GetChats()
}

func (s *Service) CreateChat(chatID string) error {
	return s.ChatStore.SaveChat(chatID)
}

func (s *Service) CreateAccount(name, email, pass string) error {
	return s.UserStore.SaveAccount(name, email, pass)
}

func (s *Service) LoginUser(name, pass string) error {
	return s.UserStore.AuthenticateAccount(name, pass)
}

func (s *Service) LoadSubscriptions(username string) ([]string, error) {
	return s.SubscriptionStore.LoadSubscriptions(username)
}

func (s *Service) SaveSubscription(username, chatID string) error {
	return s.SubscriptionStore.SaveSubscription(username, chatID)
}

func (s *Service) RetrieveChatsMessages(subs []string) (interface{}, error) {
	return s.MessageStore.GetChatsMessages(subs)
}

func (s *Service) SearchQuery(input string) ([]string, error) {
	return s.ChatStore.SearchQuery(input)
}
