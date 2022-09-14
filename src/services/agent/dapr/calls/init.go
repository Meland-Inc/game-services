package daprCalls

import (
	"game-message-core/grpc"

	"github.com/Meland-Inc/game-services/src/common/daprInvoke"
)

func InitDaprCallHandle() (err error) {
	if err = daprInvoke.AddServiceInvocationHandler(
		string(grpc.ProtoMessageActionBroadCastToClient),
		BroadCastToClientHandler,
	); err != nil {
		return err
	}

	if err = daprInvoke.AddServiceInvocationHandler(
		string(grpc.ProtoMessageActionMultipleBroadCastToClient),
		MultipleBroadCastToClientHandler,
	); err != nil {
		return err
	}

	return nil
}
