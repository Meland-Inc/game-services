package grpcPubsubEvent

import (
	"game-message-core/grpc"
	"game-message-core/grpc/pubsubEventData"

	"github.com/Meland-Inc/game-services/src/common/daprInvoke"
)

func RPCPubsubEventUseNft(env pubsubEventData.UserUseNFTEvent) error {
	return daprInvoke.PubSubEventCall(string(grpc.SubscriptionEventUseNFT), env)
}
