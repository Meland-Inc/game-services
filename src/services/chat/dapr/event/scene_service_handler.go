package daprEvent

import (
	"context"
	"fmt"
	"game-message-core/grpc"
	"game-message-core/grpc/pubsubEventData"

	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	"github.com/Meland-Inc/game-services/src/global/grpcAPI/grpcNetTool"
	"github.com/Meland-Inc/game-services/src/services/chat/msgChannel"
	"github.com/dapr/go-sdk/service/common"
	"github.com/spf13/cast"
)

func SavePlayerDataEventHandler(ctx context.Context, e *common.TopicEvent) (retry bool, err error) {
	serviceLog.Info("chat Receive SavePlayerDataEvent nft: %v", e.Data)

	input := &pubsubEventData.SavePlayerEventData{}
	err = grpcNetTool.UnmarshalGrpcTopicEvent(e, input)
	if err != nil {
		serviceLog.Error("chat SavePlayerDataEvent UnmarshalEvent fail err: %v ", err)
		return false, err
	}

	userId := cast.ToInt64(input.UserId)
	if userId < 1 {
		serviceLog.Error("chat SavePlayerEvent invalid Data[%v]", input)
		return false, fmt.Errorf("SavePlayerEvent invalid Data [%v]", input)
	}

	msgChannel.GetInstance().CallServiceMsg(&msgChannel.ServiceMsgData{
		MsgId:   string(grpc.SubscriptionEventSavePlayerData),
		MsgBody: input,
	})

	return false, nil
}
