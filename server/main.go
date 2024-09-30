package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type ApiResponse struct {
	Message string `json:"message"`
}

func main() {
	http.HandleFunc("/", homeHandler)

	err := http.ListenAndServe(":8090", nil)
	if err != nil {
		fmt.Println(err)
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "OPTIONS" {
		setOptions(w)
	} else if r.Method == "GET" {
		getHome(w, r)
	}

}

func setOptions(w http.ResponseWriter) {
	setCorsHeaders(w)
	w.WriteHeader(http.StatusOK)
}

func getHome(w http.ResponseWriter, r *http.Request) {

	setCorsHeaders(w)

	if err := r; err != nil {
		fmt.Println(err)
	}

	response := ApiResponse{"Hello from backend!"}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func setCorsHeaders(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
}
