package telegram

import (
	"bot/internal/log"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (tg TGBot) newMsg(userID int64, text string, kb *tgbotapi.InlineKeyboardMarkup) error {
	msg := tgbotapi.NewMessage(userID, text)
	msg.Text = text
	msg.ReplyMarkup = kb

	sentMsg, err := tg.bot.Send(msg)
	if err != nil {
		tg.logger.Error("newMsg: new msg procedure failed", log.Fields{
			"error": err,
		})

		return err
	}

	if _, ok := tg.cache[userID]; !ok {
		tg.cache[userID] = make(map[string]interface{})
	}

	tg.cache[userID]["msgID"] = sentMsg.MessageID

	return nil
}

func (tg TGBot) newEditMsgText(userID int64, text string, kb *tgbotapi.InlineKeyboardMarkup) error {
	updatedMsg := tgbotapi.NewEditMessageText(userID, tg.cache[userID]["msgID"].(int), text)
	updatedMsg.ReplyMarkup = kb

	_, err := tg.bot.Send(updatedMsg)
	if err != nil {
		return err
	}

	return nil
}

func (tg TGBot) newAlert(id string, text string) {
	alert := tgbotapi.NewCallbackWithAlert(id, text)
	_, _ = tg.bot.Send(alert)
}

func (tg TGBot) newDeleteMsg(userID int64) {
	delMsg := tgbotapi.NewDeleteMessage(userID, tg.cache[userID]["msgID"].(int))

	_, _ = tg.bot.Send(delMsg)
}

func (tg TGBot) newPhotoMsg(userID int64, photo tgbotapi.FilePath, caption string, kb *tgbotapi.InlineKeyboardMarkup) error {
	photoMsg := tgbotapi.NewPhoto(userID, photo)

	photoMsg.Caption = caption
	photoMsg.ReplyMarkup = kb

	sentMsg, err := tg.bot.Send(photoMsg)
	if err != nil {
		tg.logger.Error("newPhotoMsg: error sending photo msg", log.Fields{
			"error": err,
		})

		return err
	}

	if _, ok := tg.cache[userID]; !ok {
		tg.cache[userID] = make(map[string]interface{})
	}

	tg.cache[userID]["msgID"] = sentMsg.MessageID

	return nil
}

func (tg TGBot) newEditMsgByDelete(userID int64, text string, kb *tgbotapi.InlineKeyboardMarkup) error {
	tg.newDeleteMsg(userID)

	if err := tg.newMsg(userID, text, kb); err != nil {
		tg.logger.Error("newEditMsgByDelete: new msg procedure failed", log.Fields{
			"error": err,
		})

		return err
	}

	return nil
}

func (tg TGBot) newEditPhotoByDelete(userID int64, photo tgbotapi.FilePath, caption string, kb *tgbotapi.InlineKeyboardMarkup) error {
	tg.newDeleteMsg(userID)

	if err := tg.newPhotoMsg(userID, photo, caption, kb); err != nil {
		tg.logger.Error("newEditPhotoByDelete: new photo msg procedure failed", log.Fields{
			"error": err,
		})

		return err
	}

	return nil
}

func (tg TGBot) newEditMsgKeyboard(userID int64, kb *tgbotapi.InlineKeyboardMarkup) error {
	updatedMsg := tgbotapi.NewEditMessageReplyMarkup(userID, tg.cache[userID]["msgID"].(int), *kb)

	_, err := tg.bot.Send(updatedMsg)
	if err != nil {
		return err
	}

	return nil
}
