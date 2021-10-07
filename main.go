package main

import (
	"github.com/Sergei3232/bot-tl/database"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	"log"
)

type Users struct {
	Id         uint `gorm:"primaryKey"`
	UserName   string
	TelegramId uint
}

func main() {
	database.Init()

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
		var textMessageUser string
		switch update.Message.Command() {
		case "start":
			textMessageUser = "Приветствую вас на нашем канале!!!"
			AddNewUserBot(database.GetDB(), update.Message.Chat.ID, update.Message.Chat.UserName)
		case "list":
			listUser := GetUserTelegramID(database.GetDB())
			log.Printf("Список %d", listUser)
			//textMessageUser = "Тут будет список"
			textMessageUser = update.Message.CommandArguments()
		default:
			textMessageUser = "Команда не известна!!! Попробуйте задать другую команду!!!"

		}
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, textMessageUser)

		bot.Send(msg)

	}
}

func GetUserTelegramID(db *gorm.DB) []uint {
	users := []Users{}
	result := db.Find(&users)

	if result.Error != nil {
		log.Fatal(result.Error)
	}
	userId := make([]uint, len(users))
	for _, n := range users {
		userId = append(userId, n.TelegramId)
	}
	return userId
}

func AddNewUserBot(db *gorm.DB, id int64, nameUser string) {
	user := Users{}
	db.Where("Telegram_id = ?", id).First(&user)

	if user.Id == 0 {
		user = Users{UserName: nameUser, TelegramId: uint(id)}
		result := db.Create(&user)
		if result.Error != nil {
			log.Fatal(result.Error)
		}
	}
}
