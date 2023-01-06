package handlerModule

import (
	"game-message-core/grpc"

	"github.com/Meland-Inc/game-services/src/services/manager/handlerModule/serviceHandler"
)

func (p *HandlerModule) RegisterClientEvent() {
}

func (p *HandlerModule) RegisterGameServiceDaprCall() {
	p.AddGameServiceDaprCall(
		string(grpc.ManagerServiceActionRegister),
		serviceHandler.GRPCServiceRegisterHandler,
	)

	p.AddGameServiceDaprCall(
		string(grpc.ManagerServiceActionSelectService),
		serviceHandler.GRPCServiceSelectHandler,
	)

	p.AddGameServiceDaprCall(
		string(grpc.ManagerServiceActionMultiSelectService),
		serviceHandler.GRPCMultiSelectServiceHandler,
	)

	p.AddGameServiceDaprCall(
		string(grpc.ManagerServiceActionStartService),
		serviceHandler.GRPCServiceStartHandler,
	)
}

func (p *HandlerModule) RegisterGameServiceDaprEvent() {
	p.AddGameServiceDaprEvent(
		string(grpc.SubscriptionEventServiceUnregister),
		serviceHandler.GRPCServiceUnregisterEvent,
	)

}

func (p *HandlerModule) RegisterWeb3DaprCall() {

}

func (p *HandlerModule) RegisterWeb3DaprEvent() {
}
