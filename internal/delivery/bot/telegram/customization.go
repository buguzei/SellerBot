package telegram

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strconv"
)

// here is text for bot's messages

const (
	startText        = "Привет! Это официальный бот бренда \"udobno\". Здесь ты можешь быстро и удобно заказать себе наши товары с вашим нанесением и бесплатной доставкой по всей России 🤗"
	addingToCartText = "Отлично! Вы успешно добавили ваш товар в корзину!"
	customPrintText  = "Введите надпись или пришлите изображение, которое хотите видеть у себя на изделии:"
	emptyCartText    = "Ваша корзина пуста =("
	newNameText      = "Введите имя:"
	newAddressText   = "Введите адрес:"
	newPhoneText     = "Введите номер телефона:"
	createOrderText  = "Ваш заказ был успешно создан! В комментарии оплаты УКАЖИТЕ ВАШ НОМЕР ТЕЛЕФОНА, который указан у вас в профиле. Мы начнем делать ваш заказ сразу после того, как вы оплатите заказ на этот номер карты: 2200 7008 7190 5906"
	selectTypeText   = "Выберите, что вы хотите заказать"
	selectSizeText   = "Выберите размер"
	selectColorText  = "Выберите цвет"
	infoText         = "Супермегаважная инфа про нас"
)

// here is keyboards for bot's messages

func newStartKB() *tgbotapi.InlineKeyboardMarkup {
	kb := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("К покупкам!", "start_shopping"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Ваш профиль 👤", "profile"),
			tgbotapi.NewInlineKeyboardButtonData("Корзина 🛍", "cart"),
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
			tgbotapi.NewInlineKeyboardButtonData("Брюки 1990 руб.", fmt.Sprintf("%s/штаны", path)),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Толстовка 1990 руб.", fmt.Sprintf("%s/толстовка", path)),
		),
	)

	return &kb
}

func newProdColorKB(path string) *tgbotapi.InlineKeyboardMarkup {
	kb := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Белый", fmt.Sprintf("%s/белый", path)),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Черный", fmt.Sprintf("%s/черный", path)),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Синий", fmt.Sprintf("%s/синий", path)),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Серый", fmt.Sprintf("%s/серый", path)),
		),
	)

	return &kb
}

func newProdSizeKB(path string) *tgbotapi.InlineKeyboardMarkup {
	kb := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("XS", fmt.Sprintf("%s/XS", path)),
			tgbotapi.NewInlineKeyboardButtonData("S", fmt.Sprintf("%s/S", path)),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("M", fmt.Sprintf("%s/M", path)),
			tgbotapi.NewInlineKeyboardButtonData("L", fmt.Sprintf("%s/L", path)),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("XL", fmt.Sprintf("%s/XL", path)),
			tgbotapi.NewInlineKeyboardButtonData("XXL", fmt.Sprintf("%s/XXL", path)),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("3XL", fmt.Sprintf("%s/3XL", path)),
			tgbotapi.NewInlineKeyboardButtonData("4XL", fmt.Sprintf("%s/4XL", path)),
		),
	)
	return &kb
}

func cartKB(amount int) *tgbotapi.InlineKeyboardMarkup {
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

func backAndStartKB() *tgbotapi.InlineKeyboardMarkup {
	kb := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("В главное меню", "home"),
			tgbotapi.NewInlineKeyboardButtonData("Корзина 🛍", "cart"),
		),
	)
	return &kb
}

func profileKB() *tgbotapi.InlineKeyboardMarkup {
	kb := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Имя", "profile/name"),
			tgbotapi.NewInlineKeyboardButtonData("Адрес", "profile/address"),
			tgbotapi.NewInlineKeyboardButtonData("Телефон", "profile/phone"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Корзина 🛍", "cart"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("В главное меню", "home"),
		),
	)
	return &kb
}
