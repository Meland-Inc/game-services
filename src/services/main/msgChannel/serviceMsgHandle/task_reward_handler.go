package serviceMsgHandle

import (
	"game-message-core/grpc/pubsubEventData"

	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	"github.com/Meland-Inc/game-services/src/services/main/playerModel"
)

func TaskRewardHandler(iMsg interface{}) {
	input, ok := iMsg.(*pubsubEventData.UserTaskRewardEvent)
	if !ok {
		serviceLog.Error("iMsg to UserTaskRewardEvent failed")
		return
	}

	dataModel, err := playerModel.GetPlayerDataModel()
	if err != nil {
		serviceLog.Error("UserTaskRewardEvent playerDataModel not found")
		return
	}

	// call mint task reward NFT is in task service, so reward exp add in there
	if err = dataModel.AddExp(input.UserId, input.Exp); err != nil {
		serviceLog.Error("UserTaskRewardEvent  addExp err: %v", err)
		return
	}

}
