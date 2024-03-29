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
	data := update.CallbackData()
	split := strings.Split(data, "/")

	switch split[0] {
	case "info":
		tg.infoHandler(update.CallbackQuery)
	case "profile":
		if len(split) == 1 {
			tg.getProfileHandler(update.CallbackQuery)
		}

		if len(split) == 2 {
			tg.changeProfileHandler(update.CallbackQuery, data)
		}
	case "cart":
		if len(split) == 1 {
			tg.startCartHandler(update.CallbackQuery)
		}

		if len(split) == 2 {
			switch split[1] {
			case "delete_prod":
				tg.deleteProductFromCartHandler(update.CallbackQuery)
			case "increase":
				tg.increaseProductAmountHandler(update.CallbackQuery)
			case "decrease":
				tg.decreaseProductAmountHandler(update.CallbackQuery)
			case "right":
				tg.moveCartToRightHandler(update.CallbackQuery)
			case "left":
				tg.moveCartToLeftHandler(update.CallbackQuery)
			}
		}
	case "home":
		tg.backHomeHandler(update.CallbackQuery)
	case "start_shopping":
		switch len(split) {
		case 1:
			tg.startShoppingHandler(update.CallbackQuery, data)
		case 2:
			tg.productNameHandler(update.CallbackQuery, data)
		case 3:
			tg.productColorHandler(update.CallbackQuery, data)
		case 4:
			tg.productSizeHandler(update.CallbackQuery, data)
		}
	case "design":
		tg.designOrderHandler(update.CallbackQuery)
	default:
		return
	}
}

func (tg TGBot) infoHandler(callback *tgbotapi.CallbackQuery) {
	sendText := infoText
	kb := backToStartKB()
	userID := callback.From.ID

	err := tg.newEditMsgByDelete(userID, sendText, kb)
	if err != nil {
		tg.logger.Error("infoHandler: new edit msg procedure failed", log2.Fields{
			"error": err,
		})
	}
}

func (tg TGBot) getProfileHandler(callback *tgbotapi.CallbackQuery) {
	userID := callback.From.ID

	user, err := tg.svc.GetUser(userID)
	if err != nil {
		tg.logger.Error("getProfileHandler: getting user failed", log2.Fields{
			"error": err,
		})
		return
	}

	sendText := fmt.Sprintf("Ваш профиль.\n\nИмя: %s\nАдрес: %s\nТелефон: %s", user.Name, user.Address, user.Phone)
	kb := profileKB()

	err = tg.newEditMsgByDelete(userID, sendText, kb)
	if err != nil {
		tg.logger.Error("getProfileHandler: new edit msg procedure failed", log2.Fields{
			"error": err,
		})
		return
	}
}

func (tg TGBot) changeProfileHandler(callback *tgbotapi.CallbackQuery, data string) {
	split := strings.Split(data, "/")

	userID := callback.From.ID

	var sendText string
	switch split[1] {
	case "phone":
		tg.cache[userID]["lvl"] = "phone"

		sendText = newPhoneText
	case "name":
		tg.cache[userID]["lvl"] = "name"

		sendText = newNameText
	case "address":
		tg.cache[userID]["lvl"] = "address"

		sendText = newAddressText
	}

	err := tg.newMsg(userID, sendText, nil)
	if err != nil {
		tg.logger.Error("changeProfileHandler: new edit msg procedure failed", log2.Fields{
			"error": err,
		})
		return
	}
}

func (tg TGBot) startCartHandler(callback *tgbotapi.CallbackQuery) {
	userID := callback.From.ID

	var sendText string
	var kb *tgbotapi.InlineKeyboardMarkup

	cartLen, err := tg.svc.CartLen(userID)
	if err != nil {
		tg.logger.Error("startCartHandler: getting cart len failed", log2.Fields{
			"error": err,
		})
		return
	}

	if *cartLen == 0 {
		tg.newAlert(callback.ID, emptyCartText)
		return
	}

	cart, err := tg.svc.GetCart(userID)
	if err != nil {
		tg.logger.Error("startCartHandler: getting cart failed", log2.Fields{
			"error": err,
		})
		return
	}

	tg.cache[userID]["keys"] = make([]int, 0, len(cart))

	for key := range cart {
		tg.cache[userID]["keys"] = append(tg.cache[userID]["keys"].([]int), key)
	}

	tg.cache[userID]["idx"] = 0

	keys := tg.cache[userID]["keys"].([]int)
	idx := tg.cache[userID]["idx"].(int)

	currentProd, err := tg.svc.GetCartProduct(userID, keys[idx])
	if err != nil {
		tg.logger.Error("startCartHandler: getting cart product failed", log2.Fields{
			"error": err,
		})
		return
	}

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

	kb = cartKB(currentProd.Amount)

	photoFile := fmt.Sprintf("%s_%s.jpg", currentProd.Color, currentProd.Name)

	err = tg.newEditPhotoByDelete(userID, tgbotapi.FilePath("././././assets/"+photoFile), sendText, kb)
	if err != nil {
		tg.logger.Error("startCartHandler: new edit photo msg procedure failed", log2.Fields{
			"error": err,
		})
		return
	}
}

