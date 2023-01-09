package web3Handler

import (
	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	"github.com/Meland-Inc/game-services/src/global/contract"
	"github.com/Meland-Inc/game-services/src/global/serviceCnf"
	message "github.com/Meland-Inc/game-services/src/global/web3Message"
)

func Web3UpdateUserNftEvent(env contract.IModuleEventReq, curMs int64) {
	input := &message.UpdateUserNFT{}
	err := env.UnmarshalToDaprEventData(input)
	if err != nil {
		serviceLog.Error("Web3UpdateUserNft UnmarshalEvent fail err: %v ", err)
		return
	}

	if input.Etag < int(serviceCnf.GetInstance().StartMs/1000) {
		return
	}
	// TODO logic
}
