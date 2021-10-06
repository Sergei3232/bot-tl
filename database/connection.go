package database

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"log"
	"time"
)

var dbase *gorm.DB

type Users struct {
	Id         uint `gorm:"primaryKey"`
	UserName   string
	TelegramId uint
}

type AdministratorsGroup struct {
	Id uint `gorm:"primaryKey"`
}

func Init() *gorm.DB {
	db, err := gorm.Open("postgres", "user=leralarina dbname=bot_tg_db sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	db.AutoMigrate(&Users{}, &AdministratorsGroup{})
	initAdmin(db, initOneUser(db))
	return db
}

func initAdmin(db *gorm.DB, id uint) {
	adminUser := AdministratorsGroup{Id: id}
	result := db.Create(&adminUser)
	if result.Error != nil {
		log.Fatal(result.Error)
	}
}

func initOneUser(db *gorm.DB) uint {
	user := Users{UserName: "MrS1_2", TelegramId: 519588080}
	result := db.Create(&user)
	if result.Error != nil {
		log.Fatal(result.Error)
	}
	return user.Id
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

func GetUserList(db *gorm.DB) []Users {
	var user []Users
	result := db.Find(&user)
	if result.Error != nil {
		log.Fatal(result.Error)
	}
	return user
}
