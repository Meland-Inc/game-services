package serviceMsgHandle

import (
	"game-message-core/grpc/methodData"

	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	"github.com/Meland-Inc/game-services/src/global/userAgent"
)

func PlayerLeaveGameHandler(iMsg interface{}) {
	input, ok := iMsg.(*methodData.UserLeaveGameInput)
	if !ok {
		serviceLog.Error("iMsg to UserLeaveGameInput failed")
		return
	}

	agentModel := userAgent.GetUserAgentModel()
	agentModel.RemoveUserAgentRecord(input.UserId)
}
