package database

import (
	"database/sql"
	"github.com/danilobandeira29/ms-wallet/internal/entity"
	"github.com/stretchr/testify/suite"
	"testing"
)

type AccountDBTest struct {
	suite.Suite
	db        *sql.DB
	accountDB *AccountDB
}

func (a *AccountDBTest) SetupSuite() {
	db, err := sql.Open("sqlite3", ":memory:")
	a.Nil(err)
	a.db = db
	a.db.Exec("create table clients(id varchar(255) primary key, name varchar(255) not null, email varchar(255) not null, created_at date, updated_at date)")
	a.db.Exec(`
	create table accounts(
	    id varchar(255) primary key, 
	    balance int not null, 
	    client_id varchar(255),
	    created_at date not null,
	    updated_at date,
	    foreign key (client_id) references clients(id) 
	)
`)
	a.accountDB = NewAccountDB(a.db)
}

func (a *AccountDBTest) TearDownSuite() {
	defer a.db.Close()
	a.db.Exec("drop table clients")
	a.db.Exec("drop table accounts")
}

func TestAccountDBSuite(t *testing.T) {
	suite.Run(t, new(AccountDBTest))
}

func (a *AccountDBTest) TestAccountDB_Save() {
	client, _ := entity.NewClient("Danilo Bandeira", "danilo@email.com")
	account, _ := entity.NewAccount(client, 100)
	err := a.accountDB.Save(account)
	a.Nil(err)
}

func (a *AccountDBTest) TestAccountDB_FindByID() {
	client, _ := entity.NewClient("Danilo Bandeira", "danilo@email.com")
	account, _ := entity.NewAccount(client, 100)
	a.accountDB.DB.Exec("insert into clients(id, name, email, created_at) values(?, ?, ?, ?)", client.ID, client.Name, client.Email, client.CreatedAt)
	a.accountDB.Save(account)
	acc, err := a.accountDB.FindBy(account.ID)
	a.Nil(err)
	a.Equal(account.ID, acc.ID)
	a.Equal(account.BalanceInCents, acc.BalanceInCents)
	a.Equal(account.CreatedAt.UTC(), acc.CreatedAt.UTC())
	a.Equal(account.Client.ID, acc.Client.ID)
	a.Equal(account.Client.Name, acc.Client.Name)
	a.Equal(account.Client.Email, acc.Client.Email)
}
