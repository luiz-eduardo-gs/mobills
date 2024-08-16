package categories

import (
	"database/sql"
	"log"
)

type CategoryRepository interface {
	Add(c Category) error
	GetAll() ([]Category, error)
}

type SqliteCategoryRepository struct {
	DB *sql.DB
}

func (r *SqliteCategoryRepository) Add(c Category) error {
	tx, err := r.DB.Begin()
	if err != nil {
		log.Fatal(err)
		return err
	}
	stmt, err := tx.Prepare("insert into categories(name, type) values(?,?)")
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(c.Name, c.Type)
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

func (r *SqliteCategoryRepository) GetAll() ([]Category, error) {
	rows, err := r.DB.Query("select id, name, type from categories")
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer rows.Close()

	var categories []Category

	for rows.Next() {
		c := Category{}

		err = rows.Scan(&c.ID, &c.Name, &c.Type)
		if err != nil {
			log.Fatal(err)
		}

		categories = append(categories, c)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	return categories, nil
}

type InMemoryCategoryRepository struct {
	categories []Category
}

func (r *InMemoryCategoryRepository) Add(c Category) error {
	r.categories = append(r.categories, c)
	return nil
}

func (r *InMemoryCategoryRepository) GetAll() ([]Category, error) {
	return r.categories, nil
}
