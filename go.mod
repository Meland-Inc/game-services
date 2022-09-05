module github.com/Meland-Inc/game-services

go 1.17

require (
	github.com/go-sql-driver/mysql v1.6.0 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	gorm.io/driver/mysql v1.3.6 // indirect
	gorm.io/gorm v1.23.8 // indirect
)

// 使用本地go代码仓库方式: https://zhuanlan.zhihu.com/p/109828249
require game-message-core v0.0.0

replace game-message-core => ./src/game-message-core/messageGo
