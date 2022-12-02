package daprEvent

import (
	"context"
	"game-message-core/grpc"

	"github.com/Meland-Inc/game-services/src/common/daprInvoke"
	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	"github.com/Meland-Inc/game-services/src/global/component"
	"github.com/Meland-Inc/game-services/src/services/manager/controller"
	"github.com/dapr/go-sdk/service/common"
)

func InitDaprPubsubEvent() (err error) {
	serviceLog.Info(" InitDaprPubsubEvent ------ begin ------")
	if daprInvoke.AddTopicEventHandler(
		makePubsubEventHandler(string(grpc.SubscriptionEventServiceUnregister)),
	); err != nil {
		return err
	}

	serviceLog.Info(" InitDaprPubsubEvent ------ begin ------")
	return nil
}

func makePubsubEventHandler(name string) (
	string, func(ctx context.Context, e *common.TopicEvent) (retry bool, err error),
) {
	return name, func(ctx context.Context, e *common.TopicEvent) (retry bool, err error) {
		serviceLog.Info("UnRegister service [%s] env:%v", name, e.Data.(string))
		ctrlModel, _ := controller.GetControllerModel()
		ctrlModel.EventCallNoReturn(&component.ModelEventReq{
			EventType: name,
			Msg:       e,
		})
		return false, nil
	}
}
