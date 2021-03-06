package database

import (
	"fmt"
	"github.com/Sergei3232/bot-tl/pkg/model"
	"github.com/jinzhu/gorm"
	"log"
	"time"
)

var dbase *gorm.DB

func Init() *gorm.DB {
	db, err := gorm.Open("postgres", "user=leralarina dbname=bot_tg_db sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	db.AutoMigrate(&model.Users{}, &model.AdministratorsGroup{})
	initAdmin(db, initOneUser(db))
	return db
}

func initAdmin(db *gorm.DB, id uint) {
	adminUser := model.AdministratorsGroup{}
	db.Where("id = ?", id).First(&adminUser)

	if adminUser.Id == 0 {
		//adminUser = model.AdministratorsGroup{Id: id}
		adminUser.Id = id
		result := db.Create(&adminUser)
		if result.Error != nil {
			log.Fatal(result.Error)
		}
	}

}

func initOneUser(db *gorm.DB) uint {
	user := model.Users{}
	db.Where("Telegram_id = ?", 519588080).First(&user)

	if user.Id == 0 {
		//user = model.Users{UserName: "MrS1_2", TelegramId: 519588080}
		user.UserName, user.TelegramId = "MrS1_2", 519588080
		result := db.Create(&user)
		if result.Error != nil {
			log.Fatal(result.Error)
		}
	}
	return user.TelegramId
}

func GetDB() *gorm.DB {
	if dbase == nil {
		dbase = Init()
		var sleep = time.Duration(1)
		for dbase == nil {
			sleep = sleep * 2
			fmt.Printf("database is unavailable. Wait for %d sec.\n", sleep)
			time.Sleep(sleep * time.Second)
			dbase = Init()
		}
	}

	return dbase
}
