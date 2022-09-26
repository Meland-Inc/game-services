package grpcPubsubEvent

import (
	"encoding/json"
	"game-message-core/grpc"
	"game-message-core/grpc/pubsubEventData"

	"github.com/Meland-Inc/game-services/src/common/daprInvoke"
	"github.com/Meland-Inc/game-services/src/common/serviceLog"
)

func RPCPubsubEventTaskReward(env *pubsubEventData.UserTaskRewardEvent) error {
	inputBytes, err := json.Marshal(env)
	if err != nil {
		serviceLog.Error("UserTaskRewardEvent Marshal Input failed err: %+v", err)
		return err
	}
	return daprInvoke.PubSubEventCall(string(grpc.SubscriptionEventUserTaskReward), string(inputBytes))
}
