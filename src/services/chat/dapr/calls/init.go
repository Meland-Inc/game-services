package daprCalls

import (
	"context"
	"game-message-core/grpc/methodData"

	"github.com/Meland-Inc/game-services/src/common/daprInvoke"
	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	"github.com/Meland-Inc/game-services/src/global/grpcAPI/grpcNetTool"
	"github.com/dapr/go-sdk/service/common"
)

func InitDaprCallHandle() (err error) {
	daprInvoke.AddServiceInvocationHandler("DemoServiceTestCallsHandler", DemoServiceTestCallsHandler)
	if err != nil {
		return err
	}

	return nil
}

func DemoServiceTestCallsHandler(ctx context.Context, in *common.InvocationEvent) (*common.Content, error) {
	serviceLog.Info("this is DemoServiceTestCallsHandler, %s", string(in.Data))

	input := &methodData.BroadCastToClientInput{}
	err := grpcNetTool.UnmarshalGrpcData(in.Data, input)
	if err != nil {
		return nil, err
	}

	out := &common.Content{
		Data:        []byte{},
		ContentType: in.ContentType,
		DataTypeURL: in.DataTypeURL,
	}
	return out, nil
}
