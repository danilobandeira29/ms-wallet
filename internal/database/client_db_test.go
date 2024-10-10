package database

import (
	"database/sql"
	"github.com/danilobandeira29/ms-wallet/internal/entity"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/suite"
	"testing"
)

type ClientDBTest struct {
	suite.Suite
	db       *sql.DB
	clientDB *ClientDB
}

func (s *ClientDBTest) SetupSuite() {
	db, err := sql.Open("sqlite3", ":memory:")
	s.Nil(err)
	s.db = db
	_, err = s.db.Exec("create table clients(id varchar(255) primary key, name varchar(255) not null, email varchar(255) not null, created_at date, updated_at date)")
	s.clientDB = NewClientDB(db)
}

func (s *ClientDBTest) TearDownSuite() {
	defer s.db.Close()
	s.db.Exec("drop table clients")
}

func TestClientDBSuite(t *testing.T) {
	suite.Run(t, new(ClientDBTest))
}

func (s *ClientDBTest) TestClientDB_Get() {
	client, _ := entity.NewClient("Danilo Bandeira", "danilo@email.com")
	s.clientDB.Save(client)
	clientDTO, err := s.clientDB.Get(client.ID)
	s.Nil(err)
	s.Equal(client.ID, clientDTO.ID)
}

func (s *ClientDBTest) TestClientDB_Save() {
	client, _ := entity.NewClient("Danilo Bandeira", "danilo@email.com")
	err := s.clientDB.Save(client)
	s.Nil(err)
	var clientDTO struct {
		ID string
	}
	err = s.db.QueryRow("select id from clients where id = ?", client.ID).Scan(&clientDTO.ID)
	s.Nil(err)
	s.Equal(client.ID, clientDTO.ID)
}
