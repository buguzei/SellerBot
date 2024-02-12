package telegram

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strconv"
)

const (
	// msg texts
	startText        = "Дарова, бро"
	addingToCartText = "Отлично! Вы успешно добавили ваш товар в корзину!"
	customPrintText  = "Введите надпись или пришлите файл с изображением, которым хотите видеть у себя:"
	emptyCartText    = "Ваша корзина пуста =)"
	newNameText      = "Введите ваше новое имя"
	newAddressText   = "Введите ваш новый адрес"
	createOrderText  = "Ваш заказ был успешно создан! Мы начнем делать ваш заказ сразу после того, как вы оплатите заказ на этот номер телефона:"
	selectTypeText   = "Выберите, что вы хотите заказать"
	selectSizeText   = "Выберите размер"
	selectColorText  = "Выберите цвет"
	infoText         = "Супермегаважная инфа про нас"

	// from eng to rus
	trousers = "штаны"
	hoodie   = "толстовка"

	black = "Черный"
	white = "Белый"
)

func newStartKB() *tgbotapi.InlineKeyboardMarkup {
	kb := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("К покупкам!", "start_shopping"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Ваш профиль", "profile"),
			tgbotapi.NewInlineKeyboardButtonData("Корзина", "cart"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("О нас", "info"),
			tgbotapi.NewInlineKeyboardButtonURL("Поддержка", "https://t.me/buguzei"),
		),
	)

	return &kb
}

func newProdNameKB(path string) *tgbotapi.InlineKeyboardMarkup {
	kb := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Штаны", fmt.Sprintf("%s/trousers", path)),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Толстовка", fmt.Sprintf("%s/hoodie", path)),
		),
	)

	return &kb
}

func newProdColorKB(path string) *tgbotapi.InlineKeyboardMarkup {
	kb := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Белый", fmt.Sprintf("%s/white", path)),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Черный", fmt.Sprintf("%s/black", path)),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Синий", fmt.Sprintf("%s/blue", path)),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Серый", fmt.Sprintf("%s/grey", path)),
		),
	)

	return &kb
}

func newProdSizeKB(path string) *tgbotapi.InlineKeyboardMarkup {
	kb := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("XS", fmt.Sprintf("%s/XS", path)),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("S", fmt.Sprintf("%s/S", path)),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("M", fmt.Sprintf("%s/M", path)),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("L", fmt.Sprintf("%s/L", path)),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("XL", fmt.Sprintf("%s/XL", path)),
		),
	)
	return &kb
}

func newCartKB(amount int) *tgbotapi.InlineKeyboardMarkup {
	strAmount := strconv.Itoa(amount)

	kb := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("<--", "cart/left"),
			tgbotapi.NewInlineKeyboardButtonData("-", "cart/decrease"),
			tgbotapi.NewInlineKeyboardButtonData(strAmount, " "),
			tgbotapi.NewInlineKeyboardButtonData("+", "cart/increase"),
			tgbotapi.NewInlineKeyboardButtonData("-->", "cart/right"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Удалить товар", "cart/delete_prod"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("В главное меню", "home"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Оформить заказ", "design"),
		),
	)
	return &kb
}

func backToStartKB() *tgbotapi.InlineKeyboardMarkup {
	kb := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("В главное меню", "home"),
		),
	)
	return &kb
}

func profileKB() *tgbotapi.InlineKeyboardMarkup {
	kb := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Изменить имя", "profile/name"),
			tgbotapi.NewInlineKeyboardButtonData("Изменить адрес", "profile/address"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("В главное меню", "home"),
		),
	)
	return &kb
}
