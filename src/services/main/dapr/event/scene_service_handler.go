package daprEvent

import (
	"context"
	"encoding/json"
	"fmt"
	"game-message-core/grpc"
	"game-message-core/grpc/pubsubEventData"

	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	"github.com/Meland-Inc/game-services/src/global/grpcAPI/grpcNetTool"
	"github.com/Meland-Inc/game-services/src/services/main/msgChannel"
	"github.com/dapr/go-sdk/service/common"
	"github.com/spf13/cast"
)

func SavePlayerDataEventHandle(ctx context.Context, e *common.TopicEvent) (retry bool, err error) {
	inputBytes, err := json.Marshal(e.Data)
	if err != nil {
		serviceLog.Error("save player data  marshal e.Data  fail err: %+v", err)
		return false, fmt.Errorf("save player data  marshal e.Data  fail err: %+v", err)
	}

	input := &pubsubEventData.SavePlayerEventData{}
	err = grpcNetTool.UnmarshalGrpcData(inputBytes, input)
	if err != nil {
		return false, err
	}

	userId := cast.ToInt64(input.UserId)
	if userId < 1 {
		serviceLog.Error("SavePlayerEvent invalid Data[%v]", input)
		return false, fmt.Errorf("SavePlayerEvent invalid Data [%v]", input)
	}

	msgChannel.GetInstance().CallServiceMsg(&msgChannel.ServiceMsgData{
		MsgId:   string(grpc.SubscriptionEventSavePlayerData),
		MsgBody: input,
	})

	return false, nil
}

func PlayerKillMonsterEventHandle(ctx context.Context, e *common.TopicEvent) (retry bool, err error) {
	inputBytes, err := json.Marshal(e.Data)
	if err != nil {
		serviceLog.Error("KillMonsterEvent  marshal e.Data  fail err: %+v", err)
		return false, fmt.Errorf("KillMonsterEvent  marshal e.Data  fail err: %+v", err)
	}

	input := &pubsubEventData.KillMonsterEventData{}
	err = grpcNetTool.UnmarshalGrpcData(inputBytes, input)
	if err != nil {
		return false, err
	}

	userId := cast.ToInt64(input.UserId)
	if userId < 1 {
		serviceLog.Error("KillMonsterEventData invalid Data[%v]", input)
		return false, fmt.Errorf("KillMonsterEventData invalid Data [%v]", input)
	}

	msgChannel.GetInstance().CallServiceMsg(&msgChannel.ServiceMsgData{
		MsgId:   string(grpc.SubscriptionEventKillMonster),
		MsgBody: input,
	})

	return false, nil
}

func PlayerDeathEventHandle(ctx context.Context, e *common.TopicEvent) (retry bool, err error) {
	inputBytes, err := json.Marshal(e.Data)
	if err != nil {
		serviceLog.Error("PlayerDeathEvent  marshal e.Data  fail err: %+v", err)
		return false, fmt.Errorf("PlayerDeathEvent  marshal e.Data  fail err: %+v", err)
	}

	input := &pubsubEventData.PlayerDeathEventData{}
	err = grpcNetTool.UnmarshalGrpcData(inputBytes, input)
	if err != nil {
		return false, err
	}

	userId := cast.ToInt64(input.UserId)
	if userId < 1 {
		serviceLog.Error("PlayerDeathEventData invalid Data[%v]", input)
		return false, fmt.Errorf("PlayerDeathEventData invalid Data [%v]", input)
	}

	msgChannel.GetInstance().CallServiceMsg(&msgChannel.ServiceMsgData{
		MsgId:   string(grpc.SubscriptionEventPlayerDeath),
		MsgBody: input,
	})

	return false, nil
}
