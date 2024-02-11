package telegram

import (
	"bot/internal/entities"
	log2 "bot/internal/log"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"strings"
	"time"
)

func (tg TGBot) HandleCallback(update tgbotapi.Update) {
	userID := update.CallbackQuery.From.ID

	data := update.CallbackData()
	split := strings.Split(data, "/")

	var sendText string
	var kb *tgbotapi.InlineKeyboardMarkup

	switch split[0] {
	case "info":
		sendText = infoText
		kb = backToStartKB()
	case "profile":
		if len(split) == 1 {
			user := tg.svc.GetUser(userID)

			sendText = fmt.Sprintf("Ваш профиль.\n\nИмя: %s\nАдрес: %s", user.Name, user.Address)
			kb = profileKB()
		}

		if len(split) == 2 {
			switch split[1] {
			case "name":
				tg.cache[userID]["lvl"] = "name"

				sendText = newNameText
			case "address":
				tg.cache[userID]["lvl"] = "address"

				sendText = newAddressText
			}
		}
	case "cart":
		if len(split) == 1 {
			if tg.svc.CartLen(userID) == 0 {
				tg.NewAlert(update.CallbackQuery.ID, emptyCartText)
				return
			}

			cart := tg.svc.GetCart(userID)

			tg.cache[userID]["keys"] = make([]int, 0, len(cart))

			for key := range cart {
				tg.cache[userID]["keys"] = append(tg.cache[userID]["keys"].([]int), key)
			}

			tg.cache[userID]["idx"] = 0
		}

		if len(split) == 2 {
			switch split[1] {
			case "delete_prod":
				keys := tg.cache[userID]["keys"].([]int)
				idx := tg.cache[userID]["idx"].(int)

				tg.svc.DeleteProductFromCart(userID, keys[idx])

				if tg.svc.CartLen(userID) == 0 {
					sendText = startText
					kb = newStartKB()

					err := tg.newEditMsg(userID, tg.cache[userID]["msgID"].(int), sendText, kb)
					if err != nil {
						tg.logger.Error("edit msg procedure failed", log2.Fields{
							"error": err,
						})
					}

					return
				}

				cart := tg.svc.GetCart(userID)

				tg.cache[userID]["keys"] = make([]int, 0, len(cart))

				for key := range cart {
					tg.cache[userID]["keys"] = append(tg.cache[userID]["keys"].([]int), key)
				}

				tg.cache[userID]["idx"] = 0
			case "increase":
				keys := tg.cache[userID]["keys"].([]int)
				idx := tg.cache[userID]["idx"].(int)

				currentProd := tg.svc.GetCartProduct(userID, keys[idx])

				currentProd.Amount++

				tg.svc.NewCartProduct(userID, keys[idx], currentProd)
			case "decrease":
				keys := tg.cache[userID]["keys"].([]int)
				idx := tg.cache[userID]["idx"].(int)

				currentProd := tg.svc.GetCartProduct(userID, keys[idx])

				if currentProd.Amount == 1 {
					return
				}

				currentProd.Amount--

				tg.svc.NewCartProduct(userID, keys[idx], currentProd)
			case "right":
				if tg.cache[userID]["idx"] == int(tg.svc.CartLen(userID))-1 {
					tg.cache[userID]["idx"] = 0
				} else {
					tg.cache[userID]["idx"] = tg.cache[userID]["idx"].(int) + 1
				}
			case "left":
				if tg.cache[userID]["idx"] == 0 {
					tg.cache[userID]["idx"] = int(tg.svc.CartLen(userID)) - 1
				} else {
					tg.cache[userID]["idx"] = tg.cache[userID]["idx"].(int) - 1
				}
			}
		}

		keys := tg.cache[userID]["keys"].([]int)
		idx := tg.cache[userID]["idx"].(int)

		currentProd := tg.svc.GetCartProduct(userID, keys[idx])

		if currentProd.Text != "" {
			sendText = fmt.Sprintf("Ваша корзина.\n\n%d из %d\n\nКатегория: %s\nЦвет: %s\nРазмер: %s\nТекст: %s",
				tg.cache[userID]["idx"].(int)+1,
				len(keys),
				currentProd.Name,
				currentProd.Color,
				currentProd.Size,
				currentProd.Text,
			)
		}

		if currentProd.Img != "" {
			sendText = fmt.Sprintf("Ваша корзина.\n\n%d из %d\n\nКатегория: %s\nЦвет: %s\nРазмер: %s\nФото: %s",
				tg.cache[userID]["idx"].(int)+1,
				len(keys),
				currentProd.Name,
				currentProd.Color,
				currentProd.Size,
				currentProd.Img,
			)
		}

		kb = newCartKB(currentProd.Amount)
	case "home":
		sendText = startText
		kb = newStartKB()
	case "start_shopping":
		if len(split) == 1 {
			sendText = selectTypeText

			kb = newProdNameKB(data)
		}

		if len(split) == 2 {
			switch split[1] {
			case "hoodie":
				tg.cache[userID]["newProd"] = &entities.Product{Name: "hoodie"}
			case "trousers":
				tg.cache[userID]["newProd"] = &entities.Product{Name: "trousers"}
			}

			sendText = selectColorText
			kb = newProdColorKB(data)
		}

		if len(split) == 3 {
			tg.cache[userID]["newProd"].(*entities.Product).Color = split[2]

			sendText = selectSizeText
			kb = newProdSizeKB(data)
		}

		if len(split) == 4 {
			tg.cache[userID]["newProd"].(*entities.Product).Size = split[3]

			tg.cache[userID]["lvl"] = "print"
			sendText = customPrintText
		}
	case "design":
		cart := tg.svc.GetCart(userID)

		cartProducts := make([]entities.Product, len(cart))

		for _, val := range cart {
			cartProducts = append(cartProducts, val)
		}

		order := entities.CurrentOrder{
			UserID:      userID,
			Composition: cartProducts,
			Start:       time.Now(),
		}

		tg.cache[userID]["currentIdx"] = 0

		err := tg.svc.NewOrder(order)
		if err != nil {
			log.Println(err)
		}

		sendText = createOrderText
		kb = backToStartKB()
	default:
		return
	}

	err := tg.newEditMsg(userID, tg.cache[userID]["msgID"].(int), sendText, kb)
	if err != nil {
		tg.logger.Error("edit msg procedure failed", log2.Fields{
			"error": err,
		})
	}
}
