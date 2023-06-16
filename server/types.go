package main

import (
	"time"

	"github.com/google/uuid"
)

type CreateAccountRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}

type Account struct {
	ID        uuid.UUID `json:"accountId"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"createdAt"`
}

func NewAccount(username, email string) *Account {
	return &Account{
		Username:  username,
		Email:     email,
		CreatedAt: time.Now().UTC(),
	}
}
