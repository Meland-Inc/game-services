package daprCalls

import (
	"context"
	"game-message-core/grpc/methodData"

	"github.com/Meland-Inc/game-services/src/common/daprInvoke"
	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	"github.com/Meland-Inc/game-services/src/common/time_helper"
	"github.com/Meland-Inc/game-services/src/global/grpcAPI/grpcNetTool"
	"github.com/Meland-Inc/game-services/src/services/manager/controller"
	"github.com/dapr/go-sdk/service/common"
)

func RegisterServiceHandler(ctx context.Context, in *common.InvocationEvent) (*common.Content, error) {
	serviceLog.Info("service register: %s", string(in.Data))

	input := &methodData.ServiceRegisterInput{}
	err := grpcNetTool.UnmarshalGrpcData(in.Data, input)
	if err != nil {
		return nil, err
	}

	service := controller.ServiceData{
		AppId:           input.Service.AppId,
		ServiceType:     input.Service.ServiceType,
		SceneSerSubType: input.Service.SceneSerSubType,
		HomeOwner:       input.Service.Owner,
		Host:            input.Service.Host,
		Port:            input.Service.Port,
		MapId:           input.Service.MapId,
		Online:          input.Service.Online,
		MaxOnline:       input.Service.MaxOnline,
		CreateAt:        input.Service.CreatedAt,
		UpdateAt:        input.Service.UpdatedAt,
	}

	controller.GetInstance().RegisterService(service)
	output := &methodData.ServiceRegisterOutput{
		Success:    true,
		RegisterAt: input.RegisterAt,
		ManagerAt:  time_helper.NowUTCMill(),
	}
	return daprInvoke.MakeOutputContent(in, output)
}
