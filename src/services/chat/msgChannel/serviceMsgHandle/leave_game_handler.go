package serviceMsgHandle

import (
	"game-message-core/grpc/pubsubEventData"

	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	"github.com/Meland-Inc/game-services/src/global/userAgent"
	"github.com/Meland-Inc/game-services/src/services/chat/chatModel"
)

func PlayerLeaveGameHandler(iMsg interface{}) {
	input, ok := iMsg.(*pubsubEventData.UserLeaveGameEvent)
	if !ok {
		serviceLog.Error("iMsg to UserLeaveGameEvent failed")
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
