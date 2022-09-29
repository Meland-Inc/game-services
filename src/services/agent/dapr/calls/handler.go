package daprCalls

import (
	"context"
	"fmt"
	"game-message-core/grpc/methodData"
	"game-message-core/proto"
	"game-message-core/protoTool"

	"github.com/Meland-Inc/game-services/src/common/daprInvoke"
	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	"github.com/Meland-Inc/game-services/src/global/grpcAPI/grpcNetTool"
	"github.com/Meland-Inc/game-services/src/services/agent/userChannel"
	"github.com/dapr/go-sdk/service/common"
)

func BroadCastToClientHandler(ctx context.Context, in *common.InvocationEvent) (*common.Content, error) {
	// serviceLog.Info("agent received BroadCastToClient data: %v", (in.Data))
	input := &methodData.BroadCastToClientInput{}
	err := grpcNetTool.UnmarshalGrpcData(in.Data, input)
	if err != nil {
		return nil, err
	}

	resMsg, err := protoTool.UnMarshalToEnvelope(input.MsgBody)
	serviceLog.Info("BroadCastToClient msg[%+v], err:%+v", resMsg.Type, err)
	var userCh *userChannel.UserChannel
	if input.SocketId != "" {
		userCh = userChannel.GetInstance().UserChannelById(input.SocketId)
	} else if input.UserId > 0 {
		userCh = userChannel.GetInstance().UserChannelByOwner(input.UserId)
	}
	if userCh == nil {
		serviceLog.Error("BroadCastToClient userCh not found  userId[%d], socketId[%v]", input.UserId, input.SocketId)
		return nil, fmt.Errorf(" user channel is not found")
	}

	userCh.SendToUser(proto.EnvelopeType(input.MsgId), input.MsgBody)
	output := &methodData.BroadCastToClientOutput{Success: true}
	return daprInvoke.MakeOutputContent(in, output)
}

func MultipleBroadCastToClientHandler(ctx context.Context, in *common.InvocationEvent) (*common.Content, error) {
	// serviceLog.Info("agent received MultipleBroadCastToClient data: %v", string(in.Data))

	input := &methodData.MultipleBroadCastToClientInput{}
	err := grpcNetTool.UnmarshalGrpcData(in.Data, input)
	if err != nil {
		return nil, err
	}

	resMsg, err := protoTool.UnMarshalToEnvelope(input.MsgBody)
	serviceLog.Info("MultipleBroadCastToClient msg[%+v], err:%+v", resMsg.Type, err)
	for _, userId := range input.UserList {
		userCh := userChannel.GetInstance().UserChannelByOwner(userId)
		if userCh != nil {
			userCh.SendToUser(proto.EnvelopeType(input.MsgId), input.MsgBody)
		} else {
			serviceLog.Warning("UserChannel [%d] not found", userId)
		}
	}

	output := &methodData.BroadCastToClientOutput{Success: true}
	return daprInvoke.MakeOutputContent(in, output)
}
