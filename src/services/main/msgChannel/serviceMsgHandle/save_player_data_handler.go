package serviceMsgHandle

import (
	"game-message-core/grpc/pubsubEventData"

	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	"github.com/Meland-Inc/game-services/src/services/main/playerModel"
)

func SavePlayerDataHandler(iMsg interface{}) {
	input, ok := iMsg.(*pubsubEventData.SavePlayerEventData)
	if !ok {
		serviceLog.Error("iMsg to SavePlayerEventData failed")
		return
	}

	dataModel, err := playerModel.GetPlayerDataModel()
	if err != nil {
		serviceLog.Error("SavePlayerEventData playerDataModel not found")
		return
	}

	sceneData, err := dataModel.GetPlayerSceneData(input.UserId)
	if err != nil {
		serviceLog.Error("SavePlayerEventData scene Data  not found")
		return
	}

	dataModel.UpPlayerSceneData(
		input.UserId,
		input.CurHP,
		sceneData.Level,
		sceneData.Exp,
		input.MapId,
		float64(input.PosX),
		float64(input.PosY),
		float64(input.PosZ),
		float64(input.DirX),
		float64(input.DirY),
		float64(input.DirZ),
	)
}
