//go:build gore
// +build gore

package main

import (
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/create-user/{id}", CreateUserHandler)
	// /create-user/1
	// {"name":"John","age":30,"email":"test123@mail.com","password":"password"}
	mux.HandleFunc("/create-user-2", CreateUserHandler2)
	// /create-user-2?id=1&user={"name":"John","age":30,"email":"test123@mail.com","password":"password"}
	mux.HandleFunc("/get-user/{id}", GetUserHandler)
	// /get-user/1
	mux.HandleFunc("/get-user-2", GetUserHandler2)
	// /get-user-2?id=1
	log.Println("Listening on :8080")

	http.ListenAndServe(":8080", mux)
}
