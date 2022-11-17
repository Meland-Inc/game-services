package daprCalls

import (
	"context"
	"fmt"
	base_data "game-message-core/grpc/baseData"
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
	serviceData, _ := controller.GetInstance().GetAliveServiceByType(input.ServiceType, input.MapId)
	if serviceData == nil {
		output.ErrorCode = 30001
		output.ErrorMessage = fmt.Sprintf("Service [%v][%d]not found", input.ServiceType, input.MapId)
	} else {
		output.ServiceType = serviceData.ServiceType
		output.ServiceAppId = serviceData.AppId
		output.MapId = serviceData.MapId
		output.Host = serviceData.Host
		output.Port = serviceData.Port
		output.Online = serviceData.Online
		output.MaxOnline = serviceData.MaxOnline
		output.CreatedAt = serviceData.CreateAt
		output.UpdateAt = serviceData.UpdateAt
	}
	serviceLog.Info("select service res = %+v", output)
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
		ser := base_data.ServiceData{
			AppId:       s.AppId,
			ServiceType: s.ServiceType,
			Host:        s.Host,
			Port:        s.Port,
			MapId:       s.MapId,
			Online:      s.Online,
			MaxOnline:   s.MaxOnline,
			CreatedAt:   s.CreateAt,
			UpdatedAt:   s.UpdateAt,
		}
		output.Services = append(output.Services, ser)
	}
	return daprInvoke.MakeOutputContent(in, output)
}
