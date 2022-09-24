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

func toLocalServiceData(input *proto.ServiceRegisterInput) controller.ServiceData {
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
		CreateAt:    input.CreateAt,
		UpdateAt:    input.UpdateAt,
	}
}

func RegisterServiceHandler(ctx context.Context, in *common.InvocationEvent) (*common.Content, error) {
	input := &proto.ServiceRegisterInput{}
	err := protoTool.UnmarshalProto(in.Data, input)
	if err != nil {
		serviceLog.Warning("manager received ServiceRegisterInput  data: %v, err: %v", input, err)
		escStr, err := url.QueryUnescape(string(in.Data))
		serviceLog.Info("manager received ServiceRegisterInput data: %v, err: %+v", escStr, err)
		if err != nil {
			return nil, err
		}
		err = protoTool.UnmarshalProto([]byte(escStr), input)
		if err != nil {
			serviceLog.Error("Unmarshal to ServiceRegisterInput data : %+v, err: $+v", string(in.Data), err)
			return nil, fmt.Errorf("data can not unMarshal to ServiceRegisterInput")
		}
	}
	serviceLog.Info("received ServiceRegisterInput data: %v, err: %v", input, err)

	service := toLocalServiceData(input)
	controller.GetInstance().RegisterService(service)

	output := &proto.ServiceRegisterOutput{
		Success: true,
	}
	// serviceLog.Info("register service res = %+v", output)

	return daprInvoke.MakeProtoOutputContent(in, output)
}

func DestroyServiceHandler(ctx context.Context, in *common.InvocationEvent) (*common.Content, error) {
	input := &proto.ServiceRegisterInput{}
	err := protoTool.UnmarshalProto(in.Data, input)
	if err != nil {
		serviceLog.Warning("manager received DestroyService  data: %v, err: %v", input, err)
		escStr, err := url.QueryUnescape(string(in.Data))
		serviceLog.Info("manager received DestroyService data: %v, err: %+v", escStr, err)
		if err != nil {
			return nil, err
		}
		err = protoTool.UnmarshalProto([]byte(escStr), input)
		if err != nil {
			serviceLog.Error("Unmarshal to DestroyService data : %+v, err: $+v", string(in.Data), err)
			return nil, fmt.Errorf("data can not unMarshal to ServiceRegisterInput")
		}
	}
	serviceLog.Info("received DestroyService data: %v, err: %v", input, err)

	service := toLocalServiceData(input)
	controller.GetInstance().DestroyService(service)

	output := &proto.ServiceRegisterOutput{
		Success: true,
	}

	serviceLog.Info("Destroy service res = %+v", output)
	return daprInvoke.MakeProtoOutputContent(in, output)
}
