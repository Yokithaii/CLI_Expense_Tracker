package entities

import (
	"errors"
	"time"
)

type Expense struct {
	ID     int
	Date   time.Time
	Desc   string
	Amount int
}

func NewExpense(description string, amount int) (*Expense, error) {
	if amount < 0 {
		return nil, errors.New("bad amount promoted")
	}
	runes := []rune(description)
	if len(runes) > 30 {
		return nil, errors.New("your description is too long")
	}
	return &Expense{0, time.Now().UTC(), description, amount}, nil
}
