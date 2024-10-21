package main

import (
	"fmt"
	"net/http"
)

type AuthorizationMiddleware struct {
	handler http.Handler
}

type AuthHandler struct{}

func (h AuthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

}

func NewAuthMiddlewareHandler(handler http.Handler) AuthorizationMiddleware {
	return AuthorizationMiddleware{
		handler: handler,
	}
}

func (am AuthorizationMiddleware) ServeHTTP(w http.ResponseWriter, req *http.Request) {

	if req.URL.Path == "/sign_in" {
		signInHandler(w, req)
		return
	}
	if req.URL.Path == "/sign_up" {
		signUpHandler(w, req)
		return
	}
	if req.URL.Path == "/sign_out" {
		SignOutHandler(w, req)
		return
	}

	cookie, err := req.Cookie("token")

	if err != nil {
		fmt.Println(err)
		fmt.Println("redirect")
		http.Redirect(w, req, "/sign_up", http.StatusSeeOther)
	} else {
		ok := verifyToken(cookie)
		if !ok {
			http.Redirect(w, req, "/sign_up", http.StatusSeeOther)
		} else {
			chatHandler(w, req)
		}
	}

}
