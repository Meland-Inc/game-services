package serviceMsgHandle

import (
	"game-message-core/grpc/pubsubEventData"

	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	"github.com/Meland-Inc/game-services/src/global/userAgent"
)

func UserEnterGameHandle(iMsg interface{}) {
	env, ok := iMsg.(*pubsubEventData.UserEnterGameEvent)
	if !ok {
		serviceLog.Error("iMsg to UserEnterGameEvent failed")
		return
	}

	agentModel := userAgent.GetUserAgentModel()
	agent, exist := agentModel.GetUserAgent(env.UserId)
	if exist {
		agent.InSceneServiceAppId = env.SceneServiceAppId
		agent.InMapId = env.MapId
	} else {
		serviceLog.Error("UserEnterGameEvent user agent not found")
	}
}
