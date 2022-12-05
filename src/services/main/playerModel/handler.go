package playerModel

import (
	"game-message-core/grpc"

	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	"github.com/Meland-Inc/game-services/src/global/component"
	message "github.com/Meland-Inc/game-services/src/global/web3Message"
)

func (p *PlayerDataModel) OnEvent(env *component.ModelEventReq, curMs int64) {
	defer func() {
		err := recover()
		if err != nil {
			serviceLog.StackError("ControllerModel.onEvent err: %v", err)
		}
	}()

	switch env.EventType {
	case string(grpc.ProtoMessageActionPullClientMessage):
		p.clientMsgHandler(env, curMs)

	case string(message.GameDataServiceActionDeductUserExp):
		p.Web3DeductUserExpHandler(env, curMs)
	case string(message.GameDataServiceActionGetPlayerInfoByUserId):
		p.Web3GetPlayerDataHandler(env, curMs)

	case string(message.SubscriptionEventUpdateUserNFT):
		p.Web3UpdateUserNftEvent(env, curMs)
	case string(message.SubscriptionEventMultiUpdateUserNFT):
		p.Web3MultiUpdateUserNftEvent(env, curMs)

	case string(grpc.UserActionGetUserData):
		p.GRPCGetUserDataHandler(env, curMs)
	case string(grpc.MainServiceActionTakeNFT):
		p.GRPCTakeUserNftHandler(env, curMs)

	case string(grpc.SubscriptionEventUserEnterGame):
		p.GRPCUserEnterGameEvent(env, curMs)
	case string(grpc.SubscriptionEventUserLeaveGame):
		p.GRPCUserLeaveGameEvent(env, curMs)
	case string(grpc.SubscriptionEventSavePlayerData):
		p.GRPCSavePlayerDataEvent(env, curMs)
	case string(grpc.SubscriptionEventKillMonster):
		p.GRPCKillMonsterEvent(env, curMs)
	case string(grpc.SubscriptionEventPlayerDeath):
		p.GRPCPlayerDeathEvent(env, curMs)
	case string(grpc.SubscriptionEventUserTaskReward):
		p.GRPCUserTaskRewardEvent(env, curMs)
	}

}
