package models

import (
	"errors"
	"time"
)

var ErrNoRecord = errors.New("models: подходящей записи не найдено")

// Структура модели данных
type Snippet struct {
	Id      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}
