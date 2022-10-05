package msgChannel

import (
	"game-message-core/grpc"

	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	message "github.com/Meland-Inc/game-services/src/global/web3Message"
	"github.com/Meland-Inc/game-services/src/services/main/msgChannel/serviceMsgHandle"
)

type ServiceMsgData struct {
	MsgId   string      `json:"msgId"`
	MsgBody interface{} `json:"msgBody"`
}

func (ch *MsgChannel) onServiceMessage(input *ServiceMsgData) {
	serviceLog.Info("received service[%v] message: %+v", input.MsgId, input.MsgBody)

	switch input.MsgId {
	case string(grpc.SubscriptionEventUserEnterGame):
		serviceMsgHandle.UserEnterGameHandle(input.MsgBody)
	case string(message.GameServiceActionDeductUserExp):
		serviceMsgHandle.Web3DeductUserExpHandler(input.MsgBody)
	case string(message.SubscriptionEventUpdateUserNFT):
		serviceMsgHandle.Web3UpdateUserNftHandler(input.MsgBody)
	case string(message.SubscriptionEventMultiUpdateUserNFT):
		serviceMsgHandle.Web3MultiUpdateUserNftHandler(input.MsgBody)

	case string(grpc.SubscriptionEventSavePlayerData):
		serviceMsgHandle.SavePlayerDataHandler(input.MsgBody)
	case string(grpc.SubscriptionEventKillMonster):
		serviceMsgHandle.KillMonsterHandler(input.MsgBody)
	case string(grpc.SubscriptionEventPlayerDeath):
		serviceMsgHandle.PlayerDeathHandler(input.MsgBody)

	case string(grpc.UserActionLeaveGame):
		serviceMsgHandle.PlayerLeaveGameHandler(input.MsgBody)

	}
}
