module github.com/Sergei3232/bot-tl

go 1.17

require (
	github.com/go-telegram-bot-api/telegram-bot-api v4.6.4+incompatible
	github.com/jinzhu/gorm v1.9.16
	github.com/lib/pq v1.10.3
)

require (
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/technoweenie/multipartstreamer v1.0.1 // indirect
)

replace example.com/database => ../database

replace example.com/pkg/model => ../models
