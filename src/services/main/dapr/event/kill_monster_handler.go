package daprEvent

import (
	"context"
	"fmt"
	"game-message-core/grpc"
	"game-message-core/grpc/pubsubEventData"

	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	"github.com/Meland-Inc/game-services/src/global/grpcAPI/grpcNetTool"
	"github.com/Meland-Inc/game-services/src/global/serviceCnf"
	"github.com/Meland-Inc/game-services/src/services/main/msgChannel"
	"github.com/dapr/go-sdk/service/common"
	"github.com/spf13/cast"
)

func PlayerKillMonsterEventHandle(ctx context.Context, e *common.TopicEvent) (retry bool, err error) {
	env := &pubsubEventData.KillMonsterEventData{}
	err = grpcNetTool.UnmarshalGrpcTopicEvent(e, env)
	if err != nil {
		serviceLog.Error("KillMonsterEvent Unmarshal fail err: %v ", err)
		return false, nil
	}

	// 抛弃过期事件
	if env.MsgVersion < serviceCnf.GetInstance().StartMs {
		return false, nil
	}

	serviceLog.Info("Receive KillMonsterEvent : %+v", env)

	userId := cast.ToInt64(env.UserId)
	if userId < 1 {
		serviceLog.Error("KillMonsterEventData invalid Data[%v]", env)
		return false, fmt.Errorf("KillMonsterEventData invalid Data [%v]", env)
	}

	msgChannel.GetInstance().CallServiceMsg(&msgChannel.ServiceMsgData{
		MsgId:   string(grpc.SubscriptionEventKillMonster),
		MsgBody: env,
	})

	return false, nil
}
