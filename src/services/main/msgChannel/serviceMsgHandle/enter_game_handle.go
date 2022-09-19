package serviceMsgHandle

import (
	"game-message-core/grpc/pubsubEventData"

	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	"github.com/Meland-Inc/game-services/src/global/component"
	"github.com/Meland-Inc/game-services/src/global/userAgent"
)

func UserEnterGameHandle(iMsg interface{}) {
	env, ok := iMsg.(*pubsubEventData.UserEnterGameEvent)
	if !ok {
		serviceLog.Error("iMsg to UserEnterGameEvent failed")
		return
	}

	iUserAgentModel, exist := component.GetInstance().GetModel(component.MODEL_NAME_USER_AGENT)
	if !exist {
		return
	}
	agentModel := iUserAgentModel.(*userAgent.UserAgentModel)
	agent, exist := agentModel.GetUserAgent(env.BaseData.UserId)
	if exist {
		agent.InSceneServiceAppId = env.SceneServiceAppId
		agent.InMapId = env.MapId
	}
}
