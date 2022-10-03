package serviceMsgHandle

import (
	"game-message-core/grpc/pubsubEventData"

	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	"github.com/Meland-Inc/game-services/src/global/userAgent"
	"github.com/Meland-Inc/game-services/src/services/chat/chatModel"
)

func UserEnterGameHandle(iMsg interface{}) {
	input, ok := iMsg.(*pubsubEventData.UserEnterGameEvent)
	if !ok {
		serviceLog.Error("iMsg to UserEnterGameEvent failed")
		return
	}

	agentModel := userAgent.GetUserAgentModel()
	agent, exist := agentModel.GetUserAgent(input.UserId)
	if exist {
		agent.InSceneServiceAppId = input.SceneServiceAppId
		agent.SocketId = input.UserSocketId
		agent.AgentAppId = input.AgentAppId
		agent.InMapId = input.MapId
	} else {
		agent, _ = agentModel.AddUserAgentRecord(input.UserId, input.AgentAppId, input.UserSocketId)
		agent.InSceneServiceAppId = input.SceneServiceAppId
	}

	model, _ := chatModel.GetChatModel()
	if model != nil {
		if err := model.OnPlayerEnterGame(input); err != nil {
			serviceLog.Error(err.Error())
		}
	}
}
