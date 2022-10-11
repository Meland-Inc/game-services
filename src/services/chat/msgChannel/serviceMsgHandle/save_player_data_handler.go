package serviceMsgHandle

import (
	"game-message-core/grpc/pubsubEventData"

	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	"github.com/Meland-Inc/game-services/src/services/chat/chatModel"
)

func SavePlayerDataHandler(iMsg interface{}) {
	input, ok := iMsg.(*pubsubEventData.SavePlayerEventData)
	if !ok {
		serviceLog.Error("iMsg to SavePlayerEvent failed")
		return
	}

	model, _ := chatModel.GetChatModel()
	if model == nil {
		serviceLog.Error("chat model not found")
		return
	}

	model.OnUpdatePlayerData(input)
}
