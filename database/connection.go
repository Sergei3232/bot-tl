package database

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"log"
	"time"
)

var dbase *gorm.DB

type Users struct {
	Id         uint
	UserName   string
	TelegramId uint
}

type AdministratorsGroup struct {
	IdUsers uint
}

func Init() *gorm.DB {
	db, err := gorm.Open("postgres", "user=leralarina dbname=bot_tg_db sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	db.AutoMigrate(&Users{}, &AdministratorsGroup{})

	return db
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
