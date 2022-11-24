package grpcPubsubEvent

import (
	"encoding/json"
	"game-message-core/grpc"
	"game-message-core/grpc/pubsubEventData"

	"github.com/Meland-Inc/game-services/src/common/daprInvoke"
	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	"github.com/Meland-Inc/game-services/src/common/time_helper"
)

func RPCPubsubEventUserLevelUpgrade(userId int64, level int32) error {
	env := pubsubEventData.UserLevelUpgradeEvent{
		MsgVersion: time_helper.NowUTCMill(),
		UserId:     userId,
		Level:      level,
	}
	inputBytes, err := json.Marshal(env)
	if err != nil {
		serviceLog.Error("UserLevelUpgradeEvent Marshal Input failed err: %+v", err)
		return err
	}

	serviceLog.Info("pubsubEvent userLevelUp: %+v", env)

	return daprInvoke.PubSubEventCall(string(grpc.SubscriptionEventUserLevelUpgrade), string(inputBytes))
}
