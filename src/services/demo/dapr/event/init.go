package daprEvent

import (
	"context"

	"github.com/Meland-Inc/game-services/src/common/daprInvoke"
	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	"github.com/dapr/go-sdk/service/common"
)

func InitDaprPubsubEvent() (err error) {
	daprInvoke.AddTopicEventHandler("DemoServiceTestEventHandler", DemoServiceTestEventHandler)
	if err != nil {
		return err
	}

	return nil
}

func DemoServiceTestEventHandler(ctx context.Context, e *common.TopicEvent) (retry bool, err error) {
	serviceLog.Info("Receive DemoServiceTestEvent nft: %v, :%s ", e.Data, e.DataContentType)

	// input := &pubsubEventData.PlayerDeathEventData{}
	// err = grpcNetTool.UnmarshalGrpcTopicEvent(e, input)
	// if err != nil {
	// 	serviceLog.Error("PlayerDeathEvent UnmarshalEvent fail err: %v ", err)
	// 	return false, err
	// }
	return false, nil
}
