package daprCalls

import (
	"context"
	"fmt"
	"game-message-core/proto"
	"game-message-core/protoTool"

	"github.com/Meland-Inc/game-services/src/common/daprInvoke"
	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	"github.com/Meland-Inc/game-services/src/services/main/msgChannel"
	"github.com/dapr/go-sdk/service/common"
)

func ClientMessageHandler(ctx context.Context, in *common.InvocationEvent) (*common.Content, error) {
	resFunc := func(success bool, err error) (*common.Content, error) {
		out := &proto.PullClientMessageOutput{}
		out.Success = success
		if err != nil {
			out.ErrMsg = err.Error()
		}
		content, _ := daprInvoke.MakeProtoOutputContent(in, out)
		return content, err
	}

	input := &proto.PullClientMessageInput{}
	err := protoTool.UnmarshalProto(in.Data, input)
	serviceLog.Info("main service received clientPbMsg err:%+v, input:%+v, data: %v", err, input, string(in.Data))
	if err != nil {
		serviceLog.Error("main service Unmarshal to PullClientMessageInput fail err: %+v", err)
		return resFunc(false, fmt.Errorf("PullClientMessageInput unMarshal fail"))
	}

	msgChannel.GetInstance().CallClientMsg(input)
	return resFunc(true, nil)
}
