package daprCalls

import (
	"game-message-core/grpc"

	"github.com/Meland-Inc/game-services/src/common/daprInvoke"
)

func InitDaprCallHandle() (err error) {
	if err := initClientMsgCallHandle(); err != nil {
		return err
	}

	if err := initServiceGrpcCallHandle(); err != nil {
		return err
	}

	if err := initServiceGrpcPubsubEventHandle(); err != nil {
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
	// if err := daprInvoke.AddServiceInvocationHandler(
	// 	string(grpc.ProtoMessageActionPullClientMessage),
	// 	ClientMessageHandler,
	// ); err != nil {
	// 	return err
	// }

	return nil
}

func initServiceGrpcPubsubEventHandle() error {
	// if err := daprInvoke.PubSubEventCall(
	// 	string(grpc.ProtoMessageActionPullClientMessage),
	// 	ClientMessageHandler,
	// ); err != nil {
	// 	return err
	// }

	return nil
}
