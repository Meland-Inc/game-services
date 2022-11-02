package daprEvent

import (
	"context"
	"game-message-core/grpc"
	"game-message-core/grpc/pubsubEventData"

	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	"github.com/Meland-Inc/game-services/src/global/grpcAPI/grpcNetTool"
	"github.com/Meland-Inc/game-services/src/global/serviceCnf"
	"github.com/Meland-Inc/game-services/src/services/main/msgChannel"
	"github.com/dapr/go-sdk/service/common"
)

func UserLeaveGameHandler(ctx context.Context, e *common.TopicEvent) (retry bool, err error) {
	input := &pubsubEventData.UserLeaveGameEvent{}
	err = grpcNetTool.UnmarshalGrpcTopicEvent(e, input)
	if err != nil {
		serviceLog.Error("UserLeaveGame UnmarshalEvent fail err: %v ", err)
		return false, nil
	}

	// 抛弃过期事件
	if input.MsgVersion < serviceCnf.GetInstance().StartMs {
		return false, nil
	}

	serviceLog.Info("main service receive enterGame: %+v", input)

	msgChannel.GetInstance().CallServiceMsg(&msgChannel.ServiceMsgData{
		MsgId:   string(grpc.SubscriptionEventUserLeaveGame),
		MsgBody: input,
	})

	return false, nil
}
