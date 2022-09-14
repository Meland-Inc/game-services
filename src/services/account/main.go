package main

import (
	"github.com/Meland-Inc/game-services/src/application"
	"github.com/Meland-Inc/game-services/src/services/account/service"
)

func main() {
	accountSer := service.NewAccountService()
	application.Init(accountSer)
	application.Run()
}