func (tg TGBot) deleteProductFromCartHandler(callback *tgbotapi.CallbackQuery) {
	userID := callback.From.ID

	var sendText string
	var kb *tgbotapi.InlineKeyboardMarkup

	keys := tg.cache[userID]["keys"].([]int)
	idx := tg.cache[userID]["idx"].(int)

	tg.svc.DeleteProductFromCart(userID, keys[idx])

	cartLen, err := tg.svc.CartLen(userID)
	if err != nil {
		tg.logger.Error("deleteProductFromCartHandler: getting cart len failed", log2.Fields{
			"error": err,
		})
		return
	}

	if *cartLen == 0 {
		sendText = startText
		kb = newStartKB()

		err := tg.newEditMsgByDelete(userID, sendText, kb)
		if err != nil {
			tg.logger.Error("deleteProductFromCartHandler: edit msg procedure failed", log2.Fields{
				"error": err,
			})
		}

		return
	}

	cart, err := tg.svc.GetCart(userID)
	if err != nil {
		tg.logger.Error("deleteProductFromCartHandler: getting cart failed", log2.Fields{
			"error": err,
		})
	}

	tg.cache[userID]["keys"] = make([]int, 0, len(cart))

	for key := range cart {
		tg.cache[userID]["keys"] = append(tg.cache[userID]["keys"].([]int), key)
	}

	tg.cache[userID]["idx"] = 0

	keys = tg.cache[userID]["keys"].([]int)
	idx = tg.cache[userID]["idx"].(int)

	currentProd, err := tg.svc.GetCartProduct(userID, keys[idx])
	if err != nil {
		tg.logger.Error("deleteProductFromCartHandler: getting cart product failed", log2.Fields{
			"error": err,
		})
		return
	}

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

	kb = cartKB(currentProd.Amount)

	photoFile := fmt.Sprintf("%s_%s.jpg", currentProd.Color, currentProd.Name)

	err = tg.newEditPhotoByDelete(userID, tgbotapi.FilePath("././././assets/"+photoFile), sendText, kb)
	if err != nil {
		tg.logger.Error("deleteProductFromCartHandler: new edit photo msg procedure failed", log2.Fields{
			"error": err,
		})
		return
	}
}

func (tg TGBot) increaseProductAmountHandler(callback *tgbotapi.CallbackQuery) {
	userID := callback.From.ID

	keys := tg.cache[userID]["keys"].([]int)
	idx := tg.cache[userID]["idx"].(int)

	currentProd, err := tg.svc.GetCartProduct(userID, keys[idx])
	if err != nil {
		tg.logger.Error("increaseProductAmountHandler: getting cart product failed", log2.Fields{
			"error": err,
		})
		return
	}

	currentProd.Amount++

	err = tg.svc.NewCartProduct(userID, keys[idx], *currentProd)
	if err != nil {
		tg.logger.Error("increaseProductAmountHandler: new cart product failed", log2.Fields{
			"error": err,
		})
		return
	}

	kb := cartKB(currentProd.Amount)

	if err = tg.newEditMsgKeyboard(userID, kb); err != nil {
		tg.logger.Error("increaseProductAmountHandler: new edit msg keyboard procedure failed", log2.Fields{
			"error": err,
		})
		return
	}
}

