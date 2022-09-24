package grpcPubsubEvent

import (
	"encoding/json"
	"game-message-core/grpc"
	"game-message-core/grpc/pubsubEventData"

	"github.com/Meland-Inc/game-services/src/common/daprInvoke"
	"github.com/Meland-Inc/game-services/src/common/serviceLog"
)

func RPCPubsubEventUseNft(env *pubsubEventData.UserUseNFTEvent) error {
	inputBytes, err := json.Marshal(env)
	if err != nil {
		serviceLog.Error("RPCPubsubEventUseNft Marshal Input failed err: %+v", err)
		return err
	}
	return daprInvoke.PubSubEventCall(string(grpc.SubscriptionEventUseNFT), string(inputBytes))
}
