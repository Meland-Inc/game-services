package daprCalls

import (
	"context"
	"fmt"
	"game-message-core/grpc/methodData"
	"game-message-core/proto"

	"github.com/Meland-Inc/game-services/src/common/daprInvoke"
	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	"github.com/Meland-Inc/game-services/src/global/grpcAPI/grpcNetTool"
	"github.com/Meland-Inc/game-services/src/services/manager/controller"
	"github.com/dapr/go-sdk/service/common"
)

func SelectServiceHandler(ctx context.Context, in *common.InvocationEvent) (*common.Content, error) {
	serviceLog.Warning("received select service  data: %v", in.Data)
	input := &methodData.ManagerActionSelectServiceInput{}
	err := grpcNetTool.UnmarshalGrpcData(in.Data, input)
	if err != nil {
		return nil, err
	}

	serviceLog.Info("received select service  data: %v, err: %v", input, err)
	output := &methodData.ManagerActionSelectServiceOutput{}

	if serviceData, _ := controller.GetInstance().GetAliveServiceByType(
		input.ServiceType, input.SceneSerSubType, input.MapId, input.OwnerId,
	); serviceData == nil {
		output.ErrorCode = 30001
		output.ErrorMessage = fmt.Sprintf("Service [%v][%d]not found", input.ServiceType, input.MapId)
	} else {
		output.Service = serviceData.ToGrpcService()
	}

	serviceLog.Info("select service resErrMs[%s], serData = %+v", output.ErrorMessage, output.Service)
	return daprInvoke.MakeOutputContent(in, output)
}

func MultiSelectServiceHandler(ctx context.Context, in *common.InvocationEvent) (*common.Content, error) {
	input := &methodData.MultiSelectServiceInput{}
	err := grpcNetTool.UnmarshalGrpcData(in.Data, input)
	if err != nil {
		return nil, err
	}

	serviceLog.Info("received multi select service input: %+v", input)

	allService := controller.GetInstance().AllServices()
	output := &methodData.MultiSelectServiceOutput{}
	for _, s := range allService {
		if s.ServiceType != input.ServiceType {
			continue
		}
		if input.ServiceType == proto.ServiceType_ServiceTypeAgent &&
			input.MapId != s.MapId {
			continue
		}
		output.Services = append(output.Services, s.ToGrpcService())
	}
	return daprInvoke.MakeOutputContent(in, output)
}
