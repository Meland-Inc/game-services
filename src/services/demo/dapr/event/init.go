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

	// env := &pubsubEventData.PlayerDeathEventData{}
	// err = grpcNetTool.UnmarshalGrpcTopicEvent(e, env)
	// if err != nil {
	// 	serviceLog.Error("PlayerDeathEvent UnmarshalEvent fail err: %v ", err)
	// 	return false, err
	// }
	// // 抛弃过期事件
	// if env.MsgVersion < serviceCnf.GetInstance().StartMs {
	// 	return false, nil
	// }

	return false, nil
}