func (tg TGBot) decreaseProductAmountHandler(callback *tgbotapi.CallbackQuery) {
	userID := callback.From.ID

	keys := tg.cache[userID]["keys"].([]int)
	idx := tg.cache[userID]["idx"].(int)

	currentProd, err := tg.svc.GetCartProduct(userID, keys[idx])
	if err != nil {
		tg.logger.Error("decreaseProductAmountHandler: getting cart product failed", log2.Fields{
			"error": err,
		})
		return
	}

	if currentProd.Amount == 1 {
		return
	}

	currentProd.Amount--

	err = tg.svc.NewCartProduct(userID, keys[idx], *currentProd)
	if err != nil {
		tg.logger.Error("decreaseProductAmountHandler: new cart product failed", log2.Fields{
			"error": err,
		})
		return
	}

	kb := cartKB(currentProd.Amount)

	if err = tg.newEditMsgKeyboard(userID, kb); err != nil {
		tg.logger.Error("decreaseProductAmountHandler: new edit msg keyboard procedure failed", log2.Fields{
			"error": err,
		})
		return
	}
}

func (tg TGBot) moveCartToRightHandler(callback *tgbotapi.CallbackQuery) {
	userID := callback.From.ID

	var sendText string
	var kb *tgbotapi.InlineKeyboardMarkup

	cartLen, err := tg.svc.CartLen(userID)
	if err != nil {
		tg.logger.Error("moveCartToRightHandler: getting cart len failed", log2.Fields{
			"error": err,
		})
		return
	}

	if *cartLen == 1 {
		return
	}

	if tg.cache[userID]["idx"] == int(*cartLen)-1 {
		tg.cache[userID]["idx"] = 0
	} else {
		tg.cache[userID]["idx"] = tg.cache[userID]["idx"].(int) + 1
	}

	keys := tg.cache[userID]["keys"].([]int)
	idx := tg.cache[userID]["idx"].(int)

	currentProd, err := tg.svc.GetCartProduct(userID, keys[idx])
	if err != nil {
		tg.logger.Error("moveCartToRightHandler: getting cart product failed", log2.Fields{
			"error": err,
		})
		return
	}

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

	kb = cartKB(currentProd.Amount)

	photoFile := fmt.Sprintf("%s_%s.jpg", currentProd.Color, currentProd.Name)

	err = tg.newEditPhotoByDelete(userID, tgbotapi.FilePath("././././assets/"+photoFile), sendText, kb)
	if err != nil {
		tg.logger.Error("moveCartToRightHandler: new edit photo msg procedure failed", log2.Fields{
			"error": err,
		})
		return
	}
}

func (tg TGBot) moveCartToLeftHandler(callback *tgbotapi.CallbackQuery) {
	userID := callback.From.ID

	var sendText string
	var kb *tgbotapi.InlineKeyboardMarkup

	cartLen, err := tg.svc.CartLen(userID)
	if err != nil {
		tg.logger.Error("moveCartToLeftHandler: getting cart len failed", log2.Fields{
			"error": err,
		})
		return
	}

	if *cartLen == 1 {
		return
	}

	if tg.cache[userID]["idx"] == 0 {
		tg.cache[userID]["idx"] = int(*cartLen) - 1
	} else {
		tg.cache[userID]["idx"] = tg.cache[userID]["idx"].(int) - 1
	}

	keys := tg.cache[userID]["keys"].([]int)
	idx := tg.cache[userID]["idx"].(int)

	currentProd, err := tg.svc.GetCartProduct(userID, keys[idx])
	if err != nil {
		tg.logger.Error("moveCartToLeftHandler: getting cart product failed", log2.Fields{
			"error": err,
		})
		return
	}

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

	kb = cartKB(currentProd.Amount)

	photoFile := fmt.Sprintf("%s_%s.jpg", currentProd.Color, currentProd.Name)

	err = tg.newEditPhotoByDelete(userID, tgbotapi.FilePath("././././assets/"+photoFile), sendText, kb)
	if err != nil {
		tg.logger.Error("moveCartToLeftHandler: new edit photo msg procedure failed", log2.Fields{
			"error": err,
		})
		return
	}
}

func (tg TGBot) backHomeHandler(callback *tgbotapi.CallbackQuery) {
	userID := callback.From.ID

	sendText := startText
	kb := newStartKB()

	err := tg.newEditMsgByDelete(userID, sendText, kb)
	if err != nil {
		tg.logger.Error("backHomeHandler: new edit msg procedure failed", log2.Fields{
			"error": err,
		})
		return
	}
}

