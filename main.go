package main

import (
	"github.com/Sergei3232/bot-tl/database"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	_ "github.com/lib/pq"
	"log"
)

func main() {
	database.Init()
	//database.GetUserList
	//db := database.GetDB()
	bot, err := tgbotapi.NewBotAPI("2044118489:AAFf-i_MyU4vz14oovc8MEkyPd-5qelnJSY")
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}
		if update.Message.Text == "/list" {

			//database.GetUserList()
			//listUser := database.GetUserList()
			//log.Println(listUser)
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Тут будет список данных")
			msg.ReplyToMessageID = update.Message.MessageID
			bot.Send(msg)
		}
		//log.Printf("Сообщение из логов: [%s] %s", update.Message.From.UserName, update.Message.Text)
		//log.Printf("Команда %s", update.Message.Text)
		//msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
		//msg.ReplyToMessageID = update.Message.MessageID

	}
}
