package main

import (
	"github.com/Meland-Inc/game-services/src/application"
	"github.com/Meland-Inc/game-services/src/services/demo/service"
)

func main() {
	demoSer := service.NewDemoService()
	application.Init(demoSer)
	application.Run()
}
