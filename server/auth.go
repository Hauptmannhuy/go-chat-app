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

func generateToken(id string) (string, error) {
	key := []byte(os.Getenv("KEY"))
	t := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"iss": "go-chat-app",
			"id":  id,
		})
	s, err := t.SignedString(key)

	if err != nil {
		return "", err
	}
	return s, nil
}

func parseToken(tokenS string) (*jwt.Token, bool) {

	token, err := jwt.Parse(tokenS, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("KEY")), nil
	})
	if err != nil {
		log.Fatal(err)
		return nil, false
	}
	return token, true
}

func verifyToken(c *http.Cookie) bool {
	tokenS := c.Value

	token, ok := parseToken(tokenS)
	if !ok {
		return false
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		fmt.Println(claims["id"])
	}
	return true
}

func fetchUserID(tokenS string) string {
	token, _ := parseToken(tokenS)
	claims, _ := token.Claims.(jwt.MapClaims)
	val := claims["id"].(string)
	return val
}
