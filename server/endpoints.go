package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

func signUpHandler(w http.ResponseWriter, r *http.Request) {

	fmt.Println("sign up")

	data, err := io.ReadAll(r.Body)

	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var message struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	json.Unmarshal(data, &message)

	userHandler := dbManager.initializeDBhandler("user")

	id, err := userHandler.CreateUserHandler(message.Username, message.Email, message.Password)
	fmt.Println(id)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	s, err := generateToken(id)

	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	setAuthCookies(w, s, message.Username)

	w.WriteHeader(http.StatusOK)
}

func signInHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("sign in")
	userHandler := dbManager.initializeDBhandler("user")
	var message struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	data, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err = json.Unmarshal(data, &message)

	if err != nil {
		fmt.Println(err)
	}

	id, err := userHandler.LoginUserHandler(message.Username, message.Password)
	fmt.Println(id)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(fmt.Sprintf(`"%s"`, err))
		return
	} else {
		token, _ := generateToken(id)
		setAuthCookies(w, token, message.Username)
		w.WriteHeader(http.StatusOK)
	}
}

func SignOutHandler(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Expires:  time.Now(),
		HttpOnly: true,
		Path:     "/",
		SameSite: http.SameSiteStrictMode,
	})
	http.SetCookie(w, &http.Cookie{
		Name:     "username",
		Expires:  time.Now(),
		Path:     "/",
		SameSite: http.SameSiteStrictMode,
	})

	w.WriteHeader(http.StatusOK)
}

func chatHandler(w http.ResponseWriter, r *http.Request) {

	newClient := initializeWSconn(w, r)
	fmt.Println(newClient)
	newClient.sendSubscribedChats()
	newClient.sendMessageHistory()
	chatList.addClientToSubRooms(newClient)
	connSockets.AddHubMember(newClient)
	newClient.socket.SetCloseHandler(func(code int, text string) error {
		newClient.CloseConnection()
		return nil
	})

	fmt.Println(connSockets)

	go clientMessages(newClient)
}
