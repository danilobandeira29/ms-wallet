package database

import (
	"database/sql"
	"fmt"
	"github.com/danilobandeira29/ms-wallet/internal/entity"
)

type ClientDB struct {
	DB *sql.DB
}

func NewClientDB(db *sql.DB) *ClientDB {
	return &ClientDB{
		DB: db,
	}
}

func (c *ClientDB) Get(id string) (*entity.Client, error) {
	client := &entity.Client{}
	stmt, err := c.DB.Prepare("select id, name, email, created_at from clients where id = ?")
	if err != nil {
		return nil, err
	}
	defer func(stmt *sql.Stmt) {
		if errStmt := stmt.Close(); errStmt != nil {
			err = errStmt
		}
	}(stmt)
	row := stmt.QueryRow(id)
	if err = row.Scan(&client.ID, &client.Name, &client.Email, &client.CreatedAt); err != nil {
		return nil, err
	}
	return client, nil
}

func (c *ClientDB) Save(client *entity.Client) error {
	stmt, err := c.DB.Prepare("insert into clients(id, name, email, created_at) values (?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer func(stmt *sql.Stmt) {
		errStmt := stmt.Close()
		if errStmt != nil {
			err = errStmt
		}
	}(stmt)
	_, err = stmt.Exec(client.ID, client.Name, client.Email, client.CreatedAt)
	fmt.Println("exec", err)
	if err != nil {
		return err
	}
	return nil
}
