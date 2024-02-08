package postgres

import (
	"bot/internal/entities"
	"bot/internal/log"
	"fmt"
)

// user methods

func (p Postgres) InsertUser(user entities.User) error {
	_, err := p.DB.Exec("INSERT INTO users VALUES ($1) ON CONFLICT DO NOTHING;", user.ID)

	if err != nil {
		p.logger.Error("InsertUser: insert error", log.Fields{
			"error": err,
		})

		return err
	}

	return nil
}

func (p Postgres) GetUser(userID int64) *entities.User {
	var user entities.User

	row := p.DB.QueryRow("SELECT * FROM users WHERE id=($1)", userID)

	err := row.Scan(&user.ID, &user.Name, &user.Address)
	if err != nil {
		fmt.Println(err)
	}

	return &user
}

func (p Postgres) UpdateUser(user entities.User) {
	_, err := p.DB.Exec("UPDATE users SET name = ($1), address = ($2) WHERE id=($3);", user.Name, user.Address, user.ID)

	if err != nil {
		p.logger.Error("UpdateUser: update error", log.Fields{
			"error": err,
		})
	}
}

// order methods

func (p Postgres) InsertOrder(order entities.Order) (*int64, error) {
	var orderID int64

	err := p.DB.QueryRow("INSERT INTO current_orders(user_id, start) VALUES (($1), ($2)) RETURNING id;", order.UserID, order.Date).Scan(&orderID)
	if err != nil {
		p.logger.Error("InsertOrder: insert error", log.Fields{
			"error": err,
		})

		return nil, err
	}

	return &orderID, nil
}

func (p Postgres) NewCurrentProducts(order entities.Order) error {
	for _, prod := range order.Composition {
		_, err := p.DB.Exec("INSERT INTO products(current_order_id, name, size, color, text, img, amount) VALUES (($1), ($2), ($3), ($4), ($5), ($6), ($7));",
			order.ID,
			prod.Name,
			prod.Size,
			prod.Color,
			prod.Text,
			prod.Img,
			prod.Amount,
		)
		if err != nil {
			p.logger.Error("NewCurrentProducts: insert error", log.Fields{
				"error": err,
			})

			return err
		}
	}

	return nil
}

func (p Postgres) GetAllCurrentOrders() []entities.Order {
	return nil
}
