package daprEvent

import (
	"context"
	"game-message-core/grpc"
	"game-message-core/grpc/pubsubEventData"

	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	"github.com/Meland-Inc/game-services/src/global/grpcAPI/grpcNetTool"
	"github.com/Meland-Inc/game-services/src/global/serviceCnf"
	"github.com/Meland-Inc/game-services/src/services/task/msgChannel"
	"github.com/dapr/go-sdk/service/common"
)

func UserEnterGameEventHandler(ctx context.Context, e *common.TopicEvent) (retry bool, err error) {
	input := &pubsubEventData.UserEnterGameEvent{}
	err = grpcNetTool.UnmarshalGrpcTopicEvent(e, input)
	if err != nil {
		serviceLog.Error("UserEnterGame UnmarshalEvent fail err: %v ", err)
		return false, err
	}

	// 抛弃过期事件
	if input.MsgVersion < serviceCnf.GetInstance().StartMs {
		serviceLog.Warning("task service receive timeout user enterGame event data: %+v", input)
		return false, nil
	}

	serviceLog.Info("task service receive enterGameEvent: [%+v], [%+v], [%+v],[%+v] ",
		input.UserId, input.UserSocketId, input.AgentAppId, input.SceneServiceAppId)

	msgChannel.GetInstance().CallServiceMsg(&msgChannel.ServiceMsgData{
		MsgId:   string(grpc.SubscriptionEventUserEnterGame),
		MsgBody: input,
	})

	return false, nil
}
