package daprEvent

import (
	"context"
	"game-message-core/grpc"
	"game-message-core/grpc/pubsubEventData"

	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	"github.com/Meland-Inc/game-services/src/global/grpcAPI/grpcNetTool"
	"github.com/Meland-Inc/game-services/src/services/main/msgChannel"
	"github.com/dapr/go-sdk/service/common"
)

func UserEnterGameEventHandler(ctx context.Context, e *common.TopicEvent) (retry bool, err error) {
	serviceLog.Info("received enter game: %v, %s", e.Data, e.DataContentType)

	input := &pubsubEventData.UserEnterGameEvent{}
	err = grpcNetTool.UnmarshalGrpcTopicEvent(e, input)
	if err != nil {
		serviceLog.Error("UserEnterGame UnmarshalEvent fail err: %v ", err)
		return false, err
	}

	serviceLog.Info("receive enterGameData: %+v ", input)

	msgChannel.GetInstance().CallServiceMsg(&msgChannel.ServiceMsgData{
		MsgId:   string(grpc.SubscriptionEventUserEnterGame),
		MsgBody: input,
	})

	return false, nil
}
