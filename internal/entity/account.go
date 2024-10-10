package entity

import (
	"errors"
	"github.com/google/uuid"
	"time"
)

type Account struct {
	ID             string
	Client         *Client
	BalanceInCents int64
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

func NewAccount(client *Client, balanceInCents int64) (*Account, error) {
	if client == nil {
		return &Account{}, errors.New("client is mandatory")
	}
	return &Account{
		ID:             uuid.New().String(),
		Client:         client,
		BalanceInCents: balanceInCents,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}, nil
}

func (a *Account) Deposit(amount int64) {
	a.BalanceInCents += amount
}

func (a *Account) Debit(amount int64) error {
	if a.BalanceInCents < amount {
		return errors.New("balance insufficient")
	}
	a.BalanceInCents -= amount
	return nil
}
