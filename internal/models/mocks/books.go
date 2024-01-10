package mocks

import (
	"time"

	"github.com/valentedev/template-server/internal/models"
)

var mockBooks = models.Book{
	ID:      1,
	Title:   "The Gobblet of Fire",
	Author:  "J.K.Rowling",
	Created: time.Now(),
}

type BookModel struct{}

func (m *BookModel) Insert(title, author string) (int, error) {
	return 2, nil
}

func (m *BookModel) Get(id int) (models.Book, error) {
	switch id {
	case 1:
		return mockBooks, nil
	default:
		return models.Book{}, models.ErrNoRecord
	}
}

func (m *BookModel) Latest() ([]models.Book, error) {
	return []models.Book{mockBooks}, nil
}
