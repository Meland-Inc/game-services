package serviceMsgHandle

import (
	message "github.com/Meland-Inc/game-services/src/global/web3Message"
	"github.com/spf13/cast"

	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	"github.com/Meland-Inc/game-services/src/services/main/playerModel"
)

func Web3UpdateUserNftHandler(iMsg interface{}) {
	input, ok := iMsg.(*message.UpdateUserNFT)
	if !ok {
		serviceLog.Error("iMsg to UserEnterGameEvent failed")
		return
	}

	userId := cast.ToInt64(input.UserId)
	dataModel, err := playerModel.GetPlayerDataModel()
	if err != nil {
		serviceLog.Error(err.Error())
		return
	}
	dataModel.UpdatePlayerNFTs(userId, []message.NFT{input.Nft})
}

func Web3MultiUpdateUserNftHandler(iMsg interface{}) {
	input, ok := iMsg.(*message.MultiUpdateUserNFT)
	if !ok {
		serviceLog.Error("iMsg to UserEnterGameEvent failed")
		return
	}

	userId := cast.ToInt64(input.UserId)
	dataModel, err := playerModel.GetPlayerDataModel()
	if err != nil {
		serviceLog.Error(err.Error())
		return
	}
	dataModel.UpdatePlayerNFTs(userId, input.Nfts)
}
