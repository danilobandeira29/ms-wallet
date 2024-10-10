package database

import (
	"database/sql"
	"github.com/danilobandeira29/ms-wallet/internal/entity"
	"github.com/stretchr/testify/suite"
	"testing"
)

type TransactionDBTest struct {
	suite.Suite
	db            *sql.DB
	transactionDB *TransactionDB
}

func (t *TransactionDBTest) SetupSuite() {
	db, err := sql.Open("sqlite3", ":memory:")
	t.Nil(err)
	t.db = db
	t.db.Exec("create table clients(id varchar(255) primary key, name varchar(255) not null, email varchar(255) not null, created_at date, updated_at date)")
	t.db.Exec(`
	create table accounts(
	    id varchar(255) primary key, 
	    balance int not null, 
	    client_id varchar(255),
	    created_at date not null,
	    updated_at date,
	    foreign key (client_id) references clients(id) 
	)
`)
	t.db.Exec(`
	create table transactions(
	    id varchar(255) primary key,
	    account_from_id varchar(255), 
	    account_to_id varchar(255), 
	    amount int not null, 
	    created_at date,
	    foreign key (account_from_id) references accounts(id),
	    foreign key (account_to_id) references accounts(id)
	                         )`)
	t.transactionDB = NewTransactionDB(t.db)
}

func (t *TransactionDBTest) TearDownSuite() {
	defer t.db.Close()
	t.db.Exec("drop table clients")
	t.db.Exec("drop table accounts")
	t.db.Exec("drop table transactions")
}

func TestTransactionDBSuite(t *testing.T) {
	suite.Run(t, new(TransactionDBTest))
}

func (t *TransactionDBTest) TestTransactionDB_Create() {
	clientFrom, _ := entity.NewClient("Client from", "client@from.com")
	clientTo, _ := entity.NewClient("Client to", "client@to.com")
	accountFrom, _ := entity.NewAccount(clientFrom, 100)
	accountTo, _ := entity.NewAccount(clientTo, 50)
	transaction, _ := entity.NewTransaction(accountFrom, accountTo, 50)
	err := t.transactionDB.Create(transaction)
	t.Nil(err)
	transactionDTO := &entity.Transaction{
		AccountFrom: &entity.Account{},
		AccountTo:   &entity.Account{},
	}
	t.db.QueryRow("select id, account_from_id, account_to_id, amount, created_at from transactions where id = ?", transaction.ID).Scan(&transactionDTO.ID, &transactionDTO.AccountFrom.ID, &transactionDTO.AccountTo.ID, &transactionDTO.Amount, &transactionDTO.CreatedAt)
	t.Equal(transaction.ID, transactionDTO.ID)
	t.Equal(transaction.AccountFrom.ID, transactionDTO.AccountFrom.ID)
	t.Equal(transaction.AccountTo.ID, transactionDTO.AccountTo.ID)
	t.Equal(transaction.Amount, transactionDTO.Amount)
	t.Equal(transaction.CreatedAt.UTC(), transactionDTO.CreatedAt.UTC())
}
