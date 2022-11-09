package daprEvent

import (
	"context"
	"fmt"
	"game-message-core/grpc"
	"game-message-core/grpc/pubsubEventData"

	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	"github.com/Meland-Inc/game-services/src/global/grpcAPI/grpcNetTool"
	"github.com/Meland-Inc/game-services/src/services/main/msgChannel"
	"github.com/dapr/go-sdk/service/common"
	"github.com/spf13/cast"
)

func TaskRewardEventHandler(ctx context.Context, e *common.TopicEvent) (retry bool, err error) {
	input := &pubsubEventData.UserTaskRewardEvent{}
	err = grpcNetTool.UnmarshalGrpcTopicEvent(e, input)
	if err != nil {
		serviceLog.Error("UserTaskRewardEvent Unmarshal fail err: %v ", err)
		return false, nil
	}

	serviceLog.Info("Receive task reward Event: %+v", input)

	userId := cast.ToInt64(input.UserId)
	if userId < 1 {
		serviceLog.Error("SavePlayerEvent invalid Data[%v]", input)
		return false, fmt.Errorf("SavePlayerEvent invalid Data [%v]", input)
	}

	msgChannel.GetInstance().CallServiceMsg(&msgChannel.ServiceMsgData{
		MsgId:   string(grpc.SubscriptionEventUserTaskReward),
		MsgBody: input,
	})

	return false, nil
}
