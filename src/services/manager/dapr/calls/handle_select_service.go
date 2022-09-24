package daprCalls

import (
	"context"
	"fmt"
	"game-message-core/proto"
	"game-message-core/protoTool"
	"net/url"

	"github.com/Meland-Inc/game-services/src/common/daprInvoke"
	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	"github.com/Meland-Inc/game-services/src/services/manager/controller"
	"github.com/dapr/go-sdk/service/common"
)

func SelectServiceHandler(ctx context.Context, in *common.InvocationEvent) (*common.Content, error) {

	// escStr, err := url.QueryUnescape(string(in.Data))
	// if err != nil {
	// 	return nil, err
	// }

	// input := methodData.ManagerActionSelectServiceInput{}
	// err = json.Unmarshal([]byte(escStr), &input)
	// if err != nil {
	// 	serviceLog.Error("select service  data : %+v", string(in.Data))
	// 	return nil, fmt.Errorf("data can not unMarshal to select service input")
	// }

	input := &proto.ManagerActionSelectServiceInput{}
	err := protoTool.UnmarshalProto(in.Data, input)
	if err != nil {
		escStr, err := url.QueryUnescape(string(in.Data))
		serviceLog.Warning("received select service  data: %v, err: %v", input, err)
		if err != nil {
			return nil, err
		}
		err = protoTool.UnmarshalProto([]byte(escStr), input)
		if err != nil {
			serviceLog.Error("received select service  data: %v, err: %v", input, err)
			return nil, fmt.Errorf("data can not unMarshal to BroadCastToClientInput")
		}
	}

	serviceLog.Info("received select service  data: %v, err: %v", input, err)

	output := &proto.ManagerActionSelectServiceOutput{}
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
		output.CreateAt = serviceData.CreateAt
		output.UpdateAt = serviceData.UpdateAt
	}
	serviceLog.Info("select service res = %+v", output)
	return daprInvoke.MakeOutputContent(in, output)
}
