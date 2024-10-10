package gateway

import (
	"github.com/danilobandeira29/ms-wallet/internal/entity"
)

type AccountGateway interface {
	FindBy(id string) (*entity.Account, error)
	Save(a *entity.Account) error
	UpdateBalance(a *entity.Account) error
}
