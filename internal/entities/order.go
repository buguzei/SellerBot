package entities

import "time"

type CurrentOrder struct {
	UserID      int64
	ID          int64
	Start       time.Time
	Composition []Product
}

type DoneOrder struct {
	UserID      int64
	ID          int64
	Start       time.Time
	Done        time.Time
	Composition []Product
}
