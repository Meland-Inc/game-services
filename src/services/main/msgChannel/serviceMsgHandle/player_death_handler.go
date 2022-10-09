package serviceMsgHandle

import (
	"game-message-core/grpc/pubsubEventData"
	"game-message-core/proto"

	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	"github.com/Meland-Inc/game-services/src/services/main/playerModel"
)

func PlayerDeathHandler(iMsg interface{}) {
	input, ok := iMsg.(*pubsubEventData.PlayerDeathEventData)
	if !ok {
		serviceLog.Error("iMsg to PlayerDeathEventData failed")
		return
	}

	dataModel, err := playerModel.GetPlayerDataModel()
	if err != nil {
		serviceLog.Error("PlayerDeathEventData playerDataModel not found")
		return
	}

	if err = dataModel.OnPlayerDeath(
		input.UserId, &proto.Vector3{X: input.PosX, Y: input.PosY, Z: input.PosZ},
		input.KillerId, proto.EntityType(input.KillerType), input.KillerName,
	); err != nil {
		serviceLog.Error("PlayerDeathEventData OnPlayerDeath err: %v", err)
		return
	}
}
