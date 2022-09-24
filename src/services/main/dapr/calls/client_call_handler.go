package daprCalls

import (
	"context"
	"game-message-core/grpc/methodData"

	"github.com/Meland-Inc/game-services/src/common/daprInvoke"
	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	"github.com/Meland-Inc/game-services/src/global/grpcAPI/grpcNetTool"
	"github.com/Meland-Inc/game-services/src/services/main/msgChannel"
	"github.com/dapr/go-sdk/service/common"
)

func ClientMessageHandler(ctx context.Context, in *common.InvocationEvent) (*common.Content, error) {
	resFunc := func(success bool, err error) (*common.Content, error) {
		out := &methodData.PullClientMessageOutput{}
		out.Success = success
		if err != nil {
			out.ErrMsg = err.Error()
		}
		content, _ := daprInvoke.MakeOutputContent(in, out)
		return content, err
	}

	serviceLog.Info("main service received clientPbMsg data: %v", string(in.Data))

	input := &methodData.PullClientMessageInput{}
	err := grpcNetTool.UnmarshalGrpcData(in.Data, input)
	if err != nil {
		return nil, err
	}

	msgChannel.GetInstance().CallClientMsg(input)
	return resFunc(true, nil)
}
