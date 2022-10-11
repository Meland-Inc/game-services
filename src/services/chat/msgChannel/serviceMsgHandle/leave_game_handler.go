package serviceMsgHandle

import (
	"game-message-core/grpc/methodData"

	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	"github.com/Meland-Inc/game-services/src/global/serviceCnf"
	"github.com/Meland-Inc/game-services/src/global/userAgent"
	"github.com/Meland-Inc/game-services/src/services/chat/chatModel"
)

func PlayerLeaveGameHandler(iMsg interface{}) {
	input, ok := iMsg.(*methodData.UserLeaveGameInput)
	if !ok {
		serviceLog.Error("iMsg to UserLeaveGameInput failed")
		return
	}

	if input.MsgVersion < serviceCnf.GetInstance().StartMs {
		return
	}

	agentModel := userAgent.GetUserAgentModel()
	agentModel.RemoveUserAgentRecord(input.UserId)

	model, _ := chatModel.GetChatModel()
	if model == nil {
		serviceLog.Error("chat model not found")
		return
	}
	if model != nil {
		model.OnPlayerLeaveGame(input.UserId)
	}
}
