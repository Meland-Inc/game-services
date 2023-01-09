package serviceHandler

import (
	"game-message-core/grpc/pubsubEventData"

	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	"github.com/Meland-Inc/game-services/src/global/contract"
	"github.com/Meland-Inc/game-services/src/global/serviceCnf"
)

func GRPCSaveHomeDataEvent(env contract.IModuleEventReq, curMs int64) {
	input := &pubsubEventData.SaveHomeEvent{}
	err := env.UnmarshalToDaprEventData(input)
	if err != nil {
		serviceLog.Error("SaveHomeDataEvent UnmarshalEvent fail err: %v ", err)
		return
	}

	if input.MsgVersion < serviceCnf.GetInstance().StartMs {
		return
	}

	// TODO: logic
}
