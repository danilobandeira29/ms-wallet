package entity

import (
	"errors"
	"github.com/google/uuid"
	"slices"
	"strings"
	"time"
)

type Client struct {
	ID        string
	Name      string
	Email     string
	Accounts  []*Account
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewClient(name string, email string) (*Client, error) {
	if name == "" {
		return &Client{}, errors.New("name is mandatory")
	}
	if email == "" {
		return &Client{}, errors.New("email is mandatory")
	}
	return &Client{
		ID:        uuid.New().String(),
		Name:      name,
		Email:     email,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}, nil
}

func (c *Client) AddAccounts(accounts []*Account) error {
	var accountsBelongToAnotherClient []string
	var newAccounts []*Account
	for _, acc := range accounts {
		if acc.Client.ID != c.ID {
			accountsBelongToAnotherClient = append(accountsBelongToAnotherClient, acc.ID)
		} else {
			newAccounts = append(newAccounts, acc)
		}
	}
	if len(accountsBelongToAnotherClient) > 0 {
		return errors.New("This Accounts belongs to another(s) Client(s):" + strings.Join(accountsBelongToAnotherClient, ", "))
	}
	c.Accounts = slices.Concat(c.Accounts, newAccounts)
	return nil
}
