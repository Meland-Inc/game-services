package serviceMsgHandle

import (
	"game-message-core/grpc/pubsubEventData"

	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	"github.com/Meland-Inc/game-services/src/global/userAgent"
)

func UserChangeServiceHandler(iMsg interface{}) {
	input, ok := iMsg.(*pubsubEventData.UserChangeServiceEvent)
	if !ok {
		serviceLog.Error("iMsg to UserEnterGameEvent failed")
		return
	}

	agentModel := userAgent.GetUserAgentModel()
	agent, exist := agentModel.GetUserAgent(input.UserId)
	if exist {
		agent.TryUpdate(agent.UserId, agent.AgentAppId, agent.SocketId, input.ToService.AppId)
	} else {
		agent, _ = agentModel.AddUserAgentRecord(
			input.UserId,
			input.UserAgentAppId,
			input.UserSocketId,
			input.ToService.AppId,
		)
	}
}
