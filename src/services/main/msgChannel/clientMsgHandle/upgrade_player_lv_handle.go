package clientMsgHandle

import (
	"game-message-core/grpc/methodData"
	"game-message-core/proto"

	"github.com/Meland-Inc/game-services/src/services/main/playerModel"
)

func UpgradePlayerLevelHandle(input *methodData.PullClientMessageInput, msg *proto.Envelope) {
	res := &proto.UpgradePlayerLevelResponse{}
	respMsg := makeResponseMsg(msg)
	defer func() {
		if respMsg.ErrorMessage != "" {
			respMsg.ErrorCode = 20006 // TODO: USE PROTO ERROR CODE
		}
		respMsg.Payload = &proto.Envelope_UpgradePlayerLevelResponse{UpgradePlayerLevelResponse: res}
		ResponseClientMessage(input, respMsg)
	}()

	if input.UserId < 1 {
		respMsg.ErrorMessage = "Invalid User ID"
		return
	}

	dataModel, err := playerModel.GetPlayerDataModel()
	if err != nil {
		respMsg.ErrorMessage = err.Error()
		return
	}

	lv, exp, err := dataModel.UpgradePlayerLevel(input.UserId)
	if err != nil {
		respMsg.ErrorMessage = err.Error()
		return
	}
	res.CurExp = int64(exp)
	res.CurLevel = lv
}
