package daprEvent

import (
	"game-message-core/grpc"

	"github.com/Meland-Inc/game-services/src/common/daprInvoke"
	"github.com/Meland-Inc/game-services/src/common/serviceLog"
)

func InitDaprPubsubEvent() (err error) {
	serviceLog.Info(" InitDaprPubsubEvent ------ begin ------")
	daprInvoke.AddTopicEventHandler(
		string(grpc.SubscriptionEventServiceUnregister),
		ServiceUnRegisterHandler,
	)
	if err != nil {
		return err
	}
	serviceLog.Info(" InitDaprPubsubEvent ------ begin ------")
	return nil
}
