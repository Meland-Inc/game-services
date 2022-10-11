package main

import (
	"fmt"

	"github.com/Meland-Inc/game-services/src/common/daprInvoke"
	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	"github.com/Meland-Inc/game-services/src/global/grpcAPI/grpcInvoke"
)

func main() {
	daprInvoke.InitClient("5550") //

	var userId int64 = 691
	fmt.Println(" mint user ", userId, " nfts begin --- ")

	makeTaskTestItem(userId, 1, 1)
	makeTestCraftItems(userId, 1, 1)

	fmt.Println(" mint user ", userId, " nfts end --- ")
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
