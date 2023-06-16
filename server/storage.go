package main

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type Storage interface {
	CreateAccount(*Account) error
	DeleteAccount(uuid.UUID) error
	UpdateAccount(*Account) error
	GetAccountByID(uuid.UUID) (*Account, error)
	GetAccounts() ([]*Account, error)
}

type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStore() (*PostgresStore, error) {
	// TO DO remove password from here
	connStr := "user=postgres dbname=trackitdb sslmode=disable"
	db, err := sql.Open("postgres", connStr)

	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &PostgresStore{
		db: db,
	}, nil
}

// TO DO drop tables, migrations
func (s *PostgresStore) Init() error {
	return s.createAccountTable()
}

func (s *PostgresStore) createAccountTable() error {
	statement := `CREATE TABLE IF NOT EXISTS accounts (
		ID uuid PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
		Username varchar(255),
		Email varchar(255),
		CreatedAt timestamp
	);`

	_, err := s.db.Exec(statement)
	return err
}
func (s *PostgresStore) CreateAccount(a *Account) error {
	statement := `INSERT INTO accounts (Username, Email, CreatedAt) 
			VALUES ($1, $2, $3)`

	resp, err := s.db.Exec(statement, a.Username, a.Email, a.CreatedAt)

	if err != nil {
		return err
	}

	fmt.Printf("%+v\n", resp)

	return nil
}

func (s *PostgresStore) UpdateAccount(a *Account) error {
	return nil
}

func (s *PostgresStore) DeleteAccount(id uuid.UUID) error {
	return nil
}

func (s *PostgresStore) GetAccountByID(id uuid.UUID) (*Account, error) {
	return &Account{}, nil
}

func (s *PostgresStore) GetAccounts() ([]*Account, error) {
	query := `SELECT * FROM accounts`

	rows, err := s.db.Query(query)

	if err != nil {
		return nil, err
	}

	accounts := []*Account{}

	for rows.Next() {
		account := new(Account)

		err := rows.Scan(
			&account.ID,
			&account.Username,
			&account.Email,
			&account.CreatedAt)

		if err != nil {
			return nil, err
		}

		accounts = append(accounts, account)
	}

	return accounts, nil
}
