package handler

func (h *Handler) LoadSubscriptions(userID int) ([]string, error) {
	return h.SubscriptionService.LoadSubscriptions(userID)
}

func (h *Handler) SaveSubHandler(userID, chatID int) error {

	return h.SubscriptionService.SaveSubscription(userID, chatID)

}

func (h *Handler) GetPrivateChatSubs(chatName string) []int {
	return h.SubscriptionService.GetPrivateChatSubs(chatName)
}

func (h *Handler) GetGroupChatSubs(chatName string) []int {
	return h.SubscriptionService.GetGroupChatSubs(chatName)
}
