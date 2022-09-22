package daprCalls

import (
	"context"
	"encoding/json"
	"fmt"
	"game-message-core/grpc/methodData"
	"game-message-core/proto"
	"net/url"

	"github.com/Meland-Inc/game-services/src/common/daprInvoke"
	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	"github.com/Meland-Inc/game-services/src/services/agent/userChannel"
	"github.com/dapr/go-sdk/service/common"
)

func BroadCastToClientHandler(ctx context.Context, in *common.InvocationEvent) (*common.Content, error) {
	escStr, err := url.QueryUnescape(string(in.Data))
	if err != nil {
		return nil, err
	}

	serviceLog.Info("agent received BroadCastToClient data: %v", escStr)

	input := methodData.BroadCastToClientInput{}
	err = json.Unmarshal([]byte(escStr), &input)
	if err != nil {
		serviceLog.Error("Unmarshal to BroadCastToClientInput data : %+v, err: $+v", string(in.Data), err)
		return nil, fmt.Errorf("data can not unMarshal to BroadCastToClientInput")
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

	userCh.SendToUser(proto.EnvelopeType(input.MsgId), input.MsgBody)
	output := &methodData.BroadCastToClientOutput{Success: true}
	serviceLog.Info("register service res = %+v", output)
	return daprInvoke.MakeOutputContent(in, output)
}

func MultipleBroadCastToClientHandler(ctx context.Context, in *common.InvocationEvent) (*common.Content, error) {
	serviceLog.Info("agent received MultipleBroadCastToClient data: %v", string(in.Data))
	escStr, err := url.QueryUnescape(string(in.Data))
	if err != nil {
		return nil, err
	}

	serviceLog.Info("agent received MultipleBroadCastToClient data: %v", escStr)

	input := methodData.MultipleBroadCastToClientInput{}
	err = json.Unmarshal([]byte(escStr), &input)
	if err != nil {
		serviceLog.Error("Unmarshal to MultipleBroadCastToClient data : %+v, err: $+v", string(in.Data), err)
		return nil, fmt.Errorf("data can not unMarshal to MultipleBroadCastToClient")
	}

	for _, userId := range input.UserList {
		userCh := userChannel.GetInstance().UserChannelByOwner(userId)
		if userCh != nil {
			userCh.SendToUser(proto.EnvelopeType(input.MsgId), input.MsgBody)
		} else {
			serviceLog.Warning("UserChannel [%d] not found", userId)
		}
	}

	output := &methodData.BroadCastToClientOutput{Success: true}
	serviceLog.Info("register service res = %+v", output)
	return daprInvoke.MakeOutputContent(in, output)
}
