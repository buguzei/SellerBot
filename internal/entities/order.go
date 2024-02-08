package entities

import "time"

type Order struct {
	UserID      int64
	ID          int64
	Date        time.Time
	Composition []Product
}
