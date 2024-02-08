package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (tg TGBot) newMsg(userID int64, text string, kb *tgbotapi.InlineKeyboardMarkup) error {
	msg := tgbotapi.NewMessage(userID, text)
	msg.Text = text
	msg.ReplyMarkup = kb

	sentMsg, err := tg.bot.Send(msg)
	if err != nil {
		return err
	}

	if _, ok := tg.cache[userID]; !ok {
		tg.cache[userID] = make(map[string]interface{})
	}

	tg.cache[userID]["msgID"] = sentMsg.MessageID

	return nil
}

func (tg TGBot) newEditMsg(userID int64, msgID int, text string, kb *tgbotapi.InlineKeyboardMarkup) error {
	updatedMsg := tgbotapi.NewEditMessageText(userID, msgID, text)
	updatedMsg.ReplyMarkup = kb

	_, err := tg.bot.Send(updatedMsg)
	if err != nil {
		return err
	}

	return nil
}

func (tg TGBot) NewAlert(id string, text string) {
	alert := tgbotapi.NewCallbackWithAlert(id, text)
	_, _ = tg.bot.Send(alert)
}
