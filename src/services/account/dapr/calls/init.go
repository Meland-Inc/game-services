package daprCalls

import (
	"game-message-core/grpc"

	"github.com/Meland-Inc/game-services/src/common/daprInvoke"
)

func InitDaprCallHandle() (err error) {
	if daprInvoke.AddServiceInvocationHandler(
		string(grpc.ProtoMessageActionPullClientMessage),
		ClientMessageHandler,
	); err != nil {
		return err
	}

	return nil
}
