package telegram

import (
	"bot/internal/entities"
	log2 "bot/internal/log"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func (tg TGBot) HandleMessage(update tgbotapi.Update) {
	userID := update.Message.From.ID

	switch tg.cache[userID]["lvl"] {
	case "name":
		user := tg.svc.GetUser(userID)

		user.Name = update.Message.Text

		tg.svc.UpdateUser(*user)

		delete(tg.cache[userID], "lvl")
	case "address":
		user := tg.svc.GetUser(userID)

		user.Address = update.Message.Text

		tg.svc.UpdateUser(*user)

		delete(tg.cache[userID], "lvl")
	case "print":
		if update.Message.Document != nil {
			// writing file

			fileID := update.Message.Document.FileID
			fileURL, err := tg.bot.GetFileDirectURL(fileID)
			if err != nil {
				tg.logger.Error("getting file url error", log2.Fields{
					"error": err,
				})
			}

			response, err := http.Get(fileURL)
			if err != nil {
				tg.logger.Error("get http failed", log2.Fields{
					"error": err,
				})
				return
			}

			defer func() {
				err = response.Body.Close()
				if err != nil {
					tg.logger.Error("closing response body failed", log2.Fields{
						"error": err,
					})
				}
				return
			}()

			output, err := os.Create(filepath.Join("./././filesStorage/", filepath.Base(update.Message.Document.FileName)))
			if err != nil {
				tg.logger.Error("creating file failed", log2.Fields{
					"error": err,
				})
				return
			}

			defer func() {
				err = output.Close()
				if err != nil {
					tg.logger.Error("closing output body failed", log2.Fields{
						"error": err,
					})
					return
				}
			}()

			extension := filepath.Ext(update.Message.Document.FileName)
			fileName := fileID + extension

			err = os.Rename(fmt.Sprintf("./././filesStorage/%s", update.Message.Document.FileName), fileName)
			if err != nil {
				log.Println(err)
			}

			data, err := io.ReadAll(response.Body)
			if err != nil {
				tg.logger.Error("reading response body failed", log2.Fields{
					"error": err,
				})
				return
			}

			_, err = output.Write(data)
			if err != nil {
				tg.logger.Error("writing file failed", log2.Fields{
					"error": err,
				})
				return
			}

			tg.logger.Info("file installed successfully", log2.Fields{
				"file": fileName,
			})

			// adding product to cart
			tg.cache[userID]["newProd"].(*entities.Product).Img = fileName

			fmt.Println(tg.cache[userID]["newProd"].(*entities.Product))
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

		err := tg.newMsg(userID, "Отлично! Вы успешно добавили ваш товар в корзину!", backToStartKB())
		if err != nil {
			tg.logger.Error("new msg procedure failed", log2.Fields{
				"error": err,
			})
		}

		delete(tg.cache[userID], "newProd")
	default:
		switch update.Message.Command() {
		case "start":
			tg.svc.NewUser(entities.User{
				ID:      update.Message.From.ID,
				Name:    "",
				Address: "",
			})

			err := tg.newMsg(userID, "Дарова, бро", newStartKB())
			if err != nil {
				tg.logger.Error("error sending message", log2.Fields{
					"error": err,
				})
				return
			}
		case "curr_orders":
			for _, admin := range tg.conf.Admins {
				if userID == admin {

				}
			}

			return
		}
	}
}
