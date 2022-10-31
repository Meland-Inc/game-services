package daprEvent

import (
	"context"

	"game-message-core/grpc/pubsubEventData"

	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	"github.com/Meland-Inc/game-services/src/global/grpcAPI/grpcNetTool"
	"github.com/Meland-Inc/game-services/src/services/manager/controller"
	"github.com/dapr/go-sdk/service/common"
)

func ServiceUnRegisterHandler(ctx context.Context, e *common.TopicEvent) (retry bool, err error) {
	serviceLog.Info("Receive service UnRegister data: %v", e.Data)

	input := &pubsubEventData.ServiceUnregisterEvent{}
	err = grpcNetTool.UnmarshalGrpcTopicEvent(e, input)
	if err != nil {
		serviceLog.Error("ServiceUnregisterEvent UnmarshalEvent fail err: %v ", err)
		return false, err
	}

	service := controller.ServiceData{
		AppId:       input.Service.AppId,
		ServiceType: input.Service.ServiceType,
		Host:        input.Service.Host,
		Port:        input.Service.Port,
		MapId:       input.Service.MapId,
		Online:      input.Service.Online,
		MaxOnline:   input.Service.MaxOnline,
		CreateAt:    input.Service.CreatedAt,
		UpdateAt:    input.Service.UpdatedAt,
	}
	controller.GetInstance().DestroyService(service)
	return false, nil
}
