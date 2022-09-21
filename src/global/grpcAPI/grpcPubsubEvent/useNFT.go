package grpcPubsubEvent

import (
	"encoding/json"
	"game-message-core/grpc"
	"game-message-core/grpc/pubsubEventData"

	"github.com/Meland-Inc/game-services/src/common/daprInvoke"
)

func RPCPubsubEventUseNft(env pubsubEventData.UserUseNFTEvent) error {
	bs, err := json.Marshal(env)
	if err != nil {
		return err
	}
	return daprInvoke.PubSubEventCall(string(grpc.SubscriptionEventUseNFT), string(bs))
}
