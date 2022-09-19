package daprCalls

import (
	"game-message-core/grpc"

	"github.com/Meland-Inc/game-services/src/common/daprInvoke"
	message "github.com/Meland-Inc/game-services/src/global/web3Message"
)

func InitDaprCallHandle() (err error) {
	if err := initClientMsgCallHandle(); err != nil {
		return err
	}

	if err := initServiceGrpcCallHandle(); err != nil {
		return err
	}

	return nil
}

func initClientMsgCallHandle() error {
	if err := daprInvoke.AddServiceInvocationHandler(
		string(grpc.ProtoMessageActionPullClientMessage),
		ClientMessageHandler,
	); err != nil {
		return err
	}

	return nil
}

func initServiceGrpcCallHandle() error {
	if err := daprInvoke.AddServiceInvocationHandler(
		string(message.GameServiceActionDeductUserExp) ,
		Web3DeductUserExpHandler,
	); err != nil {
		return err
	}

	return nil
}
