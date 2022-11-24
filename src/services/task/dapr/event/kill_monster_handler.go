package daprEvent

import (
	"context"
	"fmt"
	"game-message-core/grpc"
	"game-message-core/grpc/pubsubEventData"

	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	"github.com/Meland-Inc/game-services/src/global/grpcAPI/grpcNetTool"
	"github.com/Meland-Inc/game-services/src/global/serviceCnf"
	"github.com/Meland-Inc/game-services/src/services/task/msgChannel"
	"github.com/dapr/go-sdk/service/common"
	"github.com/spf13/cast"
)

func PlayerKillMonsterEventHandle(ctx context.Context, e *common.TopicEvent) (retry bool, err error) {
	serviceLog.Info("task Receive KillMonsterEvent: %v ", e.Data)

	input := &pubsubEventData.KillMonsterEventData{}
	err = grpcNetTool.UnmarshalGrpcTopicEvent(e, input)
	if err != nil {
		serviceLog.Error("task KillMonsterEvent UnmarshalEvent fail err: %v ", err)
		return false, err
	}

	// 抛弃过期事件
	if input.MsgVersion < serviceCnf.GetInstance().StartMs {
		return false, nil
	}

	userId := cast.ToInt64(input.UserId)
	if userId < 1 {
		serviceLog.Error("task KillMonsterEventData invalid Data[%v]", input)
		return false, fmt.Errorf("task KillMonsterEventData invalid Data [%v]", input)
	}

	msgChannel.GetInstance().CallServiceMsg(&msgChannel.ServiceMsgData{
		MsgId:   string(grpc.SubscriptionEventKillMonster),
		MsgBody: input,
	})

	return false, nil
}
