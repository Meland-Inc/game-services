package handlerModule

import (
	"game-message-core/grpc"
	"game-message-core/proto"

	"github.com/Meland-Inc/game-services/src/services/task/handlerModule/clientHandler"
	"github.com/Meland-Inc/game-services/src/services/task/handlerModule/serviceHandler"
)

func (p *HandlerModule) RegisterClientEvent() {
	p.AddClientEvent(proto.EnvelopeType_SelfTasks, clientHandler.SelfTasksHandler)
	p.AddClientEvent(proto.EnvelopeType_AcceptTask, clientHandler.AcceptTaskHandler)
	p.AddClientEvent(proto.EnvelopeType_AbandonmentTask, clientHandler.AbandonmentTaskHandler)
	p.AddClientEvent(proto.EnvelopeType_UpgradeTaskProgress, clientHandler.UpgradeTaskProgressHandler)
	p.AddClientEvent(proto.EnvelopeType_TaskReward, clientHandler.TaskRewardHandler)
	p.AddClientEvent(proto.EnvelopeType_TaskListReward, clientHandler.TaskListRewardHandler)

}

func (p *HandlerModule) RegisterGameServiceDaprCall() {

}

func (p *HandlerModule) RegisterGameServiceDaprEvent() {
	p.AddGameServiceDaprEvent(
		string(grpc.SubscriptionEventUserEnterGame),
		serviceHandler.GRPCUserEnterGameEvent,
	)
	p.AddGameServiceDaprEvent(
		string(grpc.SubscriptionEventUserLeaveGame),
		serviceHandler.GRPCUserLeaveGameEvent,
	)
	p.AddGameServiceDaprEvent(
		string(grpc.SubscriptionEventUserChangeService),
		serviceHandler.GRPCUserChangeServiceEvent,
	)

	p.AddGameServiceDaprEvent(
		string(grpc.SubscriptionEventKillMonster),
		serviceHandler.GRPCKillMonsterEvent,
	)
	p.AddGameServiceDaprEvent(
		string(grpc.SubscriptionEventUseNFT),
		serviceHandler.GRPCPlayerUseItemEvent,
	)
	p.AddGameServiceDaprEvent(
		string(grpc.SubscriptionEventUserLevelUpgrade),
		serviceHandler.GRPCUserLevelUpgradeEvent,
	)
	p.AddGameServiceDaprEvent(
		string(grpc.SubscriptionEventSlotLevelUpgrade),
		serviceHandler.GRPCSlotLevelUpgradeEvent,
	)
	p.AddGameServiceDaprEvent(
		string(grpc.SubscriptionEventTaskFinish),
		serviceHandler.GRPCTaskFinishEvent,
	)
	p.AddGameServiceDaprEvent(
		string(grpc.SubscriptionEventTaskListFinish),
		serviceHandler.GRPCTaskListFinishEvent,
	)

}

func (p *HandlerModule) RegisterWeb3DaprCall() {

}

func (p *HandlerModule) RegisterWeb3DaprEvent() {

}
