package daprCalls

import (
	"context"
	"game-message-core/grpc"
	"game-message-core/grpc/methodData"

	"github.com/Meland-Inc/game-services/src/common/daprInvoke"
	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	"github.com/Meland-Inc/game-services/src/global/grpcAPI/grpcNetTool"
	"github.com/Meland-Inc/game-services/src/services/task/msgChannel"
	"github.com/dapr/go-sdk/service/common"
)

func UserLeaveGameHandler(ctx context.Context, in *common.InvocationEvent) (*common.Content, error) {
	resFunc := func(success bool, err error) (*common.Content, error) {
		out := &methodData.UserLeaveGameOutput{}
		out.Success = success
		if err != nil {
			out.ErrMsg = err.Error()
			serviceLog.Error("get user data err: %v", err)
		}
		content, _ := daprInvoke.MakeOutputContent(in, out)
		return content, err
	}

	serviceLog.Info("task service user Leave data: %v", string(in.Data))

	input := &methodData.UserLeaveGameInput{}
	err := grpcNetTool.UnmarshalGrpcData(in.Data, input)
	if err != nil {
		return nil, err
	}

	msgChannel.GetInstance().CallServiceMsg(&msgChannel.ServiceMsgData{
		MsgId:   string(grpc.UserActionLeaveGame),
		MsgBody: input,
	})
	return resFunc(true, nil)
}
