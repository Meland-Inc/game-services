package daprCalls

import (
	"context"
	"encoding/json"
	"fmt"
	"game-message-core/grpc/methodData"
	"net/url"

	"github.com/Meland-Inc/game-services/src/common/daprInvoke"
	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	"github.com/Meland-Inc/game-services/src/services/manager/controller"
	"github.com/dapr/go-sdk/service/common"
)

func toLocalServiceData(input methodData.ServiceDataInput) controller.ServiceData {
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
		CreatedAt:   input.CreatedAt,
		UpdatedAt:   input.UpdatedAt,
	}
}

func RegisterServiceHandler(ctx context.Context, in *common.InvocationEvent) (*common.Content, error) {
	escStr, err := url.QueryUnescape(string(in.Data))
	serviceLog.Info("received register service  data: %v, err: %v", escStr, err)
	if err != nil {
		return nil, err
	}

	input := methodData.ServiceDataInput{}
	err = json.Unmarshal([]byte(escStr), &input)
	if err != nil {
		serviceLog.Error("register service  data : %+v, err:%+v", string(escStr), err)
		return nil, fmt.Errorf("data can not unMarshal to ServiceDataInput")
	}

	service := toLocalServiceData(input)
	controller.GetInstance().RegisterService(service)

	output := methodData.ServiceDataOutput{
		MsgVersion: input.MsgVersion,
		Success:    true,
	}
	// serviceLog.Info("register service res = %+v", output)

	return daprInvoke.MakeOutputContent(in, output)
}

func DestroyServiceHandler(ctx context.Context, in *common.InvocationEvent) (*common.Content, error) {
	escStr, err := url.QueryUnescape(string(in.Data))
	serviceLog.Info("received Destroy service  data: %v, err: %v", escStr, err)
	if err != nil {
		return nil, err
	}

	input := methodData.ServiceDataInput{}
	err = json.Unmarshal([]byte(escStr), &input)
	if err != nil {
		serviceLog.Error("Destroy service  data : %+v", string(in.Data))
		return nil, fmt.Errorf("data can not unMarshal to ServiceDataInput")
	}

	service := toLocalServiceData(input)
	controller.GetInstance().DestroyService(service)

	output := methodData.ServiceDataOutput{
		MsgVersion: input.MsgVersion,
		Success:    true,
	}

	serviceLog.Info("Destroy service res = %+v", output)
	return daprInvoke.MakeOutputContent(in, output)
}
