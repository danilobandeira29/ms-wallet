package entity

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestTransaction(t *testing.T) {
	clientFrom, _ := NewClient("Client from", "client@from.com")
	clientTo, _ := NewClient("Client to", "client@to.com")
	accountFrom, _ := NewAccount(clientFrom, 100)
	accountTo, _ := NewAccount(clientTo, 50)
	_, err := NewTransaction(accountFrom, accountTo, 50)
	assert.NoError(t, err)
	assert.Equal(t, int64(100), accountTo.BalanceInCents)
	assert.Equal(t, int64(50), accountFrom.BalanceInCents)
	_, err = NewTransaction(accountFrom, accountTo, 10)
	assert.NoError(t, err)
	assert.Equal(t, int64(110), accountTo.BalanceInCents)
	assert.Equal(t, int64(40), accountFrom.BalanceInCents)
}

func TestTransactionInsufficientBalance(t *testing.T) {
	clientFrom, _ := NewClient("Client from", "client@from.com")
	clientTo, _ := NewClient("Client to", "client@to.com")
	accountFrom, _ := NewAccount(clientFrom, 100)
	accountTo, _ := NewAccount(clientTo, 100)
	transaction, err := NewTransaction(accountFrom, accountTo, 200)
	assert.NotNil(t, err)
	assert.Nil(t, transaction)
	assert.True(t, strings.Contains(err.Error(), "insufficient"))
	assert.Equal(t, int64(100), accountTo.BalanceInCents)
	assert.Equal(t, int64(100), accountFrom.BalanceInCents)
}
