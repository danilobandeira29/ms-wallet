package createaccount

import (
	"github.com/danilobandeira29/ms-wallet/internal/entity"
	"github.com/danilobandeira29/ms-wallet/internal/gateway"
)

type InputDTO struct {
	ClientID string `json:"client_id"`
}

type OutputDTO struct {
	ID string
}

type UseCase struct {
	AccountGateway gateway.AccountGateway
	ClientGateway  gateway.ClientGateway
}

func NewCreateAccountUseCase(ag gateway.AccountGateway, c gateway.ClientGateway) (*UseCase, error) {
	return &UseCase{
		AccountGateway: ag,
		ClientGateway:  c,
	}, nil
}

func (u *UseCase) Execute(input InputDTO) (*OutputDTO, error) {
	client, err := u.ClientGateway.Get(input.ClientID)
	if err != nil {
		return nil, err
	}
	account, err := entity.NewAccount(client, 100)
	if err != nil {
		return nil, err
	}
	if err = u.AccountGateway.Save(account); err != nil {
		return nil, err
	}
	return &OutputDTO{
		ID: account.ID,
	}, nil
}
