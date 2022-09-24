package grpcPubsubEvent

import (
	"game-message-core/grpc"
	"game-message-core/proto"
	"game-message-core/protoTool"

	"github.com/Meland-Inc/game-services/src/common/daprInvoke"
	"github.com/Meland-Inc/game-services/src/common/serviceLog"
)

func RPCPubsubEventUseNft(env *proto.UserUseNFTEvent) error {
	inputBytes, err := protoTool.MarshalProto(env)
	if err != nil {
		serviceLog.Error("RPCPubsubEventUseNft Marshal Input failed err: %+v", err)
		return err
	}
	return daprInvoke.PubSubEventCall(string(grpc.SubscriptionEventUseNFT), string(inputBytes))
}
