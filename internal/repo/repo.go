package repo

import (
	"bot/internal/entities"
)

//go:generate mockgen -source=./repo.go -destination=./mocks/repository-mock.go

type UserRepo interface {
	InsertUser(entities.User) error
	GetUser(int64) (*entities.User, error)
	UpdateUser(entities.User) error
}

type OrderRepo interface {
	NewCurrentOrder(entities.CurrentOrder) (*int64, error)
	NewCurrentProducts(entities.CurrentOrder) error
	GetAllCurrentOrders() ([]entities.CurrentOrder, error)
	GetAllDoneOrders() ([]entities.DoneOrder, error)
	NewDoneOrder(int64) error
}

type CartRepo interface {
	NewCartProduct(int64, int, entities.Product) error
	GetCartProduct(int64, int) (*entities.Product, error)
	CartLen(int64) (*int64, error)
	GetCart(int64) (map[int]entities.Product, error)
	ClearCart(int64)
	DeleteProductFromCart(int64, int)
}

type Repo struct {
	UserRepo
	OrderRepo
	CartRepo
}
