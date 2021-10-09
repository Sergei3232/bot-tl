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
			addNewUserBot(database.GetDB(), update.Message.Chat.ID, update.Message.Chat.UserName)
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, textMessageUser)
			_, _ = bot.Send(msg)
		case "messageUsers":
			listUser := getUserTelegramID(database.GetDB())

			for _, id := range listUser {
				msg := tgbotapi.NewMessage(int64(id), update.Message.CommandArguments())
				bot.Send(msg)
			}
		case "addAdmin":
			err := addAdmin(database.GetDB(), update.Message.Chat.ID)
			if err != nil {
				textMessageUser = "Ошибка добавления пользователя в группу Администрирования!"
			} else {
				textMessageUser = "Пользователь добавлен в группу!"
			}
		case "help":

			thisAdmin := isAnAdmin(database.GetDB(), update.Message.Chat.ID)
			if thisAdmin {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, adminMenu())
				bot.Send(msg)
			}

		default:
			textMessageUser = "Команда не известна!!! Попробуйте задать другую команду!!!"
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, textMessageUser)
			bot.Send(msg)
		}

	}
}

//Returns a list of user ID that interacted with the bot
func getUserTelegramID(db *gorm.DB) []uint {
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

//Adds a unique bot user to the database
func addNewUserBot(db *gorm.DB, id int64, nameUser string) {
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

//Checking the user for administrator rights
func isAnAdmin(db *gorm.DB, id int64) bool {
	var isAdmin bool
	adminUser := database.AdministratorsGroup{}
	db.Where("id = ?", id).First(&adminUser)

	if adminUser.Id != 0 {
		isAdmin = true
	}
	return isAdmin
}

//The text of the help on commands for the administrator
func adminMenu() (menuAdmin string) {
	menuAdmin = `Доступные команды Администратора:
	/messageUsers [Сообщение пользователю]
	Отправка сообщения всем пользователям бота
	/addAdmin [Телеграм id]
	Добавление пользователя в группу администрирования
	/deleteAdmin [Телеграм id]
	Удаляет пользователя из группу администрирования
	/ipUserHistory [Телеграм id]
	Показывает все ip из за просов пользователя`
	return
}

//Adding a user s to the admin group
func addAdmin(db *gorm.DB, id int64) (err error) {
	adminUser := database.AdministratorsGroup{}
	db.Where("id = ?", id).First(&adminUser)

	if adminUser.Id == 0 {
		adminUser = database.AdministratorsGroup{Id: uint(id)}
		result := db.Create(&adminUser)
		if result.Error != nil {
			log.Fatal(result.Error)
			err = nil
		} else {
			err = result.Error
		}
	}
	return
}
