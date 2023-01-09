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

}

func (p *HandlerModule) RegisterGameServiceDaprEvent() {

}

func (p *HandlerModule) RegisterWeb3DaprCall() {

}

func (p *HandlerModule) RegisterWeb3DaprEvent() {

}
