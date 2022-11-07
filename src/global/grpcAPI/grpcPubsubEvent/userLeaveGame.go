package grpcPubsubEvent

import (
	"encoding/json"
	"game-message-core/grpc"
	"game-message-core/grpc/pubsubEventData"

	"github.com/Meland-Inc/game-services/src/common/daprInvoke"
	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	"github.com/Meland-Inc/game-services/src/common/time_helper"
	"github.com/Meland-Inc/game-services/src/global/serviceCnf"
)

func RPCPubsubEventLeaveGame(userId int64) error {
	env := pubsubEventData.UserLeaveGameEvent{
		MsgVersion: time_helper.NowUTCMill(),
		AgentAppId: serviceCnf.GetInstance().AppId,
		UserId:     userId,
	}
	inputBytes, err := json.Marshal(env)
	if err != nil {
		serviceLog.Error("RPCPubsubEventEnterGame Marshal Input failed err: %+v", err)
		return err
	}

	serviceLog.Info("CallLeaveGame %+v", env)
	return daprInvoke.PubSubEventCall(string(grpc.SubscriptionEventUserLeaveGame), string(inputBytes))
}
