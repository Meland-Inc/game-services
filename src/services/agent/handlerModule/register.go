package handlerModule

import (
	"game-message-core/grpc"

	"github.com/Meland-Inc/game-services/src/services/agent/handlerModule/serviceHandler"
)

func (p *HandlerModule) RegisterClientEvent() {
}

func (p *HandlerModule) RegisterGameServiceDaprCall() {
	p.AddGameServiceDaprCall(
		string(grpc.ProtoMessageActionBroadCastToClient),
		serviceHandler.BroadCastToClientHandler,
	)
	p.AddGameServiceDaprCall(
		string(grpc.ProtoMessageActionMultipleBroadCastToClient),
		serviceHandler.MultipleBroadCastToClientHandler,
	)

}

func (p *HandlerModule) RegisterGameServiceDaprEvent() {
	p.AddGameServiceDaprEvent(
		string(grpc.SubscriptionEventServiceUnregister),
		serviceHandler.GRPCServiceUnRegisterEvent,
	)
	p.AddGameServiceDaprEvent(
		string(grpc.SubscriptionEventTickOutPlayer),
		serviceHandler.GRPCTickOutUserEvent,
	)
	p.AddGameServiceDaprEvent(
		string(grpc.SubscriptionEventUserChangeService),
		serviceHandler.GRPCChangeServiceEvent,
	)
}

func (p *HandlerModule) RegisterWeb3DaprCall() {

}

func (p *HandlerModule) RegisterWeb3DaprEvent() {

}
