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
	"github.com/spf13/cast"
)

func TaskFinishHandler(ctx context.Context, e *common.TopicEvent) (retry bool, err error) {
	serviceLog.Info("task Receive TaskFinishEvent: %v", e.Data)

	input := &pubsubEventData.TaskFinishEvent{}
	err = grpcNetTool.UnmarshalGrpcTopicEvent(e, input)
	if err != nil {
		serviceLog.Error("task TaskFinishEvent UnmarshalEvent fail err: %v ", err)
		return false, nil
	}

	// 抛弃过期事件
	if input.MsgVersion < serviceCnf.GetInstance().StartMs {
		return false, nil
	}

	userId := cast.ToInt64(input.UserId)
	if userId < 1 {
		serviceLog.Error("task TaskFinishEvent invalid Data[%v]", input)
		return false, nil
	}

	msgChannel.GetInstance().CallServiceMsg(&msgChannel.ServiceMsgData{
		MsgId:   string(grpc.SubscriptionEventTaskFinish),
		MsgBody: input,
	})

	return false, nil
}