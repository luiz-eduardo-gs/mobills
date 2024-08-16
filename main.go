package main

import (
	"database/sql"
	"log"
	"mobills/accounts"
	"mobills/categories"
	"mobills/tags"
	"mobills/transactions"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// os.Remove("./mobills.db")

	db, err := sql.Open("sqlite3", "./mobills.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	sqlStmt := `
	create table if not exists tags (id integer not null primary key, name text);

	create table if not exists categories (id integer not null primary key, name text, type integer);
	
	create table if not exists accounts (id integer not null primary key, balance float, description text, type integer);

	create table if not exists transactions (
		id integer not null primary key,
		value float,
		paid integer,
		date timestamp,
		description text,
		category_id integer,
		account_id integer,
		tag_id integer,
		type integer,
		fixed integer
	);

	create table if not exists fixed_transactions (
		id integer not null primary key,
		value float,
		paid integer,
		day integer,
		description text,
		category_id integer,
		account_id integer,
		tag_id integer,
		type integer
	);
	`

	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	tagService := tags.TagService{
		Repository: &tags.SqliteTagRepository{
			DB: db,
		},
	}

	categoryService := categories.CategoryService{
		Repository: &categories.SqliteCategoryRepository{
			DB: db,
		},
	}

	accountRepo := accounts.AccountService{
		Repository: &accounts.SqliteAccountRepository{
			DB: db,
		},
	}

	transactionRepo := transactions.TransactionService{
		Repo: &transactions.SqliteTransactionRepository{
			DB: db,
		},
	}

	r.Group(func(r chi.Router) {
		r.Post("/tags", tagService.CreateTag)
		r.Get("/tags", tagService.ListTags)
	})

	r.Group(func(r chi.Router) {
		r.Post("/categories", categoryService.CreateCategory)
		r.Get("/categories", categoryService.ListCategories)
	})

	r.Group(func(r chi.Router) {
		r.Post("/accounts", accountRepo.CreateAccount)
		r.Get("/accounts", accountRepo.ListAccounts)
	})

	r.Group(func(r chi.Router) {
		r.Post("/transactions", transactionRepo.CreateTransaction)
		r.Get("/transactions", transactionRepo.ListTransactions)
	})

	log.Fatal(http.ListenAndServe(":8080", r))
}
