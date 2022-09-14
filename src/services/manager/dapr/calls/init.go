package daprCalls

import (
	"game-message-core/grpc"

	"github.com/Meland-Inc/game-services/src/common/daprInvoke"
)

func InitDaprCallHandle() (err error) {
	if daprInvoke.AddServiceInvocationHandler(
		string(grpc.ManagerServiceActionRegister), RegisterServiceHandler,
	); err != nil {
		return err
	}
	if daprInvoke.AddServiceInvocationHandler(
		string(grpc.ManagerServiceActionDestroy), DestroyServiceHandler,
	); err != nil {
		return err
	}

	return nil
}
