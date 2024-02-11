package postgres

import (
	"bot/internal/entities"
	"bot/internal/log"
	"time"
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
		p.logger.Error("GetUser: scan error", log.Fields{
			"error": err,
		})
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

func (p Postgres) InsertOrder(order entities.CurrentOrder) (*int64, error) {
	var orderID int64

	err := p.DB.QueryRow("INSERT INTO current_orders(user_id, start) VALUES (($1), ($2)) RETURNING id;", order.UserID, order.Start).Scan(&orderID)
	if err != nil {
		p.logger.Error("InsertOrder: insert error", log.Fields{
			"error": err,
		})

		return nil, err
	}

	return &orderID, nil
}

func (p Postgres) NewCurrentProducts(order entities.CurrentOrder) error {
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

func (p Postgres) GetAllCurrentOrders() []entities.CurrentOrder {
	var orders []entities.CurrentOrder

	rows, err := p.DB.Query("SELECT products.id, name, size, color, text, img, amount, current_orders.id, user_id, start FROM products JOIN current_orders ON products.current_order_id = current_orders.id;")
	if err != nil {
		p.logger.Error("GetAllCurrentOrders: select join error", log.Fields{
			"error": err,
		})

		return nil
	}

	for rows.Next() {
		var order entities.CurrentOrder
		var product entities.Product

		err = rows.Scan(
			&product.ID,
			&product.Name,
			&product.Size,
			&product.Color,
			&product.Text,
			&product.Img,
			&product.Amount,
			&order.ID,
			&order.UserID,
			&order.Start,
		)
		if err != nil {
			p.logger.Error("GetAllCurrentOrders: scan error", log.Fields{
				"error": err,
			})

			return nil
		}

		order.Composition = append(order.Composition, product)

		orders = append(orders, order)
	}

	return orders
}

func (p Postgres) FromCurrentToDone(orderID int64) {
	tx, err := p.DB.Begin()
	if err != nil {
		p.logger.Error("FromCurrentToDone: begin transaction error", log.Fields{
			"error": err,
		})
	}
	// start transaction

	var order entities.CurrentOrder

	row := p.DB.QueryRow("SELECT id, user_id, start FROM current_orders WHERE id=($1);", orderID)

	err = row.Scan(&order.ID, &order.UserID, &order.Start)
	if err != nil {
		p.logger.Error("FromCurrentToDone: order scan error", log.Fields{
			"error": err,
		})
	}

	row = p.DB.QueryRow("INSERT INTO done_orders(user_id, start, done) VALUES (($1), ($2), ($3)) RETURNING done_orders.id;", order.UserID, order.Start, time.Now())

	var doneID int64

	err = row.Scan(&doneID)
	if err != nil {
		p.logger.Error("FromCurrentToDone: done id scan error", log.Fields{
			"error": err,
		})
	}

	_, err = p.DB.Exec("UPDATE products SET current_order_id=null, done_order_id=($1) WHERE current_order_id=($2);", doneID, order.ID)
	if err != nil {
		p.logger.Error("FromCurrentToDone: updating products error", log.Fields{
			"error": err,
		})
	}

	_, err = p.DB.Exec("DELETE FROM current_orders WHERE id=($1);", orderID)
	if err != nil {
		p.logger.Error("FromCurrentToDone: deleting current_order error", log.Fields{
			"error": err,
		})
	}
	// end transaction
	err = tx.Commit()
	if err != nil {
		p.logger.Error("FromCurrentToDone: commit transaction error", log.Fields{
			"error": err,
		})
	}
}

func (p Postgres) GetAllDoneOrders() []entities.DoneOrder {
	var orders []entities.DoneOrder

	rows, err := p.DB.Query("SELECT products.id, name, size, color, text, img, amount, done_orders.id, user_id, start, done FROM products JOIN done_orders ON products.done_order_id = done_orders.id;")
	if err != nil {
		p.logger.Error("GetAllDoneOrders: select join error", log.Fields{
			"error": err,
		})

		return nil
	}

	for rows.Next() {
		var order entities.DoneOrder
		var product entities.Product

		err = rows.Scan(
			&product.ID,
			&product.Name,
			&product.Size,
			&product.Color,
			&product.Text,
			&product.Img,
			&product.Amount,
			&order.ID,
			&order.UserID,
			&order.Start,
			&order.Done,
		)
		if err != nil {
			p.logger.Error("GetAllDoneOrders: scan error", log.Fields{
				"error": err,
			})

			return nil
		}

		order.Composition = append(order.Composition, product)

		orders = append(orders, order)
	}

	return orders
}
