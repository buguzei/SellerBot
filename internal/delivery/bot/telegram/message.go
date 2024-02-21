package telegram

import (
	"bot/internal/entities"
	log2 "bot/internal/log"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strconv"
	"strings"
)

func (tg TGBot) HandleMessage(update tgbotapi.Update) {
	userID := update.Message.From.ID

	switch tg.cache[userID]["lvl"] {
	case "phone":
		tg.phoneLvlHandler(update.Message)
	case "name":
		tg.nameLvlHandler(update.Message)
	case "address":
		tg.addressLvlHandler(update.Message)
	case "print":
		tg.printLvlHandler(update)
	default:
		split := strings.Split(update.Message.Command(), "_")

		switch split[0] {
		case "start":
			tg.startCommandHandler(update.Message)
		case "current":
			var access bool

			for _, admin := range tg.conf.Admins {
				if userID == admin {
					access = true
				}
			}

			if access {
				tg.currentCommandHandler(update.Message)
			}
		case "done":
			var access bool

			for _, admin := range tg.conf.Admins {
				if userID == admin {
					access = true
				}
			}

			if access {
				if len(split) == 1 {
					tg.doneCommandHandler(update.Message)
				}

				if len(split) == 2 {
					tg.done_xCommandHandler(split)
				}
			}
		}
	}
}

func (tg TGBot) phoneLvlHandler(message *tgbotapi.Message) {
	userID := message.From.ID

	user, err := tg.svc.GetUser(userID)
	if err != nil {
		tg.logger.Error("phoneLvlHandler: error getting user", log2.Fields{
			"error": err,
		})
	}

	user.Phone = message.Text

	err = tg.svc.UpdateUser(*user)
	if err != nil {
		tg.logger.Error("phoneLvlHandler: error updating user", log2.Fields{
			"error": err,
		})
	}

	user, err = tg.svc.GetUser(userID)
	if err != nil {
		tg.logger.Error("phoneLvlHandler: error getting user", log2.Fields{
			"error": err,
		})
	}

	sendText := fmt.Sprintf("Ваш профиль.\n\nИмя: %s\nАдрес: %s\nТелефон: %s", user.Name, user.Address, user.Phone)
	kb := profileKB()

	err = tg.newMsg(userID, sendText, kb)
	if err != nil {
		tg.logger.Error("phoneLvlHandler: error sending message", log2.Fields{
			"error": err,
		})
	}

	delete(tg.cache[userID], "lvl")
}

func (tg TGBot) nameLvlHandler(message *tgbotapi.Message) {
	userID := message.From.ID

	user, err := tg.svc.GetUser(userID)
	if err != nil {
		tg.logger.Error("nameLvlHandler: error getting user", log2.Fields{
			"error": err,
		})
	}

	user.Name = message.Text

	err = tg.svc.UpdateUser(*user)
	if err != nil {
		tg.logger.Error("nameLvlHandler: error updating user", log2.Fields{
			"error": err,
		})
	}

	user, err = tg.svc.GetUser(userID)
	if err != nil {
		tg.logger.Error("nameLvlHandler: error getting user", log2.Fields{
			"error": err,
		})
	}

	sendText := fmt.Sprintf("Ваш профиль.\n\nИмя: %s\nАдрес: %s\nТелефон: %s", user.Name, user.Address, user.Phone)
	kb := profileKB()

	err = tg.newMsg(userID, sendText, kb)
	if err != nil {
		tg.logger.Error("nameLvlHandler: error sending message", log2.Fields{
			"error": err,
		})
	}

	delete(tg.cache[userID], "lvl")
}

func (tg TGBot) addressLvlHandler(message *tgbotapi.Message) {
	userID := message.From.ID

	user, err := tg.svc.GetUser(userID)
	if err != nil {
		tg.logger.Error("addressLvlHandler: error getting user", log2.Fields{
			"error": err,
		})
		return
	}

	user.Address = message.Text

	err = tg.svc.UpdateUser(*user)
	if err != nil {
		tg.logger.Error("addressLvlHandler: error updating user", log2.Fields{
			"error": err,
		})
		return
	}

	user, err = tg.svc.GetUser(userID)
	if err != nil {
		tg.logger.Error("addressLvlHandler: error getting user", log2.Fields{
			"error": err,
		})
		return
	}

	sendText := fmt.Sprintf("Ваш профиль.\n\nИмя: %s\nАдрес: %s\n Телефон: %s", user.Name, user.Address, user.Phone)
	kb := profileKB()

	err = tg.newMsg(userID, sendText, kb)
	if err != nil {
		tg.logger.Error("addressLvlHandler: error sending message", log2.Fields{
			"error": err,
		})
		return
	}

	delete(tg.cache[userID], "lvl")
}

func (tg TGBot) startCommandHandler(message *tgbotapi.Message) {
	userID := message.From.ID

	err := tg.svc.NewUser(entities.User{
		ID:      message.From.ID,
		Name:    " ",
		Address: " ",
		Phone:   " ",
	})
	if err != nil {
		tg.logger.Error("startCommandHandler: error creating new user", log2.Fields{
			"error": err,
		})
		return
	}

	err = tg.newMsg(userID, startText, newStartKB())
	if err != nil {
		tg.logger.Error("startCommandHandler: error sending message", log2.Fields{
			"error": err,
		})
		return
	}
}

