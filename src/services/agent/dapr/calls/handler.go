package daprCalls

import (
	"context"
	"fmt"
	"game-message-core/proto"
	"game-message-core/protoTool"
	"net/url"

	"github.com/Meland-Inc/game-services/src/common/daprInvoke"
	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	"github.com/Meland-Inc/game-services/src/services/agent/userChannel"
	"github.com/dapr/go-sdk/service/common"
)

func BroadCastToClientHandler(ctx context.Context, in *common.InvocationEvent) (*common.Content, error) {
	input := &proto.BroadCastToClientInput{}
	err := protoTool.UnmarshalProto(in.Data, input)
	if err != nil {
		escStr, err := url.QueryUnescape(string(in.Data))
		serviceLog.Info("agent received BroadCastToClient QueryUnescape data: %v, err: %+v", escStr, err)
		if err != nil {
			return nil, err
		}
		err = protoTool.UnmarshalProto([]byte(escStr), input)
		if err != nil {
			serviceLog.Error("Unmarshal to BroadCastToClientInput data : %+v, err: $+v", string(in.Data), err)
			return nil, fmt.Errorf("data can not unMarshal to BroadCastToClientInput")
		}
	}

	var userCh *userChannel.UserChannel
	if input.SocketId != "" {
		userCh = userChannel.GetInstance().UserChannelById(input.SocketId)
	} else if input.UserId > 0 {
		userCh = userChannel.GetInstance().UserChannelByOwner(input.UserId)
	}
	if userCh == nil {
		serviceLog.Error("BroadCastToClient userCha not found  userId[%d], socketId[%v]", input.UserId, input.SocketId)
		return nil, fmt.Errorf(" user channel is not found")
	}

	userCh.SendToUser(input.Msg.Type, input.Msg)
	output := &proto.BroadCastToClientOutput{Success: true}
	serviceLog.Info("register service res = %+v", output)
	return daprInvoke.MakeOutputContent(in, output)
}

func MultipleBroadCastToClientHandler(ctx context.Context, in *common.InvocationEvent) (*common.Content, error) {
	input := &proto.MultipleBroadCastToClientInput{}
	err := protoTool.UnmarshalProto(in.Data, input)
	if err != nil {
		escStr, err := url.QueryUnescape(string(in.Data))
		serviceLog.Info("agent received MultipleBroadCastToClient QueryUnescape data: %v, err: %+v", escStr, err)
		if err != nil {
			return nil, err
		}
		err = protoTool.UnmarshalProto([]byte(escStr), input)
		if err != nil {
			serviceLog.Error("Unmarshal to MultipleBroadCastToClient data : %+v, err: $+v", string(in.Data), err)
			return nil, fmt.Errorf("data can not unMarshal to BroadCastToClientInput")
		}
	}

	for _, userId := range input.UserList {
		userCh := userChannel.GetInstance().UserChannelByOwner(userId)
		if userCh != nil {
			userCh.SendToUser(proto.EnvelopeType(input.MsgId), input.Msg)
		} else {
			serviceLog.Warning("UserChannel [%d] not found", userId)
		}
	}

	output := &proto.MultipleBroadCastToClientOutput{Success: true}
	serviceLog.Info("register service res = %+v", output)
	return daprInvoke.MakeOutputContent(in, output)
}
