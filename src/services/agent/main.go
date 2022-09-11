package main

import (
	"github.com/Meland-Inc/game-services/src/application"
	"github.com/Meland-Inc/game-services/src/services/agent/service"
)

func main() {
	agentSer := service.NewAgentService()
	application.Init(agentSer)
	application.Run()
}
