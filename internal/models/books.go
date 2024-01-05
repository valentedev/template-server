package models

import (
	"database/sql"
	"errors"
	"time"
)

type Book struct {
	ID      int
	Created time.Time
	Title   string
	Author  string
}

type BookModel struct {
	DB *sql.DB
}

func (m *BookModel) Insert(title, author string) (int, error) {

	var id int

	query := `
		INSERT INTO books(title, author) VALUES ($1, $2) RETURNING id;
	`

	err := m.DB.QueryRow(query, title, author).Scan(&id)

	return int(id), err
}

func (m *BookModel) Get(id int) (Book, error) {

	query := `
		SELECT id, created, title, author FROM books WHERE id=$1;
	`

	row := m.DB.QueryRow(query, id)
	var b Book

	err := row.Scan(&b.ID, &b.Created, &b.Title, &b.Author)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Book{}, ErrNoRecord
		} else {
			return Book{}, err
		}
	}

	return b, err
}

func (m *BookModel) Latest() ([]Book, error) {
	query := `
		SELECT title, author, id FROM books ORDER BY id DESC LIMIT 10;
	`

	rows, err := m.DB.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var books []Book

	for rows.Next() {
		var b Book
		err = rows.Scan(&b.Title, &b.Author, &b.ID)
		if err != nil {
			return nil, err
		}
		books = append(books, b)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return books, err
}
