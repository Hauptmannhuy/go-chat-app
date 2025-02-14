package handler

import "go-chat-app/dbmanager/store"

func (h *Handler) CreateUserHandler(username, email, pass string) (string, error) {

	if username == "" || pass == "" || email == "" {
		var err = &argError{"Username, email or password fields"}
		return "", err
	}
	return h.UserService.CreateAccount(username, email, pass)

}

func (h *Handler) LoginUserHandler(username, pass string) (string, error) {
	if username == "" || pass == "" {
		var err = &argError{"Username or password fields"}
		return "", err
	}
	return h.UserService.LoginUser(username, pass)
}

func (h *Handler) SearchUser(input string, userID int) (map[string]store.UserContainerData, error) {

	if input == "" {
		return nil, &argError{"Input should not be empty"}
	}
	return h.UserService.SearchUser(input, userID)
}
