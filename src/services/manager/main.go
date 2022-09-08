package main

import (
	"github.com/Meland-Inc/game-services/src/application"
	mgrService "github.com/Meland-Inc/game-services/src/services/manager/service"
)

func main() {
	mgrSer := mgrService.NewManagerService()
	application.Init(mgrSer)
	application.Run()
}
