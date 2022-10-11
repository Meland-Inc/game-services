package serviceMsgHandle

import (
	"game-message-core/grpc/pubsubEventData"

	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	"github.com/Meland-Inc/game-services/src/global/grpcAPI/grpcInvoke"
	"github.com/Meland-Inc/game-services/src/services/main/playerModel"
)

func KillMonsterHandler(iMsg interface{}) {
	input, ok := iMsg.(*pubsubEventData.KillMonsterEventData)
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
			int32(input.PosX), int32(input.PosZ),
		); err != nil {
			serviceLog.Error("mint nft[%d] failed: %v", drop.Cid, err)
			return
		}
	}

	//     ---------- task   test data --------------------
	go func() {
		// makeTaskTestItem(input.UserId, int32(input.PosX), int32(input.PosZ))
		// makeTestCraftItems(input.UserId, int32(input.PosX), int32(input.PosZ))
	}()
}

func makeTaskTestItem(userId int64, posX, posY int32) {
	for itemId := 3010101; itemId <= 3010102; itemId++ {
		err := grpcInvoke.MintNFT(userId, int32(itemId), 10, 1, posX, posY)
		if err != nil {
			serviceLog.Error("mint nft[%d] failed: %v", itemId, err)
		}
	}

	for itemId := 3010201; itemId <= 3010202; itemId++ {
		err := grpcInvoke.MintNFT(userId, int32(itemId), 10, 1, posX, posY)
		if err != nil {
			serviceLog.Error("mint nft[%d] failed: %v", itemId, err)
		}
	}

	for itemId := 1010001; itemId <= 1010012; itemId++ {
		err := grpcInvoke.MintNFT(userId, int32(itemId), 1, 1, posX, posY)
		if err != nil {
			serviceLog.Error("mint nft[%d] failed: %v", itemId, err)
		}
	}

}

func makeTestCraftItems(userId int64, posX, posY int32) {
	for itemId := 3010101; itemId <= 3010102; itemId++ {
		err := grpcInvoke.MintNFT(userId, int32(itemId), 10, 1, posX, posY)
		if err != nil {
			serviceLog.Error("mint nft[%d] failed: %v", itemId, err)
		}
	}

	for itemId := 3010201; itemId <= 3010202; itemId++ {
		err := grpcInvoke.MintNFT(userId, int32(itemId), 10, 1, posX, posY)
		if err != nil {
			serviceLog.Error("mint nft[%d] failed: %v", itemId, err)
		}
	}

	for itemId := 60000001; itemId <= 60000020; itemId++ {
		err := grpcInvoke.MintNFT(userId, int32(itemId), 10, 1, posX, posY)
		if err != nil {
			serviceLog.Error("mint nft[%d] failed: %v", itemId, err)
		}
	}

	for itemId := 4010001; itemId <= 4010012; itemId++ {
		err := grpcInvoke.MintNFT(userId, int32(itemId), 1, 1, posX, posY)
		if err != nil {
			serviceLog.Error("mint nft[%d] failed: %v", itemId, err)
		}
	}

	for itemId := 4010901; itemId <= 4010902; itemId++ {
		err := grpcInvoke.MintNFT(userId, int32(itemId), 1, 1, posX, posY)
		if err != nil {
			serviceLog.Error("mint nft[%d] failed: %v", itemId, err)
		}
	}

	for itemId := 4020001; itemId <= 4020002; itemId++ {
		err := grpcInvoke.MintNFT(userId, int32(itemId), 1, 1, posX, posY)
		if err != nil {
			serviceLog.Error("mint nft[%d] failed: %v", itemId, err)
		}
	}

}
