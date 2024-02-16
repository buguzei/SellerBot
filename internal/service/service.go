package service

import (
	"bot/internal/entities"
	"bot/internal/log"
	"bot/internal/repo"
)

type Service struct {
	db     repo.Repo
	logger log.Logger
}

// NewService is a constructor for Service
func NewService(db repo.Repo) Service {
	logger := log.NewLogrus("debug")

	return Service{db: db, logger: logger}
}

//user methods

func (s Service) NewUser(user entities.User) error {
	err := s.db.InsertUser(user)
	if err != nil {
		s.logger.Error("NewUser: adding new user to db failed", log.Fields{
			"error": err,
		})

		return err
	}

	return nil
}

func (s Service) GetUser(userID int64) (*entities.User, error) {
	user, err := s.db.GetUser(userID)
	if err != nil {
		s.logger.Error("GetUser: getting user from db failed", log.Fields{
			"error": err,
		})

		return nil, err
	}

	return user, nil
}

func (s Service) UpdateUser(user entities.User) error {
	err := s.db.UpdateUser(user)
	if err != nil {
		s.logger.Error("UpdateUser: updating user failed", log.Fields{
			"error": err,
		})

		return err
	}

	return nil
}

// cart methods

func (s Service) NewCartProduct(userID int64, idx int, product entities.Product) error {
	err := s.db.NewCartProduct(userID, idx, product)
	if err != nil {
		s.logger.Error("NewCartProduct: new product in cart failed", log.Fields{
			"error": err,
		})

		return err
	}

	return nil
}

func (s Service) CartLen(userID int64) (*int64, error) {
	length, err := s.db.CartLen(userID)
	if err != nil {
		s.logger.Error("CartLen: getting cart len failed", log.Fields{
			"error": err,
		})

		return nil, err
	}

	return length, nil
}

func (s Service) GetCartProduct(userID int64, idx int) (*entities.Product, error) {
	prod, err := s.db.GetCartProduct(userID, idx)
	if err != nil {
		s.logger.Error("GetCartProduct: getting product failed", log.Fields{
			"error": err,
		})

		return nil, err
	}

	return prod, nil
}

func (s Service) GetCart(userID int64) (map[int]entities.Product, error) {
	cart, err := s.db.GetCart(userID)
	if err != nil {
		s.logger.Error("GetCart: getting cart failed", log.Fields{
			"error": err,
		})

		return nil, err
	}

	return cart, nil
}

func (s Service) DeleteProductFromCart(userID int64, idx int) {
	s.db.DeleteProductFromCart(userID, idx)
}

// order methods

func (s Service) NewCurrentOrder(order entities.CurrentOrder) error {
	orderID, err := s.db.NewCurrentOrder(order)
	if err != nil {
		s.logger.Error("NewCurrentOrder: creating order failed", log.Fields{
			"error": err,
		})

		return err
	}

	order.ID = *orderID

	err = s.db.NewCurrentProducts(order)
	if err != nil {
		s.logger.Error("NewCurrentOrder: creating product failed", log.Fields{
			"error": err,
		})

		return err
	}

	s.db.ClearCart(order.UserID)

	return nil
}

func (s Service) GetAllCurrentOrders() ([]entities.CurrentOrder, error) {
	orders, err := s.db.GetAllCurrentOrders()
	if err != nil {
		s.logger.Error("GetAllCurrentOrders: getting orders from db failed", log.Fields{
			"error": err,
		})

		return nil, err
	}

	return orders, nil
}

func (s Service) NewDoneOrder(orderID int64) error {
	err := s.db.NewDoneOrder(orderID)
	if err != nil {
		s.logger.Error("NewDoneOrder: new done order error", log.Fields{
			"error": err,
		})

		return err
	}

	return nil
}

func (s Service) GetAllDoneOrders() ([]entities.DoneOrder, error) {
	orders, err := s.db.GetAllDoneOrders()
	if err != nil {
		s.logger.Error("GetAllDoneOrders: getting done orders failed", log.Fields{
			"error": err,
		})

		return nil, err
	}

	return orders, nil
}
