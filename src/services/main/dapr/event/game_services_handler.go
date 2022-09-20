package daprEvent

import (
	"context"
	"encoding/json"
	"fmt"
	"game-message-core/grpc"
	"game-message-core/grpc/pubsubEventData"

	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	"github.com/Meland-Inc/game-services/src/services/main/msgChannel"
	"github.com/dapr/go-sdk/service/common"
	"github.com/spf13/cast"
)

func SavePlayerDataEventHandle(ctx context.Context, e *common.TopicEvent) (retry bool, err error) {
	serviceLog.Info("Receive save player data: %v %s \n", e.Data, e.DataContentType)

	bs, err := json.Marshal(e.Data)
	input := &pubsubEventData.SavePlayerEventData{}
	err = json.Unmarshal(bs, input)
	if err != nil {
		serviceLog.Error("not math to dapr msg SavePlayerEventData data : %+v", e.Data)
		return false, fmt.Errorf("not math to dapr msg SavePlayerEventData")
	}

	userId := cast.ToInt64(input.UserId)
	if userId < 1 {
		serviceLog.Error("SavePlayerEventData invalid Data[%v]", input)
		return false, fmt.Errorf("SavePlayerEventData invalid Data [%v]", input)
	}

	msgChannel.GetInstance().CallServiceMsg(&msgChannel.ServiceMsgData{
		MsgId:   string(grpc.SubscriptionEventSavePlayerData),
		MsgBody: input,
	})

	return false, nil
}

func PlayerKillMonsterEventHandle(ctx context.Context, e *common.TopicEvent) (retry bool, err error) {
	serviceLog.Info("Receive  player kill monster data: %v %s \n", e.Data, e.DataContentType)

	bs, err := json.Marshal(e.Data)
	input := &pubsubEventData.KillMonsterEventData{}
	err = json.Unmarshal(bs, input)
	if err != nil {
		serviceLog.Error("not math to dapr msg KillMonsterEventData data : %+v", e.Data)
		return false, fmt.Errorf("not math to dapr msg KillMonsterEventData")
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
	serviceLog.Info("Receive player death data: %v %s \n", e.Data, e.DataContentType)

	bs, err := json.Marshal(e.Data)
	input := &pubsubEventData.PlayerDeathEventData{}
	err = json.Unmarshal(bs, input)
	if err != nil {
		serviceLog.Error("not math to dapr msg PlayerDeathEventData data : %+v", e.Data)
		return false, fmt.Errorf("not math to dapr msg PlayerDeathEventData")
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
