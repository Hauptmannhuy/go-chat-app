package main

import (
	"fmt"
	"net/http"
)

type AuthorizationMiddleware struct {
	verifier tokenVerifier
	handler  http.Handler
}

type tokenVerifier func(c *http.Cookie) bool

func NewAuthMiddlewareHandler() AuthorizationMiddleware {
	mux := http.NewServeMux()
	mux.HandleFunc("/sign_in", signInHandler)
	mux.HandleFunc("/sign_up", signUpHandler)
	mux.HandleFunc("/sign_out", SignOutHandler)
	mux.HandleFunc("/chat", chatHandler)

	return AuthorizationMiddleware{
		verifier: verifyToken,
		handler:  mux,
	}
}

func (am AuthorizationMiddleware) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	cookie, err := req.Cookie("token")
	tokenValid := err == nil && am.verifier(cookie)
	path := req.URL.Path
	fmt.Println(tokenValid, path)
	switch {
	case tokenValid && path == "/checkauth":
		w.WriteHeader(http.StatusOK)
	case !tokenValid && (path == "/sign_in" || path == "/sign_up"):
		am.handler.ServeHTTP(w, req)
	case !tokenValid:
		fmt.Println("redirect")
		http.Redirect(w, req, "/sign_up", http.StatusUnauthorized)
	case tokenValid && (path == "/sign_up" || path == "/sign_in"):
		http.Redirect(w, req, "/chat", http.StatusSeeOther)
	case tokenValid && path == "/sign_out":
		fmt.Println("sign out")
		am.handler.ServeHTTP(w, req)
	case tokenValid && path == "/chat":
		am.handler.ServeHTTP(w, req)
	default:
		http.NotFound(w, req)
	}
}
