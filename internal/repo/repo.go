package repo

import (
	"bot/internal/entities"
)

type Repo interface {
	UserRepo
	CartRepo
	OrderRepo
}

type UserRepo interface {
	InsertUser(entities.User) error
	GetUser(int64) *entities.User
	UpdateUser(entities.User)
}

type OrderRepo interface {
	InsertOrder(entities.CurrentOrder) (*int64, error)
	NewCurrentProducts(order entities.CurrentOrder) error
	GetAllCurrentOrders() []entities.CurrentOrder
	GetAllDoneOrders() []entities.DoneOrder
	FromCurrentToDone(int64)
}

type CartRepo interface {
	NewCartProduct(int64, int, entities.Product)
	GetCartProduct(int64, int) entities.Product
	CartLen(int64) int64
	GetCart(int64) []entities.Product
	ClearCart(int64)
}

type CombineRepos struct {
	UserRepo
	OrderRepo
	CartRepo
}
