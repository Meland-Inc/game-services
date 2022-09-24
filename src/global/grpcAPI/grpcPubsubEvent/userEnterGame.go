package grpcPubsubEvent

import (
	"game-message-core/grpc"
	"game-message-core/proto"
	"game-message-core/protoTool"

	"github.com/Meland-Inc/game-services/src/common/daprInvoke"
	"github.com/Meland-Inc/game-services/src/common/serviceLog"
)

func RPCPubsubEventEnterGame(env *proto.UserEnterGameEvent) error {
	inputBytes, err := protoTool.MarshalProto(env)
	if err != nil {

		serviceLog.Error("RPCPubsubEventEnterGame Marshal Input failed err: %+v", err)
		return err
	}
	return daprInvoke.PubSubEventCall(string(grpc.SubscriptionEventUserEnterGame), string(inputBytes))
}
