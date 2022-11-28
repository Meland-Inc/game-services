package daprEvent

import (
	"context"
	"game-message-core/grpc/pubsubEventData"
	"game-message-core/proto"
	"game-message-core/protoTool"

	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	"github.com/Meland-Inc/game-services/src/global/grpcAPI/grpcNetTool"
	"github.com/Meland-Inc/game-services/src/global/serviceCnf"
	"github.com/Meland-Inc/game-services/src/services/agent/userChannel"

	"github.com/dapr/go-sdk/service/common"
)

func TickOutUserHandler(ctx context.Context, e *common.TopicEvent) (retry bool, err error) {
	input := &pubsubEventData.TickOutPlayerEvent{}
	err = grpcNetTool.UnmarshalGrpcTopicEvent(e, input)
	if err != nil {
		serviceLog.Error("TickOutPlayer UnmarshalEvent fail err: %v ", err)
		return false, nil
	}

	// 抛弃过期事件
	if input.MsgVersion < serviceCnf.GetInstance().StartMs {
		return false, nil
	}

	serviceLog.Info("agent service receive TickOutPlayer: %+v", input)

	var userCh *userChannel.UserChannel
	if input.SocketId != "" {
		userCh = userChannel.GetInstance().UserChannelById(input.SocketId)
	} else if input.UserId > 0 {
		userCh = userChannel.GetInstance().UserChannelByOwner(input.UserId)
	}
	if userCh == nil {
		return false, nil
	}

	broadcastTickOutPlayer(userCh, input.TickOutCode)
	userCh.Stop()
	return false, nil
}

func broadcastTickOutPlayer(userCh *userChannel.UserChannel, tickCode proto.TickOutType) {
	if userCh == nil {
		return
	}
	msg := &proto.Envelope{
		Type: proto.EnvelopeType_BroadCastTickOut,
		Payload: &proto.Envelope_BroadCastTickOutResponse{
			BroadCastTickOutResponse: &proto.BroadCastTickOutResponse{
				Kind: tickCode,
			},
		},
	}
	msgBody, err := protoTool.MarshalEnvelope(msg)
	if err != nil {
		serviceLog.Error(err.Error())
		return
	}
	serviceLog.Debug("broad cast tick out user[%v][%v][%v]", userCh.GetOwner(), tickCode, userCh.GetSession().SessionId())
	userCh.SendToUser(msg.Type, msgBody)
}
