package database

import (
	"database/sql"
	"github.com/danilobandeira29/ms-wallet/internal/entity"
)

type AccountDB struct {
	DB *sql.DB
}

func NewAccountDB(db *sql.DB) *AccountDB {
	return &AccountDB{
		DB: db,
	}
}

func (a *AccountDB) FindBy(id string) (*entity.Account, error) {
	var account entity.Account
	var client entity.Client
	account.Client = &client
	stmt, err := a.DB.Prepare("select a.id, a.client_id, a.balance, a.created_at, c.id, c.name, c.email, c.created_at from accounts a join clients c on c.id = a.client_id where a.id = ?")
	if err != nil {
		return nil, err
	}
	row := stmt.QueryRow(id)
	err = row.Scan(
		&account.ID,
		&account.Client.ID,
		&account.BalanceInCents,
		&account.CreatedAt,
		&client.ID,
		&client.Name,
		&client.Email,
		&client.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &account, nil
}

func (a *AccountDB) Save(account *entity.Account) error {
	stmt, err := a.DB.Prepare("insert into accounts(id, client_id, balance, created_at) values(?, ?, ?, ?)")
	defer func(s *sql.Stmt) {
		errClose := s.Close()
		if errClose != nil {
			err = errClose
		}
	}(stmt)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(account.ID, account.Client.ID, account.BalanceInCents, account.CreatedAt)
	if err != nil {
		return err
	}
	return nil
}

func (a *AccountDB) UpdateBalance(account *entity.Account) error {
	stmt, err := a.DB.Prepare("update accounts set balance = ? where id = ?")
	if err != nil {
		return err
	}
	defer func(stmt2 *sql.Stmt) {
		errClose := stmt2.Close()
		if err != nil {
			err = errClose
		}
	}(stmt)
	_, err = stmt.Exec(account.BalanceInCents, account.ID)
	if err != nil {
		return err
	}
	return nil
}
