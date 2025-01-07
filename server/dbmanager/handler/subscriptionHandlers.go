package handler

func (h *Handler) LoadSubscriptions(username string) ([]string, error) {
	return h.SubscriptionService.LoadSubscriptions(username)
}

func (h *Handler) SaveSubHandler(userID, chatID string) error {

	return h.SubscriptionService.SaveSubscription(userID, chatID)

}

func (h *Handler) GetPrivateChatSubs(chatName, sender string) []string {
	return h.SubscriptionService.GetPrivateChatSubs(chatName, sender)
}

func (h *Handler) GetGroupChatSubs(chatName, sender string) []string {
	return h.SubscriptionService.GetGroupChatSubs(chatName, sender)
}
