package main

import (
	"github.com/Meland-Inc/game-services/src/application"
	"github.com/Meland-Inc/game-services/src/services/main/service"
)

func main() {
	mainSer := service.NewMainService()
	application.Init(mainSer)
	application.Run()
}
