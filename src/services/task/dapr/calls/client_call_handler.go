package daprCalls

import (
	"context"
	"game-message-core/grpc/methodData"
	"game-message-core/proto"

	"github.com/Meland-Inc/game-services/src/common/daprInvoke"
	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	"github.com/Meland-Inc/game-services/src/global/grpcAPI/grpcNetTool"
	"github.com/Meland-Inc/game-services/src/services/task/msgChannel"
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

	input := &methodData.PullClientMessageInput{}
	err := grpcNetTool.UnmarshalGrpcData(in.Data, input)
	if err != nil {
		serviceLog.Error("task service Unmarshal clientPbMsg failed: %v", err)
		return nil, err
	}

	serviceLog.Info("task service received clientMsg [%v]", proto.EnvelopeType(input.MsgId))
	msgChannel.GetInstance().CallClientMsg(input)
	return resFunc(true, nil)
}
