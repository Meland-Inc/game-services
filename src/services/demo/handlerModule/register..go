package handlerModule

import (
	"game-message-core/grpc"
	"game-message-core/proto"

	message "github.com/Meland-Inc/game-services/src/global/web3Message"
	"github.com/Meland-Inc/game-services/src/services/demo/handlerModule/clientHandler"
	"github.com/Meland-Inc/game-services/src/services/demo/handlerModule/serviceHandler"
	"github.com/Meland-Inc/game-services/src/services/demo/handlerModule/web3Handler"
)

func (p *HandlerModule) RegisterClientEvent() {
	// sing in message
	p.AddClientEvent(proto.EnvelopeType_SigninPlayer, clientHandler.SingInHandler)

}

func (p *HandlerModule) RegisterGameServiceDaprCall() {
	p.AddGameServiceDaprCall(
		string(grpc.MainServiceActionGetHomeData),
		serviceHandler.GRPCGetHomeDataHandler,
	)
}

func (p *HandlerModule) RegisterGameServiceDaprEvent() {
	p.AddGameServiceDaprEvent(
		string(grpc.SubscriptionEventSaveHomeData),
		serviceHandler.GRPCSaveHomeDataEvent,
	)

}

func (p *HandlerModule) RegisterWeb3DaprCall() {
	p.AddWeb3DaprCall(
		string(message.GameDataServiceActionDeductUserExp),
		web3Handler.Web3DeductUserExpHandler,
	)

}

func (p *HandlerModule) RegisterWeb3DaprEvent() {
	p.AddWeb3DaprEvent(
		string(message.SubscriptionEventUpdateUserNFT),
		web3Handler.Web3UpdateUserNftEvent,
	)

}
