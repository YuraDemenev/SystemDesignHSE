package main

import (
	authservice "auth-service/handlers"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/register", authservice.RegisterHandler)
	http.HandleFunc("/login", authservice.LoginHandler)
	log.Println("Auth service running on port 8000")
	log.Println(http.ListenAndServe(":8000", nil))
}
