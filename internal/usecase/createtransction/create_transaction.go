package createtransction

import (
	"context"
	"github.com/danilobandeira29/ms-wallet/internal/entity"
	"github.com/danilobandeira29/ms-wallet/internal/gateway"
	"github.com/danilobandeira29/ms-wallet/pkg/events"
	"github.com/danilobandeira29/ms-wallet/pkg/uow"
)

type Input struct {
	AccountIDFrom string `json:"account_id_from"`
	AccountIDTo   string `json:"account_id_to"`
	Amount        int64  `json:"amount"`
}

type Output struct {
	ID            string `json:"id"`
	AccountIDFrom string `json:"account_id_from"`
	AccountIDTo   string `json:"account_id_to"`
	Amount        int64  `json:"amount"`
}

type BalanceUpdatedDto struct {
	AccountIDFrom        string `json:"account_id_from"`
	AccountIDTo          string `json:"account_id_to"`
	BalanceAccountIDFrom int64  `json:"balance_account_id_from"`
	BalanceAccountIDTo   int64  `json:"balance_account_id_to"`
}

type UseCase struct {
	Uow                uow.UowInterface
	EventDispatcher    events.EventDispatcherInterface
	TransactionCreated events.EventInterface
	BalanceUpdated     events.EventInterface
}

func NewCreateTransactionUseCase(uow uow.UowInterface, ed events.EventDispatcherInterface, e events.EventInterface, bu events.EventInterface) (*UseCase, error) {
	return &UseCase{
		Uow:                uow,
		EventDispatcher:    ed,
		TransactionCreated: e,
		BalanceUpdated:     bu,
	}, nil
}

func (u *UseCase) Execute(ctx context.Context, input Input) (*Output, error) {
	output := &Output{}
	var balanceUpdatedDto BalanceUpdatedDto
	err := u.Uow.Do(ctx, func(_ *uow.Uow) error {
		accountGateway := u.getAccountGateway(ctx)
		transactionGateway := u.getTransactionGateway(ctx)
		accountFrom, err := accountGateway.FindBy(input.AccountIDFrom)
		if err != nil {
			return err
		}
		accountTo, err := accountGateway.FindBy(input.AccountIDTo)
		if err != nil {
			return err
		}
		transaction, err := entity.NewTransaction(accountFrom, accountTo, input.Amount)
		if err != nil {
			return err
		}
		err = accountGateway.UpdateBalance(accountFrom)
		if err != nil {
			return err
		}
		err = accountGateway.UpdateBalance(accountTo)
		if err != nil {
			return err
		}
		err = transactionGateway.Create(transaction)
		if err != nil {
			return err
		}
		output.ID = transaction.ID
		output.AccountIDFrom = transaction.AccountFrom.ID
		output.AccountIDTo = transaction.AccountTo.ID
		output.Amount = transaction.Amount
		balanceUpdatedDto.AccountIDFrom = transaction.AccountFrom.ID
		balanceUpdatedDto.AccountIDTo = transaction.AccountTo.ID
		balanceUpdatedDto.BalanceAccountIDFrom = accountFrom.BalanceInCents
		balanceUpdatedDto.BalanceAccountIDTo = accountTo.BalanceInCents
		return nil
	})
	if err != nil {
		return nil, err
	}
	u.TransactionCreated.SetPayload(output)
	err = u.EventDispatcher.Dispatch(u.TransactionCreated)
	if err != nil {
		return nil, err
	}
	u.BalanceUpdated.SetPayload(balanceUpdatedDto)
	err = u.EventDispatcher.Dispatch(u.BalanceUpdated)
	if err != nil {
		return nil, err
	}
	return output, nil
}

func (u *UseCase) getAccountGateway(ctx context.Context) gateway.AccountGateway {
	repo, err := u.Uow.GetRepository(ctx, "AccountDB")
	if err != nil {
		panic(err)
	}
	return repo.(gateway.AccountGateway)
}

func (u *UseCase) getTransactionGateway(ctx context.Context) gateway.TransactionGatewayInterface {
	repo, err := u.Uow.GetRepository(ctx, "TransactionDB")
	if err != nil {
		panic(err)
	}
	return repo.(gateway.TransactionGatewayInterface)
}
