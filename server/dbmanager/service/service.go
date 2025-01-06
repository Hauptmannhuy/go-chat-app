package service

import "go-chat-app/dbmanager/store"

type Service struct {
	MessageStore      store.MessageStore
	ChatStore         store.ChatStore
	UserStore         store.UserStore
	SubscriptionStore store.SubscriptionStore
}

func (s *Service) SearchUser(input, userID string) (map[string]store.UserContainerData, error) {
	return s.UserStore.SearchUser(input, userID)
}

func (s *Service) CreateAccount(name, email, pass string) (string, error) {
	return s.UserStore.SaveAccount(name, email, pass)
}

func (s *Service) LoginUser(name, pass string) (string, error) {
	return s.UserStore.AuthenticateAccount(name, pass)
}

func (s *Service) LoadSubscriptions(username string) ([]string, error) {
	return s.SubscriptionStore.LoadSubscriptions(username)
}

func (s *Service) SaveSubscription(userID, chatID string) error {
	return s.SubscriptionStore.SaveSubscription(userID, chatID)
}

func (s *Service) GetPrivateChatSubs(chatName, sender string) []string {
	return s.SubscriptionStore.GetPrivateChatSubs(chatName, sender)
}

func (s *Service) RetrieveChatsMessages(subs []string) (interface{}, error) {
	return s.MessageStore.GetChatsMessages(subs)
}

func (s *Service) CreateMessage(body, chatID, userID string) (int, error) {
	return s.MessageStore.SaveMessage(body, chatID, userID)
}

func (s *Service) SearchChat(input, userID string) (interface{}, error) {
	return s.ChatStore.SearchChat(input, userID)
}

func (s *Service) GetAllChats() (store.Chats, error) {
	return s.ChatStore.GetChats()
}

func (s *Service) CreateChat(name, creatorID string) (string, error) {
	return s.ChatStore.SaveChat(name, creatorID)
}

func (s *Service) LoadSubscribedChats(id string) ([]interface{}, error) {
	return s.ChatStore.LoadSubscribedChats(id)
}

func (s *Service) CreatePrivateChat(user1id, user2id string) (interface{}, error) {
	return s.ChatStore.SavePrivateChat(user1id, user2id)
}

func (s *Service) LoadSubscribedPrivateChats(id string) (interface{}, error) {
	return s.ChatStore.LoadSubscribedPrivateChats(id)
}

func (s *Service) RetrieveGroupChatCreatorID(chatID string) string {
	return s.ChatStore.RetrieveGroupChatCreatorID(chatID)
}
