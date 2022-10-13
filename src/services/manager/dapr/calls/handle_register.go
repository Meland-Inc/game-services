package daprCalls

import (
	"context"
	"game-message-core/grpc/methodData"

	"github.com/Meland-Inc/game-services/src/common/daprInvoke"
	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	"github.com/Meland-Inc/game-services/src/global/grpcAPI/grpcNetTool"
	"github.com/Meland-Inc/game-services/src/services/manager/controller"
	"github.com/dapr/go-sdk/service/common"
)

func toLocalServiceData(input *methodData.ServiceDataInput) controller.ServiceData {
	return controller.ServiceData{
		Id:          input.Id,
		Name:        input.Name,
		AppId:       input.AppId,
		ServiceType: input.ServiceType,
		Host:        input.Host,
		Port:        input.Port,
		MapId:       input.MapId,
		Online:      input.Online,
		MaxOnline:   input.MaxOnline,
		CreateAt:    input.CreatedAt,
		UpdateAt:    input.UpdatedAt,
	}
}

func RegisterServiceHandler(ctx context.Context, in *common.InvocationEvent) (*common.Content, error) {
	input := &methodData.ServiceDataInput{}
	err := grpcNetTool.UnmarshalGrpcData(in.Data, input)
	if err != nil {
		return nil, err
	}

	service := toLocalServiceData(input)
	serviceLog.Info("received register service: %v", service)
	controller.GetInstance().RegisterService(service)

	output := &methodData.ServiceDataOutput{
		Success: true,
	}
	// serviceLog.Info("register service res = %+v", output)

	return daprInvoke.MakeOutputContent(in, output)
}

func DestroyServiceHandler(ctx context.Context, in *common.InvocationEvent) (*common.Content, error) {
	serviceLog.Info("received Destroy service  data: %v", string(in.Data))
	input := &methodData.ServiceDataInput{}
	err := grpcNetTool.UnmarshalGrpcData(in.Data, input)
	if err != nil {
		return nil, err
	}

	service := toLocalServiceData(input)
	controller.GetInstance().DestroyService(service)

	output := &methodData.ServiceDataOutput{
		Success: true,
	}

	serviceLog.Info("Destroy service res = %+v", output)
	return daprInvoke.MakeOutputContent(in, output)
}
