package serviceMsgHandle

import (
	"game-message-core/proto"

	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	"github.com/Meland-Inc/game-services/src/services/main/playerModel"
)

func PlayerDeathHandler(iMsg interface{}) {
	input, ok := iMsg.(*proto.PlayerDeathEvent)
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
		input.UserId, input.Position, input.KillerId,
		proto.EntityType(input.KillerType), input.KillerName,
	); err != nil {
		serviceLog.Error("PlayerDeathEventData OnPlayerDeath err: %v", err)
		return
	}
}
