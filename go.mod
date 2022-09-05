module github.com/Meland-Inc/game-services

go 1.17

require (
	github.com/dapr/dapr v1.8.0 // indirect
	github.com/dapr/go-sdk v1.5.0 // indirect
	github.com/go-sql-driver/mysql v1.6.0 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	golang.org/x/net v0.0.0-20220621193019-9d032be2e588 // indirect
	golang.org/x/sys v0.0.0-20220520151302-bc2c85ada10a // indirect
	golang.org/x/text v0.3.7 // indirect
	google.golang.org/genproto v0.0.0-20220622171453-ea41d75dfa0f // indirect
	google.golang.org/grpc v1.47.0 // indirect
	google.golang.org/protobuf v1.28.1 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	gorm.io/driver/mysql v1.3.6 // indirect
	gorm.io/gorm v1.23.8 // indirect
)

// 使用本地go代码仓库方式: https://zhuanlan.zhihu.com/p/109828249
require game-message-core v0.0.0

replace game-message-core => ./src/game-message-core/messageGo
