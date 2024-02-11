package service

import (
	"bot/internal/entities"
	"bot/internal/log"
)

//user methods

func (s Svc) NewUser(user entities.User) error {
	err := s.db.InsertUser(user)
	if err != nil {
		s.logger.Error("adding new user to db failed", log.Fields{
			"error": err,
		})

		return err
	}

	return nil
}

func (s Svc) GetUser(userID int64) *entities.User {
	user := s.db.GetUser(userID)

	return user
}

func (s Svc) UpdateUser(user entities.User) {
	s.db.UpdateUser(user)
}

// cart methods

func (s Svc) NewCartProduct(userID int64, idx int, product entities.Product) {
	s.db.NewCartProduct(userID, idx, product)
}

func (s Svc) CartLen(userID int64) int64 {
	length := s.db.CartLen(userID)

	return length
}

func (s Svc) GetCartProduct(userID int64, idx int) entities.Product {
	prod := s.db.GetCartProduct(userID, idx)

	return prod
}

func (s Svc) GetCart(userID int64) map[int]entities.Product {
	cart := s.db.GetCart(userID)
	return cart
}

func (s Svc) DeleteProductFromCart(userID int64, idx int) {
	s.db.DeleteProductFromCart(userID, idx)
}

// order methods

func (s Svc) NewOrder(order entities.CurrentOrder) error {
	orderID, err := s.db.InsertOrder(order)
	if err != nil {
		s.logger.Error("NewOrder: creating order failed", log.Fields{
			"error": err,
		})

		return err
	}

	order.ID = *orderID

	err = s.db.NewCurrentProducts(order)
	if err != nil {
		s.logger.Error("NewOrder: creating product failed", log.Fields{
			"error": err,
		})

		return err
	}

	s.db.ClearCart(order.UserID)

	return nil
}

func (s Svc) GetAllCurrentOrders() []entities.CurrentOrder {
	orders := s.db.GetAllCurrentOrders()

	return orders
}

func (s Svc) FromCurrentToDone(orderID int64) {
	s.db.FromCurrentToDone(orderID)
}

func (s Svc) GetAllDoneOrders() []entities.DoneOrder {
	orders := s.db.GetAllDoneOrders()

	return orders
}
