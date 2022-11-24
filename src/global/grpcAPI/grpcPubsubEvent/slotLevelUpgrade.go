package grpcPubsubEvent

import (
	"encoding/json"
	"game-message-core/grpc"
	"game-message-core/grpc/pubsubEventData"

	"github.com/Meland-Inc/game-services/src/common/daprInvoke"
	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	"github.com/Meland-Inc/game-services/src/common/time_helper"
)

func RPCPubsubEventSlotLevelUpgrade(userId int64, slotPos, level int32) error {
	env := pubsubEventData.SlotLevelUpgradeEvent{
		MsgVersion: time_helper.NowUTCMill(),
		UserId:     userId,
		SlotPos:    slotPos,
		Level:      level,
	}
	inputBytes, err := json.Marshal(env)
	if err != nil {
		serviceLog.Error("SlotLevelUpgradeEvent Marshal Input failed err: %+v", err)
		return err
	}

	serviceLog.Info("pubsubEvent Slot Level Upgrade Event : %+v", env)

	return daprInvoke.PubSubEventCall(string(grpc.SubscriptionEventSlotLevelUpgrade), string(inputBytes))
}
