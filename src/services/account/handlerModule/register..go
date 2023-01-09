package handlerModule

import (
	"game-message-core/proto"

	"github.com/Meland-Inc/game-services/src/services/account/handlerModule/clientHandler"
)

func (p *HandlerModule) RegisterClientEvent() {
	// player message
	p.AddClientEvent(proto.EnvelopeType_CreatePlayer, clientHandler.CreatePlayerHandler)
	p.AddClientEvent(proto.EnvelopeType_QueryPlayer, clientHandler.QueryPlayerHandler)

}

func (p *HandlerModule) RegisterGameServiceDaprCall() {

}

func (p *HandlerModule) RegisterGameServiceDaprEvent() {

}

func (p *HandlerModule) RegisterWeb3DaprCall() {

}

func (p *HandlerModule) RegisterWeb3DaprEvent() {

}
