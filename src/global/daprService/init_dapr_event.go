package daprService

import (
	"context"

	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	"github.com/Meland-Inc/game-services/src/global/contract"
	"github.com/Meland-Inc/game-services/src/global/globalModule"
	"github.com/Meland-Inc/game-services/src/global/module"
	"github.com/dapr/go-sdk/service/common"

	"github.com/Meland-Inc/game-services/src/common/daprInvoke"
)

func InitDaprPubsubEvent() (err error) {
	serEventModel, err := globalModule.GetServiceEventModel()
	if err != nil {
		return err
	}

	for _, eventName := range serEventModel.GetWeb3DaprEventTypes() {
		err = daprInvoke.AddTopicEventHandler(makePubsubEventHandle(eventName, serEventModel))
		if err != nil {
			return err
		}
	}

	for _, eventName := range serEventModel.GetGameServiceDaprEventTypes() {
		err = daprInvoke.AddTopicEventHandler(makePubsubEventHandle(eventName, serEventModel))
		if err != nil {
			return err
		}
	}

	return nil
}

func makePubsubEventHandle(name string, serEventModel contract.IServiceEvent) (
	string, func(ctx context.Context, e *common.TopicEvent) (retry bool, err error),
) {
	serviceLog.Info("listen dapr pubsubEvent [ %s ]", name)

	handler := func(ctx context.Context, e *common.TopicEvent) (retry bool, err error) {
		// 	serviceLog.Debug("pubsub event[%s]  data:%v", name, e.Data)

		serEventModel.EventCallNoReturn(module.NewModuleEventReq(name, e, false, nil))
		return false, nil
	}
	return name, handler
}
