package main

import (
	"github.com/Meland-Inc/game-services/src/application"
	"github.com/Meland-Inc/game-services/src/services/chat/service"
)

func main() {
	chatSer := service.NewChatService()
	application.Init(chatSer)
	application.Run()
}
