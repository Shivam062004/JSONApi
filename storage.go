package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

type Storage interface {
	storeAccount(*Account) (int64, error)
	getAccountByID(int) (*Account, error)
	getAccounts() ([]*Account, error)
	deleteAccountByID(int) error
	setUp() error
}

type Postgres struct {
	db *sql.DB
}

func createPostGresStore() (*Postgres, error) {
	pass := os.Getenv("DB_PASS")
	connStr := fmt.Sprintf("user=postgres password=%s dbname=postgres host=localhost port=5432 sslmode=disable", pass)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return &Postgres{
		db: db,
	}, nil
}

func (p *Postgres) setUp() error {
	query := `CREATE TABLE IF NOT EXISTS account (
    id BIGSERIAL PRIMARY KEY,
    firstname VARCHAR(50) NOT NULL,
    lastname VARCHAR(50) NOT NULL,
    balance BIGINT NOT NULL,
    number BIGINT NOT NULL,
	password VARCHAR(100) NOT NULL
	)`
	_, err := p.db.Exec(query)
	return err
}

func (p *Postgres) storeAccount(account *Account) (int64, error) {
	query := `INSERT INTO account (firstname, lastname, balance, number, password) VALUES ($1, $2, $3, $4, $5) RETURNING id`
	var id int64
	err := p.db.QueryRow(query, account.FirstName, account.LastName, account.Balance, account.Number, account.Password).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (p *Postgres) getAccounts() ([]*Account, error) {
	query := `SELECT * FROM account`
	rows, err := p.db.Query(query)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	var accounts []*Account = []*Account{}

	for rows.Next() {
		newAcc := &Account{}
		err := rows.Scan(&newAcc.ID, &newAcc.FirstName, &newAcc.LastName, &newAcc.Balance, &newAcc.Number, &newAcc.Password)
		if err != nil {
			log.Fatal(err)
		}
		accounts = append(accounts, newAcc)
	}

	return accounts, nil
}

func (p *Postgres) getAccountByID(id int) (*Account, error) {
	query := `SELECT * FROM account WHERE ID = ($1)`
	acc := &Account{}
	err := p.db.QueryRow(query, id).Scan(&acc.ID, &acc.FirstName, &acc.LastName, &acc.Balance, &acc.Number, &acc.Password)
	if err != nil {
		return nil, err
	}
	return acc, nil
}

func (p *Postgres) deleteAccountByID(id int) error {
	query := `DELETE FROM account WHERE ID = $1`
	_, err := p.db.Exec(query, id)
	return err
}
