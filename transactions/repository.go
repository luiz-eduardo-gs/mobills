package transactions

import (
	"database/sql"
	"fmt"
	"log"
	"time"
)

type TransactionRepository interface {
	Add(t Transaction) error
	AddFixed(t Transaction) error
	List(from time.Time, to time.Time) ([]Transaction, error)
}

type SqliteTransactionRepository struct {
	DB *sql.DB
}

func (r *SqliteTransactionRepository) Add(t Transaction) error {
	tx, err := r.DB.Begin()
	if err != nil {
		log.Fatal(err)
		return err
	}
	stmt, err := tx.Prepare("insert into transactions(value,paid,date,description,category_id,account_id,tag_id,type,fixed) values(?,?,?,?,?,?,?,?,?)")
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec()
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

func (r *SqliteTransactionRepository) AddFixed(t Transaction) error {
	tx, err := r.DB.Begin()
	if err != nil {
		log.Fatal(err)
		return err
	}
	stmt, err := tx.Prepare("insert into fixed_transactions(value,paid,day,description,category_id,account_id,tag_id,type) values(?,?,?,?,?,?,?,?)")
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer stmt.Close()

	day := t.Date.Day()

	_, err = stmt.Exec(t.Value, t.Paid, day, t.Description, t.CategoryID, t.AccountID, t.TagID, t.Type)
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

func (r *SqliteTransactionRepository) List(from time.Time, to time.Time) ([]Transaction, error) {
	// dayFrom := from.Day()
	// dayTo := to.Day()

	rows, err := r.DB.Query(fmt.Sprintf(`
	select 
		id,
		value,
		paid,
		day,
		description,
		category_id,
		account_id,
		tag_id,
		type
	from fixed_transactions
	`))
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var transactions []Transaction

	for rows.Next() {
		t := Transaction{}

		var day int

		err = rows.Scan(&t.ID, &t.Value, &t.Paid, &day, &t.Description, &t.CategoryID, &t.AccountID, &t.TagID, &t.Type)
		if err != nil {
			log.Fatal(err)
		}

		t.Fixed = true

		for i := int(from.Month()); i <= int(to.Month()); i++ {
			log.Println("i: ", i)
			if i == int(from.Month()) && from.Day() > day {
				continue
			}

			t.Date = CustomDate{
				time.Date(from.Year(), time.Month(i), day, 0, 0, 0, 0, time.UTC),
			}
			transactions = append(transactions, t)
		}
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	return transactions, nil
}

type InMemoryTransactionRepository struct {
	transactions []Transaction
}

func (r *InMemoryTransactionRepository) Add(t Transaction) error {
	r.transactions = append(r.transactions, t)
	return nil
}
