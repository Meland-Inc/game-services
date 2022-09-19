package grpcPubsubEvent

import (
	"game-message-core/grpc"
	"game-message-core/grpc/pubsubEventData"

	"github.com/Meland-Inc/game-services/src/common/daprInvoke"
)

func RPCPubsubEventEnterGame(env pubsubEventData.UserEnterGameEvent) error {
	return daprInvoke.PubSubEventCall(string(grpc.SubscriptionEventUserEnterGame), env)
}
