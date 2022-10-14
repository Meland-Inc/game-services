package grpcPubsubEvent

import (
	"encoding/json"
	"game-message-core/grpc"
	"game-message-core/grpc/pubsubEventData"

	"github.com/Meland-Inc/game-services/src/common/daprInvoke"
	"github.com/Meland-Inc/game-services/src/common/serviceLog"
)

func RPCPubsubEventEnterGame(env *pubsubEventData.UserEnterGameEvent) error {
	inputBytes, err := json.Marshal(env)
	if err != nil {
		serviceLog.Error("RPCPubsubEventEnterGame Marshal Input failed err: %+v", err)
		return err
	}

	serviceLog.Info(
		"[%s] CallEnterGame user[%d], socketId[%s], sceneService[%s]",
		env.AgentAppId, env.UserId, env.UserSocketId, env.SceneServiceAppId,
	)

	return daprInvoke.PubSubEventCall(string(grpc.SubscriptionEventUserEnterGame), string(inputBytes))
}
