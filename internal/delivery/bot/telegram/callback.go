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
	fmt.Println(tg.conf.Admins)
	userID := update.CallbackQuery.From.ID

	data := update.CallbackData()
	split := strings.Split(data, "/")

	var sendText string
	var kb *tgbotapi.InlineKeyboardMarkup

	switch split[0] {
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

				sendText = "Введите ваше новое имя"
			case "address":
				tg.cache[userID]["lvl"] = "address"

				sendText = "Введите ваш новый адрес"
			}
		}
	case "cart":
		var prodAmount int64
		var currentProd entities.Product

		if len(split) == 1 {
			if tg.svc.CartLen(userID) == 0 {
				tg.NewAlert(update.CallbackQuery.ID, "Ваша корзина пуста =)")
				return
			}

			tg.cache[userID]["currentIdx"] = 1

			prodAmount = tg.svc.CartLen(userID)
			currentProd = tg.svc.GetCartProduct(userID, tg.cache[userID]["currentIdx"].(int))
		}

		if len(split) == 2 {
			prodAmount = tg.svc.CartLen(userID)

			switch split[1] {
			case "increase":
				currentProd = tg.svc.GetCartProduct(userID, tg.cache[userID]["currentIdx"].(int))
				currentProd.Amount++

				tg.svc.NewCartProduct(userID, tg.cache[userID]["currentIdx"].(int), currentProd)
			case "decrease":

				currentProd = tg.svc.GetCartProduct(userID, tg.cache[userID]["currentIdx"].(int))
				currentProd.Amount--

				tg.svc.NewCartProduct(userID, tg.cache[userID]["currentIdx"].(int), currentProd)
			case "right":
				if tg.cache[userID]["currentIdx"] == int(prodAmount) {
					tg.cache[userID]["currentIdx"] = 1
				} else {
					tg.cache[userID]["currentIdx"] = tg.cache[userID]["currentIdx"].(int) + 1
				}
			case "left":
				if tg.cache[userID]["currentIdx"] == 1 {
					tg.cache[userID]["currentIdx"] = int(prodAmount)
				} else {
					tg.cache[userID]["currentIdx"] = tg.cache[userID]["currentIdx"].(int) - 1
				}
			}
			currentProd = tg.svc.GetCartProduct(userID, tg.cache[userID]["currentIdx"].(int))
		}

		sendText = fmt.Sprintf("Ваша корзина.\n\n%d из %d\n\nКатегория: %s\nЦвет: %s\nРазмер: %s\n",
			tg.cache[userID]["currentIdx"],
			prodAmount,
			currentProd.Name,
			currentProd.Color,
			currentProd.Size,
		)

		kb = newCartKB(currentProd.Amount)

	case "home":
		sendText = "Дарова, бро!"
		kb = newStartKB()
	case "start_shopping":
		if len(split) == 1 {
			sendText = "Выберите, что вы хотите заказать"

			kb = newProdNameKB(data)
		}

		if len(split) == 2 {
			switch split[1] {
			case "hoodie":
				tg.cache[userID]["newProd"] = &entities.Product{Name: "hoodie"}

				sendText = "Выберите цвет вашей толстовки"
				kb = newProdColorKB(data)
			case "trousers":
				tg.cache[userID]["newProd"] = &entities.Product{Name: "trousers"}

				sendText = "Выберите цвет ваших штанов"
				kb = newProdColorKB(data)
			}
		}

		if len(split) == 3 {
			tg.cache[userID]["newProd"].(*entities.Product).Color = split[2]

			sendText = "Выберите размер"
			kb = newProdSizeKB(data)
		}

		if len(split) == 4 {
			tg.cache[userID]["newProd"].(*entities.Product).Size = split[3]

			tg.cache[userID]["lvl"] = "print"
			sendText = "Введите надпись или пришлите файл с изображением, которым хотите видеть у себя:"
		}
	case "design":
		cart := tg.svc.GetCart(userID)

		order := entities.CurrentOrder{
			UserID:      userID,
			Composition: cart,
			Start:       time.Now(),
		}

		tg.cache[userID]["currentIdx"] = 0

		err := tg.svc.NewOrder(order)
		if err != nil {
			log.Println(err)
		}

		sendText = "Ваш заказ был успешно создан! Мы начнем делать ваш заказ сразу после того, как вы оплатите заказ на этот номер телефона:"
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
