package serviceMsgHandle

import (
	"game-message-core/grpc/pubsubEventData"

	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	"github.com/Meland-Inc/game-services/src/global/userAgent"
)

func UserEnterGameHandle(iMsg interface{}) {
	input, ok := iMsg.(*pubsubEventData.UserEnterGameEvent)
	if !ok {
		serviceLog.Error("iMsg to UserEnterGameEvent failed")
		return
	}

	agentModel := userAgent.GetUserAgentModel()
	agent, exist := agentModel.GetUserAgent(input.UserId)
	if !exist {
		agentModel.AddUserAgentRecord(input.UserId, input.AgentAppId, input.UserSocketId)
	} else {
		agent.TryUpdate(input.UserId, input.AgentAppId, input.UserSocketId)
	}
}
