package msgChannel

import (
	"game-message-core/grpc"

	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	"github.com/Meland-Inc/game-services/src/services/task/msgChannel/serviceMsgHandle"
)

type ServiceMsgData struct {
	MsgId   string      `json:"msgId"`
	MsgBody interface{} `json:"msgBody"`
}

func (ch *MsgChannel) onServiceMessage(input *ServiceMsgData) {
	serviceLog.Info("received service[%v] message: %v", input.MsgId, input.MsgBody)

	switch input.MsgId {
	case string(grpc.SubscriptionEventUserEnterGame):
		serviceMsgHandle.UserEnterGameHandle(input.MsgBody)

	case string(grpc.SubscriptionEventUserLeaveGame):
		serviceMsgHandle.PlayerLeaveGameHandler(input.MsgBody)

	case string(grpc.SubscriptionEventKillMonster):
		serviceMsgHandle.KillMonsterHandler(input.MsgBody)

	case string(grpc.SubscriptionEventUseNFT):
		serviceMsgHandle.PlayerUseItemHandler(input.MsgBody)

	case string(grpc.SubscriptionEventUserLevelUpgrade):
		serviceMsgHandle.UserLevelUpgradeHandler(input.MsgBody)

	case string(grpc.SubscriptionEventSlotLevelUpgrade):
		serviceMsgHandle.SlotLevelUpgradeHandler(input.MsgBody)

	case string(grpc.SubscriptionEventTaskFinish):
		serviceMsgHandle.TaskFinishHandler(input.MsgBody)

	case string(grpc.SubscriptionEventTaskListFinish):
		serviceMsgHandle.TaskListFinishHandler(input.MsgBody)

	}
}
