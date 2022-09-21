package daprEvent

import (
	"context"
	"encoding/json"
	"fmt"
	"game-message-core/grpc"
	"game-message-core/grpc/pubsubEventData"
	"net/url"

	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	"github.com/Meland-Inc/game-services/src/services/main/msgChannel"
	"github.com/dapr/go-sdk/service/common"
	"github.com/spf13/cast"
)

func SavePlayerDataEventHandle(ctx context.Context, e *common.TopicEvent) (retry bool, err error) {
	bs, err := json.Marshal(e.Data)
	if err != nil {
		serviceLog.Error("save player data  marshal e.Data  fail err: %+v", err)
		return false, fmt.Errorf("save player data  marshal e.Data  fail err: %+v", err)
	}

	escStr, err := url.QueryUnescape(string(bs))
	serviceLog.Info("Receive save player data: %v, err: %v", escStr, err)

	input := &pubsubEventData.SavePlayerEventData{}
	err = json.Unmarshal([]byte(escStr), input)
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
	bs, err := json.Marshal(e.Data)
	if err != nil {
		serviceLog.Error("KillMonsterEventData  marshal e.Data  fail err: %+v", err)
		return false, fmt.Errorf("KillMonsterEventData  marshal e.Data  fail err: %+v", err)
	}

	escStr, err := url.QueryUnescape(string(bs))
	serviceLog.Info("Receive KillMonsterEventData data: %v, err: %v", escStr, err)

	input := &pubsubEventData.KillMonsterEventData{}
	err = json.Unmarshal([]byte(escStr), input)
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
	bs, err := json.Marshal(e.Data)
	if err != nil {
		serviceLog.Error("PlayerDeathEventData  marshal e.Data  fail err: %+v", err)
		return false, fmt.Errorf("PlayerDeathEventData  marshal e.Data  fail err: %+v", err)
	}

	escStr, err := url.QueryUnescape(string(bs))
	serviceLog.Info("Receive PlayerDeathEventData data: %v, err: %v", escStr, err)

	input := &pubsubEventData.PlayerDeathEventData{}
	err = json.Unmarshal([]byte(escStr), input)
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
