package serviceMsgHandle

import (
	"game-message-core/proto"

	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	"github.com/Meland-Inc/game-services/src/services/main/playerModel"
)

func SavePlayerDataHandler(iMsg interface{}) {
	input, ok := iMsg.(*proto.SavePlayerEvent)
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
		input.CurHp,
		sceneData.Level,
		sceneData.Exp,
		input.MapId,
		input.Position.X,
		input.Position.Y,
		input.Position.Z,
		input.Dir.X,
		input.Dir.Y,
		input.Dir.Z,
	)
}
