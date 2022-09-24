package serviceMsgHandle

import (
	"game-message-core/proto"

	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	"github.com/Meland-Inc/game-services/src/global/grpcAPI/grpcInvoke"
	"github.com/Meland-Inc/game-services/src/services/main/playerModel"
)

func KillMonsterHandler(iMsg interface{}) {
	input, ok := iMsg.(*proto.KillMonsterEvent)
	if !ok {
		serviceLog.Error("iMsg to KillMonsterEvent failed")
		return
	}

	dataModel, err := playerModel.GetPlayerDataModel()
	if err != nil {
		serviceLog.Error("KillMonsterEvent playerDataModel not found")
		return
	}

	err = dataModel.AddExp(input.UserId, input.Exp)
	if err != nil {
		serviceLog.Error("KillMonsterEvent add exp failed: %v", err)
		return
	}

	for _, drop := range input.DropList {
		if err := grpcInvoke.MintNFT(
			input.UserId,
			drop.Cid, drop.Num, drop.Quality,
			int32(input.Position.X), int32(input.Position.Z),
		); err != nil {
			serviceLog.Error("mint nft[%d] failed: %v", drop.Cid, err)
			return
		}
	}
}
