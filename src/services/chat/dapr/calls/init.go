package daprCalls

import (
	"game-message-core/grpc"

	"github.com/Meland-Inc/game-services/src/common/daprInvoke"
)

func InitDaprCallHandle() (err error) {
	if err := initClientMsgCallHandle(); err != nil {
		return err
	}

	return nil
}

func initClientMsgCallHandle() error {
	if err := daprInvoke.AddServiceInvocationHandler(
		string(grpc.UserActionLeaveGame),
		UserLeaveGameHandler,
	); err != nil {
		return err
	}

	return nil
}
