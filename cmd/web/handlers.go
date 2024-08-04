//go:build gore
// +build gore

package main

import (
	"fmt"

	"github.com/ic-it/goreapi"
)

type User struct {
	Id       int    `json:"id,omitempty"`
	Name     string `json:"name"`
	Age      int    `json:"age"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

var (
	CreateUserHandler = goreapi.H(
		goreapi.Param("id", &goreapi.Int{}),
		goreapi.Body("user", &User{}),
		goreapi.Response(&goreapi.Int{}),
		func(id int, user *User) int {
			fmt.Printf("Creating user: %s\n", user.Name)
			return id
		},
	)

	CreateUserHandler2 = goreapi.H(
		goreapi.Query("id", &goreapi.Int{}),
		goreapi.Query("user", &User{}),
		goreapi.Response(&goreapi.Int{}),
		func(id int, user *User) int {
			fmt.Printf("Creating user: %s\n", user.Name)
			return id
		},
	)

	GetUserHandler = goreapi.H(
		goreapi.Param("id", &goreapi.Int{}),
		goreapi.Response(&User{}),
		func(id int) *User {
			return &User{
				Id:       id,
				Name:     "John",
				Age:      30,
				Email:    "test@mail.com",
				Password: "password",
			}
		},
	)

	GetUserHandler2 = goreapi.H(
		goreapi.Query("id", &goreapi.Int{}),
		goreapi.Response(&User{}),
		func(id int) *User {
			return &User{
				Id:       id,
				Name:     "John",
				Age:      30,
				Email:    "asdasd@asdasd.com",
				Password: "password",
			}
		},
	)
)
