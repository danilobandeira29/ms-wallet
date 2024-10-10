package gateway

import "github.com/danilobandeira29/ms-wallet/internal/entity"

type TransactionGatewayInterface interface {
	Create(t *entity.Transaction) error
}
