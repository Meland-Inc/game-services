package serviceMsgHandle

import (
	"game-message-core/grpc/pubsubEventData"

	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	"github.com/Meland-Inc/game-services/src/global/userAgent"
	login_model "github.com/Meland-Inc/game-services/src/services/main/loginModel"
)

func PlayerLeaveGameHandler(iMsg interface{}) {
	input, ok := iMsg.(*pubsubEventData.UserLeaveGameEvent)
	if !ok {
		serviceLog.Error("iMsg to UserLeaveGameInput failed")
		return
	}

	agentModel := userAgent.GetUserAgentModel()
	agentModel.RemoveUserAgentRecord(input.UserId)

	loginModel, _ := login_model.GetLoginModel()
	loginModel.OnLogOut(input.UserId)
}
