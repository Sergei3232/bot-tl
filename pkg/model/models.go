package model

type Users struct {
	Id         uint `gorm:"primaryKey"`
	UserName   string
	TelegramId uint
}

type AdministratorsGroup struct {
	Id uint `gorm:"primaryKey"`
}
