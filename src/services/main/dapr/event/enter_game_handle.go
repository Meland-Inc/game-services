package daprEvent

import (
	"context"
	"encoding/json"
	"game-message-core/grpc"
	"game-message-core/grpc/pubsubEventData"

	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	"github.com/Meland-Inc/game-services/src/services/main/msgChannel"
	"github.com/dapr/go-sdk/service/common"
)

func UserEnterGameEventHandler(ctx context.Context, e *common.TopicEvent) (retry bool, err error) {
	serviceLog.Info("received enter game: %v, %s", e.Data, e.DataContentType)

	inputBytes, err := json.Marshal(e.Data)
	if err != nil {
		serviceLog.Info("enter game Marshal(e.Data) fail err: %v ", err)
		return false, err
	}

	input := &pubsubEventData.UserEnterGameEvent{}
	err = json.Unmarshal(inputBytes, &input)
	if err != nil {
		serviceLog.Info("enter game Marshal to enterGameInput fail err: %v ", err)
		return false, err
	}

	serviceLog.Info("receive enterGameData: %+v ", input)

	msgChannel.GetInstance().CallServiceMsg(&msgChannel.ServiceMsgData{
		MsgId:   string(grpc.SubscriptionEventUserEnterGame),
		MsgBody: input,
	})

	return false, nil
}
