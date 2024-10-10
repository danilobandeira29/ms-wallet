package entity

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"time"
)

type Transaction struct {
	ID          string
	AccountFrom *Account
	AccountTo   *Account
	Amount      int64
	CreatedAt   time.Time
}

func NewTransaction(accountFrom *Account, accountTo *Account, amount int64) (*Transaction, error) {
	if amount <= 0 {
		return nil, errors.New("amount have to be greater than zero")
	}
	if err := accountFrom.Debit(amount); err != nil {
		return nil, err
	}
	accountTo.Deposit(amount)
	return &Transaction{
		ID:          uuid.New().String(),
		AccountFrom: accountFrom,
		AccountTo:   accountTo,
		Amount:      amount,
		CreatedAt:   time.Now(),
	}, nil
}

func div(a int, b int) (int, error) {
	if b == 0 {
		return 0, errors.New("division zero cannot be done")
	}
	return a / b, nil
}

func main() {
	resultado, err := div(1, 0)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Resultado %d", resultado)
}
