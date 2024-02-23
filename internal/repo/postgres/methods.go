package postgres

import (
	"bot/internal/entities"
	"bot/internal/log"
	"database/sql"
	"time"
)

// user methods

func (p Postgres) InsertUser(user entities.User) error {
	_, err := p.DB.Exec("INSERT INTO users VALUES (($1), ($2), ($3), ($4)) ON CONFLICT DO NOTHING;", user.ID, user.Name, user.Phone, user.Address)

	if err != nil {
		p.logger.Error("InsertUser: insert error", log.Fields{
			"error": err,
		})

		return err
	}

	return nil
}

type result struct {
	name    sql.NullString
	phone   sql.NullString
	address sql.NullString
}

func (p Postgres) GetUser(userID int64) (*entities.User, error) {
	var user entities.User
	var res result

	row := p.DB.QueryRow("SELECT * FROM users WHERE id=($1)", userID)

	err := row.Scan(&user.ID, &res.name, &res.phone, &res.address)
	if err != nil {
		p.logger.Error("GetUser: scan error", log.Fields{
			"error": err,
		})

		return nil, err
	}

	if !res.name.Valid {
		user.Name = ""
	} else {
		user.Name = res.name.String
	}

	if !res.phone.Valid {
		user.Phone = ""
	} else {
		user.Phone = res.phone.String
	}

	if !res.address.Valid {
		user.Address = ""
	} else {
		user.Address = res.address.String
	}

	return &user, nil
}

func (p Postgres) UpdateUser(user entities.User) error {
	_, err := p.DB.Exec("UPDATE users SET name = ($1), address = ($2), phone=($3) WHERE id=($4);", user.Name, user.Address, user.Phone, user.ID)

	if err != nil {
		p.logger.Error("UpdateUser: update error", log.Fields{
			"error": err,
		})

		return err
	}

	return nil
}

// order methods

func (p Postgres) NewCurrentOrder(order entities.CurrentOrder) (*int64, error) {
	var orderID int64

	err := p.DB.QueryRow("INSERT INTO current_orders(user_id, start) VALUES (($1), ($2)) RETURNING id;", order.UserID, order.Start).Scan(&orderID)
	if err != nil {
		p.logger.Error("NewCurrentOrder: insert error", log.Fields{
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

func (p Postgres) GetAllCurrentOrders() ([]entities.CurrentOrder, error) {
	var resOrders []entities.CurrentOrder
	orderMap := make(map[int64]entities.CurrentOrder)

	rows, err := p.DB.Query("SELECT products.id, name, size, color, text, img, amount, current_orders.id, user_id, start FROM products JOIN current_orders ON products.current_order_id = current_orders.id;")
	if err != nil {
		p.logger.Error("GetAllCurrentOrders: select join error", log.Fields{
			"error": err,
		})

		return nil, err
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

			return nil, err
		}

		if value, ok := orderMap[order.ID]; !ok {
			order.Composition = append(order.Composition, product)

			orderMap[order.ID] = order
		} else {
			value.Composition = append(value.Composition, product)

			orderMap[order.ID] = value
		}
	}

	for _, value := range orderMap {
		resOrders = append(resOrders, value)
	}

	return resOrders, nil
}

func (p Postgres) NewDoneOrder(orderID int64) error {
	tx, err := p.DB.Begin()
	if err != nil {
		p.logger.Error("NewDoneOrder: begin transaction error", log.Fields{
			"error": err,
		})

		return err
	}

	defer func() {
		err = tx.Rollback()
		if err != nil {
			p.logger.Error("NewDoneOrder: rollback error", log.Fields{
				"error": err,
			})
		}
	}()

	// start transaction

	var order entities.CurrentOrder

	row := p.DB.QueryRow("SELECT id, user_id, start FROM current_orders WHERE id=($1);", orderID)

	err = row.Scan(&order.ID, &order.UserID, &order.Start)
	if err != nil {
		p.logger.Error("NewDoneOrder: order scan error", log.Fields{
			"error": err,
		})

		return err
	}

	row = p.DB.QueryRow("INSERT INTO done_orders(user_id, start, done) VALUES (($1), ($2), ($3)) RETURNING done_orders.id;", order.UserID, order.Start, time.Now())

	var doneID int64

	err = row.Scan(&doneID)
	if err != nil {
		p.logger.Error("NewDoneOrder: done id scan error", log.Fields{
			"error": err,
		})

		return err
	}

	_, err = p.DB.Exec("UPDATE products SET current_order_id=null, done_order_id=($1) WHERE current_order_id=($2);", doneID, order.ID)
	if err != nil {
		p.logger.Error("NewDoneOrder: updating products error", log.Fields{
			"error": err,
		})

		return err
	}

	_, err = p.DB.Exec("DELETE FROM current_orders WHERE id=($1);", orderID)
	if err != nil {
		p.logger.Error("NewDoneOrder: deleting current_order error", log.Fields{
			"error": err,
		})

		return err
	}
	// end transaction
	err = tx.Commit()
	if err != nil {
		p.logger.Error("NewDoneOrder: commit transaction error", log.Fields{
			"error": err,
		})

		return err
	}

	return nil
}

func (p Postgres) GetAllDoneOrders() ([]entities.DoneOrder, error) {
	var resOrders []entities.DoneOrder
	orderMap := make(map[int64]entities.DoneOrder)

	rows, err := p.DB.Query("SELECT products.id, name, size, color, text, img, amount, done_orders.id, user_id, start, done FROM products JOIN done_orders ON products.done_order_id = done_orders.id;")
	if err != nil {
		p.logger.Error("GetAllDoneOrders: select join error", log.Fields{
			"error": err,
		})

		return nil, err
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

			return nil, err
		}

		if value, ok := orderMap[order.ID]; !ok {
			order.Composition = append(order.Composition, product)

			orderMap[order.ID] = order
		} else {
			value.Composition = append(value.Composition, product)

			orderMap[order.ID] = value
		}
	}

	for _, value := range orderMap {
		resOrders = append(resOrders, value)
	}

	return resOrders, nil
}
