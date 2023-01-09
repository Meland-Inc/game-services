package serviceHandler

import (
	"game-message-core/grpc/pubsubEventData"

	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	"github.com/Meland-Inc/game-services/src/global/contract"
	"github.com/Meland-Inc/game-services/src/global/serviceCnf"
	"github.com/Meland-Inc/game-services/src/services/manager/controller"
)

func GRPCServiceUnregisterEvent(env contract.IModuleEventReq, curMs int64) {
	input := &pubsubEventData.ServiceUnregisterEvent{}
	err := env.UnmarshalToDaprEventData(input)
	if err != nil {
		serviceLog.Error("ServiceUnregisterEvent UnmarshalEvent fail err: %v ", err)
		return
	}

	if input.MsgVersion < serviceCnf.GetInstance().StartMs {
		return
	}

	serviceLog.Info("service UnRegister: %v", input)

	ctlModel, _ := controller.GetControllerModel()
	service := controller.ToServiceData(input.Service)
	ctlModel.DestroyService(service)
}
