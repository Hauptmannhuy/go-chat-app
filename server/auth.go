package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
)

func setAuthCookies(w http.ResponseWriter, s string, name string) {

	http.SetCookie(w, &http.Cookie{
		Name:     "username",
		Value:    name,
		Expires:  time.Now().AddDate(0, 0, 14),
		Path:     "/",
		SameSite: http.SameSiteStrictMode,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    s,
		Expires:  time.Now().AddDate(0, 0, 14),
		HttpOnly: true,
		Path:     "/",
		SameSite: http.SameSiteStrictMode,
	})
}

func generateToken() (string, error) {
	key := []byte(os.Getenv("KEY"))
	t := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"iss": "go-chat-app",
		})
	s, err := t.SignedString(key)

	if err != nil {
		return "", err
	}
	return s, nil
}

func verifyToken(c *http.Cookie) bool {
	tokenS := c.Value

	token, err := jwt.Parse(tokenS, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("KEY")), nil
	})

	if err != nil {
		log.Fatal(err)
		return false
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		fmt.Println(claims["iss"])
	} else {
		fmt.Println(err)
	}
	return true
}