func (tg TGBot) currentCommandHandler(message *tgbotapi.Message) {
	userID := message.From.ID

	orders, err := tg.svc.GetAllCurrentOrders()
	if err != nil {
		tg.logger.Error("currentCommandHandler: error getting current orders", log2.Fields{
			"error": err,
		})
		return
	}

	for _, order := range orders {
		user, err := tg.svc.GetUser(order.UserID)
		if err != nil {
			tg.logger.Error("currentCommandHandler: error getting user", log2.Fields{
				"error": err,
			})
			return
		}

		var productText string

		for _, product := range order.Composition {
			productText += fmt.Sprintf("\n\nТип: %s\nРазмер: %s\nЦвет %s\nТекст: %s\nФото: %s\nКол-во: %d",
				product.Name,
				product.Size,
				product.Color,
				product.Text,
				product.Img,
				product.Amount,
			)
		}

		sendText := fmt.Sprintf("Имя: %s\nАдрес: %s\nТелефон: %s\nДата заказа: %s%v\n\n/done_%d", user.Name, user.Address, user.Phone, fmt.Sprintf("%d.%d.%d", order.Start.Day(), order.Start.Month(), order.Start.Year()), productText, order.ID)

		err = tg.newMsg(userID, sendText, nil)
		if err != nil {
			tg.logger.Error("currentCommandHandler: error sending message", log2.Fields{
				"error": err,
			})
			return
		}
	}
}

func (tg TGBot) doneCommandHandler(message *tgbotapi.Message) {
	userID := message.From.ID

	orders, err := tg.svc.GetAllDoneOrders()
	if err != nil {
		tg.logger.Error("doneCommandHandler: error getting all done orders", log2.Fields{
			"error": err,
		})
		return
	}

	for _, order := range orders {
		user, err := tg.svc.GetUser(order.UserID)
		if err != nil {
			tg.logger.Error("doneCommandHandler: error getting user", log2.Fields{
				"error": err,
			})
			return
		}

		var productText string

		for _, product := range order.Composition {
			productText += fmt.Sprintf("\n\nТип: %s\nРазмер: %s\nЦвет %s\nТекст: %s\nФото: %s\nКол-во: %d",
				product.Name,
				product.Size,
				product.Color,
				product.Text,
				product.Img,
				product.Amount,
			)
		}

		sendText := fmt.Sprintf("Имя: %s\nАдрес: %s\n\nЗаказ:\nДата создания заказа: %s\nДата выполнения заказа: %s%s", user.Name, user.Address, fmt.Sprintf("%d.%d.%d", order.Start.Day(), order.Start.Month(), order.Start.Year()), fmt.Sprintf("%d.%d.%d", order.Done.Day(), order.Done.Month(), order.Done.Year()), productText)

		err = tg.newMsg(userID, sendText, nil)
		if err != nil {
			tg.logger.Error("doneCommandHandler: error sending message", log2.Fields{
				"error": err,
			})
			return
		}
	}
}

func (tg TGBot) done_xCommandHandler(split []string) {
	orderID, err := strconv.Atoi(split[1])
	if err != nil {
		tg.logger.Error("done_xCommandHandler: error converting string to int", log2.Fields{
			"error": err,
		})
		return
	}

	err = tg.svc.NewDoneOrder(int64(orderID))
	if err != nil {
		tg.logger.Error("done_xCommandHandler: error creating done order", log2.Fields{
			"error": err,
		})
		return
	}
}

func (tg TGBot) printLvlHandler(update tgbotapi.Update) {
	userID := update.Message.From.ID

	if update.Message.Document != nil {
		// writing file

		fileID := update.Message.Document.FileID
		fileURL, err := tg.bot.GetFileDirectURL(fileID)
		if err != nil {
			tg.logger.Error("printLvlHandler: getting file url for document error", log2.Fields{
				"error": err,
			})
			return
		}

		// adding product to cart
		tg.cache[userID]["newProd"].(*entities.Product).Img = fileURL
	}

	if update.Message.Text != "" {
		tg.cache[userID]["newProd"].(*entities.Product).Text = update.Message.Text
	}

	if update.Message.Photo != nil {
		fileID := update.Message.Photo[len(update.Message.Photo)-1].FileID

		fileURL, err := tg.bot.GetFileDirectURL(fileID)
		if err != nil {
			tg.logger.Error("printLvlHandler: getting file url for photo error", log2.Fields{
				"error": err,
			})
			return
		}

		tg.cache[userID]["newProd"].(*entities.Product).Img = fileURL
	}

	tg.cache[userID]["newProd"].(*entities.Product).Amount = 1

	// indexing products in cart for each user
	if _, ok := tg.cache[userID]["prodIdx"].(int); !ok {
		tg.cache[userID]["prodIdx"] = 1
	} else {
		tg.cache[userID]["prodIdx"] = tg.cache[userID]["prodIdx"].(int) + 1
	}

	// adding product to cart
	err := tg.svc.NewCartProduct(userID, tg.cache[userID]["prodIdx"].(int), *tg.cache[userID]["newProd"].(*entities.Product))
	if err != nil {
		tg.logger.Error("printLvlHandler: new product in cart error", log2.Fields{
			"error": err,
		})
		return
	}

	err = tg.newMsg(userID, addingToCartText, backAndStartKB())
	if err != nil {
		tg.logger.Error("printLvlHandler: new msg procedure failed", log2.Fields{
			"error": err,
		})
		return
	}

	delete(tg.cache[userID], "newProd")
	delete(tg.cache[userID], "lvl")
}
