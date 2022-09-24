package daprEvent

import (
	"context"
	"encoding/json"
	"fmt"
	"game-message-core/grpc"
	"game-message-core/proto"
	"game-message-core/protoTool"
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

	input := &proto.SavePlayerEvent{}
	err = protoTool.UnmarshalProto(bs, input)
	if err != nil {
		escStr, err := url.QueryUnescape(string(bs))
		serviceLog.Info("received SavePlayerEvent QueryUnescape data: %v, err: %+v", escStr, err)
		if err != nil {
			return false, err
		}
		err = protoTool.UnmarshalProto([]byte(escStr), input)
		if err != nil {
			serviceLog.Error("Unmarshal to SavePlayerEvent data : %+v, err: $+v", string(bs), err)
			return false, fmt.Errorf("data can not unMarshal to SavePlayerEvent")
		}
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
	bs, err := json.Marshal(e.Data)
	if err != nil {
		serviceLog.Error("KillMonsterEvent  marshal e.Data  fail err: %+v", err)
		return false, fmt.Errorf("KillMonsterEvent  marshal e.Data  fail err: %+v", err)
	}

	input := &proto.KillMonsterEvent{}
	err = protoTool.UnmarshalProto(bs, input)
	if err != nil {
		escStr, err := url.QueryUnescape(string(bs))
		serviceLog.Info("received KillMonsterEvent QueryUnescape data: %v, err: %+v", escStr, err)
		if err != nil {
			return false, err
		}
		err = protoTool.UnmarshalProto([]byte(escStr), input)
		if err != nil {
			serviceLog.Error("Unmarshal to KillMonsterEvent data : %+v, err: $+v", string(bs), err)
			return false, fmt.Errorf("data can not unMarshal to KillMonsterEvent")
		}
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
		serviceLog.Error("PlayerDeathEvent  marshal e.Data  fail err: %+v", err)
		return false, fmt.Errorf("PlayerDeathEvent  marshal e.Data  fail err: %+v", err)
	}

	input := &proto.PlayerDeathEvent{}
	err = protoTool.UnmarshalProto(bs, input)
	if err != nil {
		escStr, err := url.QueryUnescape(string(bs))
		serviceLog.Info("received PlayerDeathEvent QueryUnescape data: %v, err: %+v", escStr, err)
		if err != nil {
			return false, err
		}
		err = protoTool.UnmarshalProto([]byte(escStr), input)
		if err != nil {
			serviceLog.Error("Unmarshal to PlayerDeathEvent data : %+v, err: $+v", string(bs), err)
			return false, fmt.Errorf("data can not unMarshal to PlayerDeathEvent")
		}
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
