package database

import (
	"database/sql"

	"github.com/NatanSiilva/ms-wallet/internal/entity"
)

type ClientDB struct {
	DB *sql.DB
}

func NewClientDB(db *sql.DB) *ClientDB {
	return &ClientDB{DB: db}
}

func (db *ClientDB) Get(id string) (*entity.Client, error) {
	client := &entity.Client{}
	stmt, err := db.DB.Prepare("SELECT id, name, email, created_at FROM clients WHERE id = $1")

	if err != nil {
		return nil, err
	}

	defer stmt.Close()

	row := stmt.QueryRow(id)

	if err := row.Scan(&client.ID, &client.Name, &client.Email, &client.CreatedAt); err != nil {
		return nil, err
	}

	return client, nil
}

func (db *ClientDB) Save(client *entity.Client) error {
	stmt, err := db.DB.Prepare("INSERT INTO clients (id, name, email, created_at) VALUES ($1, $2, $3, $4)")

	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(client.ID, client.Name, client.Email, client.CreatedAt)

	if err != nil {
		return err
	}

	return nil
}
