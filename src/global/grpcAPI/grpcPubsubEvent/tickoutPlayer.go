package grpcPubsubEvent

import (
	"encoding/json"
	"game-message-core/grpc"
	"game-message-core/grpc/pubsubEventData"
	"game-message-core/proto"

	"github.com/Meland-Inc/game-services/src/common/daprInvoke"
	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	"github.com/Meland-Inc/game-services/src/common/time_helper"
)

func RPCPubsubEventTickOutPlayer(
	userId int64,
	agentAppId, socketId, sceneServiceAppId string,
	code proto.TickOutType,
) error {
	env := pubsubEventData.TickOutPlayerEvent{
		MsgVersion:        time_helper.NowUTCMill(),
		UserId:            userId,
		AgentAppId:        agentAppId,
		SocketId:          socketId,
		SceneServiceAppId: sceneServiceAppId,
		TickOutCode:       code,
	}
	inputBytes, err := json.Marshal(env)
	if err != nil {
		serviceLog.Error("TickOutPlayerEvent Marshal failed err: %+v", err)
		return err
	}

	serviceLog.Info("Call TickOutPlayerEvent %+v", env)
	return daprInvoke.PubSubEventCall(string(grpc.SubscriptionEventTickOutPlayer), string(inputBytes))
}
