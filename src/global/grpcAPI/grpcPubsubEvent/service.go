package grpcPubsubEvent

import (
	"encoding/json"
	"game-message-core/grpc"
	base_data "game-message-core/grpc/baseData"
	"game-message-core/grpc/pubsubEventData"

	"github.com/Meland-Inc/game-services/src/common/daprInvoke"
	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	"github.com/Meland-Inc/game-services/src/common/time_helper"
)

func RPCPubsubEventServiceUnregister(service base_data.ServiceData) error {
	env := pubsubEventData.ServiceUnregisterEvent{
		MsgVersion: time_helper.NowUTCMill(),
		Service:    service,
	}

	inputBytes, err := json.Marshal(env)
	if err != nil {
		serviceLog.Error("RPCPubsub Unregister service Marshal failed err: %+v", err)
		return err
	}

	serviceLog.Info("pubsub event service unregister %+v", service)

	return daprInvoke.PubSubEventCall(string(grpc.SubscriptionEventServiceUnregister), string(inputBytes))
}
