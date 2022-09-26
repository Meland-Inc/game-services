package main

import (
	"github.com/Meland-Inc/game-services/src/application"
	"github.com/Meland-Inc/game-services/src/services/task/service"
)

func main() {
	taskSer := service.NewTaskService()
	application.Init(taskSer)
	application.Run()
}
