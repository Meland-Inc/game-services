package daprCalls

import (
	"context"

	"fmt"
	"game-message-core/proto"
	"game-message-core/protoTool"

	"github.com/Meland-Inc/game-services/src/common/daprInvoke"
	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	"github.com/Meland-Inc/game-services/src/services/account/msgChannel"
	"github.com/dapr/go-sdk/service/common"
)

func ClientMessageHandler(ctx context.Context, in *common.InvocationEvent) (*common.Content, error) {
	resFunc := func(success bool, err error) (*common.Content, error) {
		output := &proto.PullClientMessageOutput{}
		output.Success = success
		if err != nil {
			output.ErrMsg = err.Error()
		}
		content, _ := daprInvoke.MakeProtoOutputContent(in, output)
		return content, err
	}

	input := &proto.PullClientMessageInput{}
	err := protoTool.UnmarshalProto(in.Data, input)
	serviceLog.Info("account received clientPbMsg input:%+v, err: %v", input, err)
	if err != nil {
		serviceLog.Error("account Unmarshal to PullClientMessageInput fail err: %+v", err)
		return resFunc(false, fmt.Errorf("PullClientMessageInput unMarshal fail"))
	}

	msgChannel.GetInstance().CallClientMsg(input)
	return resFunc(true, nil)
}
