package daprCalls

import (
	"context"
	"game-message-core/grpc/methodData"
	"game-message-core/proto"
	"game-message-core/protoTool"

	"github.com/Meland-Inc/game-services/src/common/daprInvoke"
	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	"github.com/Meland-Inc/game-services/src/global/grpcAPI/grpcNetTool"
	"github.com/Meland-Inc/game-services/src/global/grpcAPI/grpcPubsubEvent"
	"github.com/Meland-Inc/game-services/src/services/agent/userChannel"
	"github.com/dapr/go-sdk/service/common"
)

func ignoreMsgLog(msgType proto.EnvelopeType) bool {
	switch msgType {
	case proto.EnvelopeType_BroadCastItemAdd,
		proto.EnvelopeType_BroadCastEntityMove,
		proto.EnvelopeType_BroadCastMapEntityUpdate:
		return true
	}
	return false
}

func getUserChannel(userId int64, socketId string) *userChannel.UserChannel {
	var userCh *userChannel.UserChannel
	if socketId != "" {
		userCh = userChannel.GetInstance().UserChannelById(socketId)
	} else if userId > 0 {
		userCh = userChannel.GetInstance().UserChannelByOwner(userId)
	}
	return userCh
}

func BroadCastToClientHandler(ctx context.Context, in *common.InvocationEvent) (*common.Content, error) {
	// serviceLog.Info("agent received BroadCastToClient data: %v", string(in.Data))
	input := &methodData.BroadCastToClientInput{}
	err := grpcNetTool.UnmarshalGrpcData(in.Data, input)
	if err != nil {
		return nil, err
	}

	resMsg, err := protoTool.UnMarshalToEnvelope(input.MsgBody)
	if !ignoreMsgLog(resMsg.Type) {
		serviceLog.Info("BroadCastToClient user[%d] msg[%+v], err:%+v", input.UserId, resMsg.Type, err)
	}

	userCh := getUserChannel(input.UserId, input.SocketId)
	if userCh == nil {
		serviceLog.Error("BroadCastToClient userCh not found  userId[%d], socketId[%v]", input.UserId, input.SocketId)
		if input.UserId > 0 {
			grpcPubsubEvent.RPCPubsubEventLeaveGame(input.UserId)
		}
		return nil, nil
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
	if !ignoreMsgLog(resMsg.Type) {
		serviceLog.Info("MultipleBroadCastToClient Users:%v, msg[%+v], err:%+v", input.UserList, resMsg.Type, err)
	}
	for _, userId := range input.UserList {
		userCh := getUserChannel(userId, "")
		if userCh != nil {
			userCh.SendToUser(proto.EnvelopeType(input.MsgId), input.MsgBody)
		} else {
			serviceLog.Warning("UserChannel [%d] not found", userId)
			grpcPubsubEvent.RPCPubsubEventLeaveGame(userId)
		}
	}

	output := &methodData.BroadCastToClientOutput{Success: true}
	return daprInvoke.MakeOutputContent(in, output)
}
