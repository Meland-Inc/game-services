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

func UserChangeServiceHandler(ctx context.Context, e *common.TopicEvent) (retry bool, err error) {
	input := &pubsubEventData.UserChangeServiceEvent{}
	err = grpcNetTool.UnmarshalGrpcTopicEvent(e, input)
	if err != nil {
		serviceLog.Error("UserChangeServiceEvent UnmarshalEvent fail err: %v ", err)
		return false, err
	}

	// 抛弃过期事件
	if input.MsgVersion < serviceCnf.GetInstance().StartMs {
		serviceLog.Warning("task service receive timeout user change Service event data: %+v", input)
		return false, nil
	}

	serviceLog.Info("task service receive UserChangeServiceEvent:  %+v", input)

	msgChannel.GetInstance().CallServiceMsg(&msgChannel.ServiceMsgData{
		MsgId:   string(grpc.SubscriptionEventUserChangeService),
		MsgBody: input,
	})

	return false, nil
}
