package handlerModule

import (
	"game-message-core/grpc"
	"game-message-core/proto"

	"github.com/Meland-Inc/game-services/src/services/chat/handlerModule/clientHandler"
	"github.com/Meland-Inc/game-services/src/services/chat/handlerModule/serviceHandler"
)

func (p *HandlerModule) RegisterClientEvent() {
	p.AddClientEvent(proto.EnvelopeType_SendChatMessage, clientHandler.ChatMsgHandle)

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
		string(grpc.SubscriptionEventSavePlayerData),
		serviceHandler.GRPCSavePlayerDataEvent,
	)
	p.AddGameServiceDaprEvent(
		string(grpc.SubscriptionEventTickOutPlayer),
		serviceHandler.GRPCTickOutUserEvent,
	)
}

func (p *HandlerModule) RegisterWeb3DaprCall() {

}

func (p *HandlerModule) RegisterWeb3DaprEvent() {

}
