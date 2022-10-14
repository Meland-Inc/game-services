package daprCalls

import (
	"context"
	"game-message-core/grpc/methodData"

	"github.com/Meland-Inc/game-services/src/common/daprInvoke"
	"github.com/Meland-Inc/game-services/src/global/grpcAPI/grpcNetTool"
	"github.com/Meland-Inc/game-services/src/services/account/msgChannel"
	"github.com/dapr/go-sdk/service/common"
)

func ClientMessageHandler(ctx context.Context, in *common.InvocationEvent) (*common.Content, error) {
	resFunc := func(success bool, err error) (*common.Content, error) {
		output := &methodData.PullClientMessageOutput{}
		output.Success = success
		if err != nil {
			output.ErrMsg = err.Error()
		}
		content, _ := daprInvoke.MakeOutputContent(in, output)
		return content, err
	}

	input := &methodData.PullClientMessageInput{}
	err := grpcNetTool.UnmarshalGrpcData(in.Data, input)
	if err != nil {
		return nil, err
	}

	msgChannel.GetInstance().CallClientMsg(input)
	return resFunc(true, nil)
}
