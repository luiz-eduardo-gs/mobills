package accounts

import (
	"database/sql"
	"log"
)

type AccountRepository interface {
	Add(a Account) error
	GetAll() ([]Account, error)
}

type SqliteAccountRepository struct {
	DB *sql.DB
}

func (r *SqliteAccountRepository) Add(a Account) error {
	tx, err := r.DB.Begin()
	if err != nil {
		log.Fatal(err)
		return err
	}
	stmt, err := tx.Prepare("insert into accounts(balance, description, type) values(?,?,?)")
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(a.Balance, a.Description, a.Type)
	if err != nil {
		log.Fatal(err)
		return err
	}
	err = tx.Commit()
	if err != nil {
		log.Fatal(err)
		return err
	}

	return nil
}

func (r *SqliteAccountRepository) GetAll() ([]Account, error) {
	rows, err := r.DB.Query("select id, balance, description, type from accounts")
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer rows.Close()

	var accounts []Account

	for rows.Next() {
		a := Account{}

		err = rows.Scan(&a.ID, &a.Balance, &a.Description, &a.Type)
		if err != nil {
			log.Fatal(err)
		}

		accounts = append(accounts, a)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	return accounts, nil
}

type InMemoryAccountRepository struct {
	accounts []Account
}

func (r *InMemoryAccountRepository) Add(a Account) error {
	r.accounts = append(r.accounts, a)
	return nil
}

func (r *InMemoryAccountRepository) GetAll() ([]Account, error) {
	return r.accounts, nil
}
