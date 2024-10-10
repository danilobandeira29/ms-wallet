package createclient

import (
	"github.com/danilobandeira29/ms-wallet/internal/entity"
	"github.com/danilobandeira29/ms-wallet/internal/gateway"
	"time"
)

type InputDTO struct {
	Name  string
	Email string
}

type OutputDTO struct {
	ID        string
	Name      string
	Email     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type CreateClientUseCase struct {
	ClientGateway gateway.ClientGateway
}

func NewCreateClientUseCase(cg gateway.ClientGateway) *CreateClientUseCase {
	return &CreateClientUseCase{
		ClientGateway: cg,
	}
}

func (u *CreateClientUseCase) Execute(input InputDTO) (*OutputDTO, error) {
	client, err := entity.NewClient(input.Name, input.Email)
	if err != nil {
		return nil, err
	}
	err = u.ClientGateway.Save(client)
	if err != nil {
		return nil, err
	}
	return &OutputDTO{
		ID:        client.ID,
		Name:      client.Name,
		Email:     client.Email,
		CreatedAt: client.CreatedAt,
		UpdatedAt: client.UpdatedAt,
	}, nil
}