func (tg TGBot) startShoppingHandler(callback *tgbotapi.CallbackQuery, data string) {
	userID := callback.From.ID

	sendText := selectTypeText
	kb := newProdNameKB(data)

	err := tg.newEditMsgByDelete(userID, sendText, kb)
	if err != nil {
		tg.logger.Error("startShoppingHandler: new edit msg procedure failed", log2.Fields{
			"error": err,
		})
		return
	}
}

func (tg TGBot) productNameHandler(callback *tgbotapi.CallbackQuery, data string) {
	split := strings.Split(data, "/")

	userID := callback.From.ID

	var photoFile string

	switch split[1] {
	case "толстовка":
		tg.cache[userID]["newProd"] = &entities.Product{Name: "толстовка"}
		photoFile = "белый_толстовка.jpg"
	case "штаны":
		tg.cache[userID]["newProd"] = &entities.Product{Name: "штаны"}
		photoFile = "белый_штаны.jpg"
	}

	sendText := selectColorText
	kb := newProdColorKB(data)

	err := tg.newEditPhotoByDelete(userID, tgbotapi.FilePath("././././assets/"+photoFile), sendText, kb)
	if err != nil {
		tg.logger.Error("productNameHandler: new edit photo msg procedure failed", log2.Fields{
			"error": err,
		})
		return
	}
}

func (tg TGBot) productColorHandler(callback *tgbotapi.CallbackQuery, data string) {
	split := strings.Split(data, "/")

	userID := callback.From.ID

	tg.cache[userID]["newProd"].(*entities.Product).Color = split[2]

	photoFile := fmt.Sprintf("%s_%s.jpg", split[2], split[1])
	sendText := selectSizeText
	kb := newProdSizeKB(data)

	err := tg.newEditPhotoByDelete(userID, tgbotapi.FilePath("././././assets/"+photoFile), sendText, kb)
	if err != nil {
		tg.logger.Error("productColorHandler: new edit photo msg procedure failed", log2.Fields{
			"error": err,
		})
		return
	}
}

func (tg TGBot) productSizeHandler(callback *tgbotapi.CallbackQuery, data string) {
	split := strings.Split(data, "/")

	userID := callback.From.ID

	tg.cache[userID]["newProd"].(*entities.Product).Size = split[3]

	tg.cache[userID]["lvl"] = "print"

	photoFile := fmt.Sprintf("%s_%s.jpg", split[2], split[1])
	sendText := customPrintText

	err := tg.newEditPhotoByDelete(userID, tgbotapi.FilePath("././././assets/"+photoFile), sendText, nil)
	if err != nil {
		tg.logger.Error("productSizeHandler: new edit photo msg procedure failed", log2.Fields{
			"error": err,
		})
		return
	}
}

func (tg TGBot) designOrderHandler(callback *tgbotapi.CallbackQuery) {
	userID := callback.From.ID

	// check user
	user, err := tg.svc.GetUser(userID)
	if err != nil {
		tg.logger.Error("designOrderHandler: getting user failed", log2.Fields{
			"error": err,
		})
		return
	}

	if user.Name == "" || user.Address == "" || user.Phone == "" {
		tg.getProfileHandler(callback)
		return
	}

	// creating order
	cart, err := tg.svc.GetCart(userID)
	if err != nil {
		tg.logger.Error("designOrderHandler: getting cart failed", log2.Fields{
			"error": err,
		})
		return
	}

	cartProducts := make([]entities.Product, 0, len(cart))

	for _, val := range cart {
		cartProducts = append(cartProducts, val)
	}

	order := entities.CurrentOrder{
		UserID:      userID,
		Composition: cartProducts,
		Start:       time.Now(),
	}

	tg.cache[userID]["idx"] = 0

	err = tg.svc.NewCurrentOrder(order)
	if err != nil {
		log.Println(err)
	}

	sendText := createOrderText
	kb := backToStartKB()

	err = tg.newEditMsgByDelete(userID, sendText, kb)
	if err != nil {
		tg.logger.Error("designOrderHandler: new edit msg procedure failed", log2.Fields{
			"error": err,
		})
		return
	}

	for _, admin := range tg.conf.Admins {
		err = tg.newMsg(admin, "Поступил новый заказ!", nil)
		if err != nil {
			tg.logger.Error("designOrderHandler: new msg procedure failed", log2.Fields{
				"error": err,
			})
			return
		}
	}
}
