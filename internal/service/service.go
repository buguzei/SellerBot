package service

import (
	"bot/internal/entities"
	log2 "bot/internal/log"
	"bot/internal/repo"
)

type Service interface {
	UserSvc
	OrderSvc
	CartSvc
}

type Svc struct {
	db     repo.Repo
	logger log2.Logger
}

func NewService(db repo.Repo) Svc {
	logger := log2.NewLogrus("debug")

	return Svc{db: db, logger: logger}
}

type UserSvc interface {
	NewUser(entities.User) error
	GetUser(int64) *entities.User
	UpdateUser(entities.User)
}

type OrderSvc interface {
	NewCurrentOrder(entities.CurrentOrder) error
	GetAllCurrentOrders() []entities.CurrentOrder
	GetAllDoneOrders() []entities.DoneOrder
	NewDoneOrder(int64)
}

type CartSvc interface {
	NewCartProduct(int64, int, entities.Product)
	CartLen(int64) int64
	GetCartProduct(int64, int) entities.Product
	GetCart(int64) map[int]entities.Product
	DeleteProductFromCart(int64, int)
}
