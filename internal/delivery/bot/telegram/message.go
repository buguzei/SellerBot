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

func (tg TGBot) nameLvlHandler(message *tgbotapi.Message) {
	userID := message.From.ID

	user := tg.svc.GetUser(userID)

	user.Name = message.Text

	tg.svc.UpdateUser(*user)

	user = tg.svc.GetUser(userID)

	sendText := fmt.Sprintf("Ваш профиль.\n\nИмя: %s\nАдрес: %s", user.Name, user.Address)
	kb := profileKB()

	err := tg.newMsg(userID, sendText, kb)
	if err != nil {
		tg.logger.Error("nameLvlHandler: error sending message", log2.Fields{
			"error": err,
		})
		return
	}

	delete(tg.cache[userID], "lvl")
}

func (tg TGBot) addressLvlHandler(message *tgbotapi.Message) {
	userID := message.From.ID

	user := tg.svc.GetUser(userID)

	user.Address = message.Text

	tg.svc.UpdateUser(*user)

	user = tg.svc.GetUser(userID)

	sendText := fmt.Sprintf("Ваш профиль.\n\nИмя: %s\nАдрес: %s", user.Name, user.Address)
	kb := profileKB()

	err := tg.newMsg(userID, sendText, kb)
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

	tg.svc.NewUser(entities.User{
		ID:      message.From.ID,
		Name:    "",
		Address: "",
	})

	err := tg.newMsg(userID, startText, newStartKB())
	if err != nil {
		tg.logger.Error("startCommandHandler: error sending message", log2.Fields{
			"error": err,
		})
		return
	}
}

func (tg TGBot) currentCommandHandler(message *tgbotapi.Message) {
	userID := message.From.ID

	orders := tg.svc.GetAllCurrentOrders()

	for _, order := range orders {

		user := tg.svc.GetUser(order.UserID)

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

		sendText := fmt.Sprintf("Имя: %s\nАдрес: %s\n\nЗаказ:\nДата заказа: %v%s\n\n/done_%d", user.Name, user.Address, order.Start, productText, order.ID)

		err := tg.newMsg(userID, sendText, nil)
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

	orders := tg.svc.GetAllDoneOrders()

	for _, order := range orders {
		user := tg.svc.GetUser(order.UserID)

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

		sendText := fmt.Sprintf("Имя: %s\nАдрес: %s\n\nЗаказ:\nДата создания заказа: %v\nДата выполнения заказа: %v%s", user.Name, user.Address, order.Start, order.Done, productText)

		err := tg.newMsg(userID, sendText, nil)
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
	}

	tg.svc.NewDoneOrder(int64(orderID))
}

func (tg TGBot) printLvlHandler(update tgbotapi.Update) {
	userID := update.Message.From.ID

	if update.Message.Document != nil {
		// writing file

		fileID := update.Message.Document.FileID
		fileURL, err := tg.bot.GetFileDirectURL(fileID)
		if err != nil {
			tg.logger.Error("printLvlHandler: getting file url error", log2.Fields{
				"error": err,
			})
		}

		// adding product to cart
		tg.cache[userID]["newProd"].(*entities.Product).Img = fileURL
	}

	if update.Message.Text != "" {
		tg.cache[userID]["newProd"].(*entities.Product).Text = update.Message.Text
	}

	tg.cache[userID]["newProd"].(*entities.Product).Amount = 1

	// indexing products in cart for each user
	if _, ok := tg.cache[userID]["prodIdx"].(int); !ok {
		tg.cache[userID]["prodIdx"] = 1
	} else {
		tg.cache[userID]["prodIdx"] = tg.cache[userID]["prodIdx"].(int) + 1
	}

	// adding product to cart
	tg.svc.NewCartProduct(userID, tg.cache[userID]["prodIdx"].(int), *tg.cache[userID]["newProd"].(*entities.Product))

	err := tg.newMsg(userID, addingToCartText, backToStartKB())
	if err != nil {
		tg.logger.Error("printLvlHandler: new msg procedure failed", log2.Fields{
			"error": err,
		})
	}

	delete(tg.cache[userID], "newProd")
	delete(tg.cache[userID], "lvl")
}
