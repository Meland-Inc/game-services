package serviceMsgHandle

import (
	"game-message-core/grpc/pubsubEventData"

	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	"github.com/Meland-Inc/game-services/src/services/main/playerModel"
)

func SavePlayerDataHandler(iMsg interface{}) {
	input, ok := iMsg.(*pubsubEventData.SavePlayerEventData)
	if !ok {
		serviceLog.Error("iMsg to SavePlayerEvent failed")
		return
	}

	dataModel, err := playerModel.GetPlayerDataModel()
	if err != nil {
		serviceLog.Error("SavePlayerEvent playerDataModel not found")
		return
	}

	sceneData, err := dataModel.GetPlayerSceneData(input.UserId)
	if err != nil {
		serviceLog.Error("SavePlayerEvent scene Data  not found")
		return
	}

	dataModel.UpPlayerSceneData(
		input.UserId,
		input.CurHP,
		sceneData.Level,
		sceneData.Exp,
		input.MapId,
		input.PosX,
		input.PosY,
		input.PosZ,
		input.DirX,
		input.DirY,
		input.DirZ,
	)
}
