package main

import (
	"github.com/Sergei3232/bot-tl/database"
	"github.com/Sergei3232/bot-tl/pkg/model"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
	"log"
	"strconv"
)

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
		userAdmin := isAnAdmin(database.GetDB(), update.Message.Chat.ID)

		commandMessage := update.Message.Command()
		switch {
		case commandMessage == "start":
			textMessageUser = `Приветствую вас на нашем канале!!!
Для получения справки по командам введите /help`

			addNewUserBot(database.GetDB(), update.Message.Chat.ID, update.Message.Chat.UserName)
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, textMessageUser)
			bot.Send(msg)
		case userAdmin && "messageUsers" == commandMessage:
			listUser := getUserTelegramID(database.GetDB())
			for _, id := range listUser {
				msg := tgbotapi.NewMessage(int64(id), update.Message.CommandArguments())
				bot.Send(msg)
			}
		case userAdmin && "addAdmin" == commandMessage:
			id, err := stringToUint(update.Message.CommandArguments())
			if err != nil {
				log.Fatal(err)
			}

			err = addAdmin(database.GetDB(), id)
			if err != nil {
				textMessageUser = "Ошибка добавления пользователя в группу Администрирования!"
			} else {
				textMessageUser = "Пользователь добавлен в группу!"
			}
		case userAdmin && "deleteAdmin" == commandMessage:
			id, err := stringToUint(update.Message.CommandArguments())
			if err != nil {
				log.Fatal(err)
			}
			deleteAdminUser(database.GetDB(), id)
			thisAdmin := isAnAdmin(database.GetDB(), int64(id))
			if thisAdmin {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Операция выполнена!")
				bot.Send(msg)
			}

		case "help" == commandMessage:
			if userAdmin {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, adminMenu())
				bot.Send(msg)
			}
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, userMenu())
			bot.Send(msg)

		default:
			textMessageUser = "Команда не известна!!! Попробуйте задать другую команду!!!"
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, textMessageUser)
			bot.Send(msg)
		}

	}
}

//Returns a list of user ID that interacted with the bot
func getUserTelegramID(db *gorm.DB) []uint {
	users := []model.Users{}
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
	user := model.Users{}
	db.Where("Telegram_id = ?", id).First(&user)

	if user.Id == 0 {
		//user = Users{UserName: nameUser, TelegramId: uint(id)}
		user.UserName, user.TelegramId = nameUser, uint(id)
		result := db.Create(&user)
		if result.Error != nil {
			log.Fatal(result.Error)
		}
	}
}

//Checking the user for administrator rights
func isAnAdmin(db *gorm.DB, id int64) bool {
	var isAdmin bool
	adminUser := model.AdministratorsGroup{}
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
	/ipUserHistory [Телеграм id] (в разработке)
	Показывает все ip из за просов пользователя`
	return
}

func userMenu() (menuUser string) {
	menuUser = `Доступные команды Пользователя:
	/ip [URL для проверки]
	Отправка url на проверку`
	return
}

//Adding a user s to the admin group
func addAdmin(db *gorm.DB, id uint) (err error) {
	adminUser := model.AdministratorsGroup{}
	db.Where("id = ?", id).First(&adminUser)

	if adminUser.Id == 0 {
		adminUser = model.AdministratorsGroup{Id: uint(id)}
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

//Removes a user from the administrators group
func deleteAdminUser(db *gorm.DB, id uint) {
	admin := model.AdministratorsGroup{}
	db.Where("Id = ?", id).Delete(&admin)
}

func stringToUint(stringId string) (id uint, err error) {
	if n, err := strconv.Atoi(stringId); err == nil {
		id = uint(n)
	} else {
		log.Fatal(err)
	}
	return
}
