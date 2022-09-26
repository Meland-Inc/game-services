package daprCalls

import (
	"context"
	"fmt"
	"game-message-core/grpc/methodData"

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
