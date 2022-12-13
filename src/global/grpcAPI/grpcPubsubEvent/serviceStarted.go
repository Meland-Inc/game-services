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

func RPCPubsubEventServiceStarted(ser base_data.ServiceData) error {
	env := pubsubEventData.ServiceStartedEvent{
		MsgVersion: time_helper.NowUTCMill(),
		Service:    ser,
	}

	inputBytes, err := json.Marshal(env)
	if err != nil {
		serviceLog.Error("ServiceStartedEvent Marshal failed err: %+v", err)
		return err
	}

	serviceLog.Info("pubsubEvent ServiceStartedEvent: %+v", env)

	return daprInvoke.PubSubEventCall(string(grpc.SubscriptionEventServiceStarted), string(inputBytes))
}
