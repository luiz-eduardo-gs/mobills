package tags

import (
	"database/sql"
	"log"
)

type TagRepository interface {
	Add(t Tag) error
	GetAll() ([]Tag, error)
}

type SqliteTagRepository struct {
	DB *sql.DB
}

func (r *SqliteTagRepository) Add(t Tag) error {
	tx, err := r.DB.Begin()
	if err != nil {
		log.Fatal(err)
		return err
	}
	stmt, err := tx.Prepare("insert into tags(name) values(?)")
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(t.Name)
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

func (r *SqliteTagRepository) GetAll() ([]Tag, error) {
	rows, err := r.DB.Query("select id, name from tags")
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer rows.Close()

	var tags []Tag

	for rows.Next() {
		var id uint64
		var name string

		err = rows.Scan(&id, &name)
		if err != nil {
			log.Fatal(err)
		}

		tags = append(tags, Tag{ID: id, Name: name})
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	return tags, nil
}

type InMemoryTagRepository struct {
	tags []Tag
}

func (r *InMemoryTagRepository) Add(t Tag) error {
	r.tags = append(r.tags, t)
	return nil
}

func (r *InMemoryTagRepository) GetAll() ([]Tag, error) {
	return r.tags, nil
}
