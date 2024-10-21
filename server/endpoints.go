package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
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

	userHandler := initializeDBhandler(db, "user")

	err = userHandler.CreateUserHandler(message.Username, message.Email, message.Password)

	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	s, err := generateToken()

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
	userHandler := initializeDBhandler(db, "user")
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

	err = userHandler.LoginUserHandler(message.Username, message.Password)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(fmt.Sprintf(`"%s"`, err))
		return
	} else {
		token, _ := generateToken()
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
	conn, err := upgrader.Upgrade(w, r, nil)
	var newClient = &Client{
		Socket:    conn,
		Connected: true,
	}

	connSockets.AddConection(newClient)

	newClient.Socket.SetCloseHandler(func(code int, text string) error {
		newClient.CloseConnection()
		return nil
	})

	if err != nil {
		log.Println(err)
		return
	}

	fmt.Println(connSockets)

	go clientMessages(newClient)
}
