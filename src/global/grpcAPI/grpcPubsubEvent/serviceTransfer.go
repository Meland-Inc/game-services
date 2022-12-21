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

func RPCPubsubEventServiceTransfer(ser base_data.ServiceData) error {
	env := pubsubEventData.ServiceTransferEvent{
		MsgVersion: time_helper.NowUTCMill(),
		Service:    ser,
	}

	inputBytes, err := json.Marshal(env)
	if err != nil {
		serviceLog.Error("ServiceTransferEvent Marshal failed err: %+v", err)
		return err
	}

	serviceLog.Info("pubsubEvent ServiceTransferEvent: %+v", env)

	return daprInvoke.PubSubEventCall(string(grpc.SubscriptionEventServiceTransfer), string(inputBytes))
}
