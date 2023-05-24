package main

import "github.com/google/uuid"

type Account struct {
	ID       uuid.UUID `json:"account_id"`
	Username string    `json:"username"`
	Email    string    `json:"email"`
}

func NewAccount(username, email string) *Account {
	return &Account{
		ID:       uuid.New(),
		Username: username,
		Email:    email,
	}
}
