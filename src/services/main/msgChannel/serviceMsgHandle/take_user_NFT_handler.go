package serviceMsgHandle

import (
	"game-message-core/grpc/methodData"

	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	"github.com/Meland-Inc/game-services/src/services/main/playerModel"
)

func TakeUserNftHandler(iMsg interface{}) {
	input, ok := iMsg.(*methodData.MainServiceActionTakeNftInput)
	if !ok {
		serviceLog.Error("iMsg to MainServiceActionTakeNftInput failed")
		return
	}

	dataModel, err := playerModel.GetPlayerDataModel()
	if err != nil {
		serviceLog.Error("Main Service Take Nft playerDataModel not found")
		return
	}

	for _, takeNft := range input.TakeNfts {
		if takeNft.NftId != "" {
			err = dataModel.TakeNftById(input.UserId, takeNft.NftId, takeNft.Num)
		} else {
			err = dataModel.TakeNftByItemCid(input.UserId, takeNft.ItemCid, takeNft.Num)
		}
		if err != nil {
			serviceLog.Error(err.Error())
		}
	}
}
