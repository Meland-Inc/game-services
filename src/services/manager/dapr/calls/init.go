package daprCalls

import (
	"game-message-core/grpc"

	"github.com/Meland-Inc/game-services/src/common/daprInvoke"
	"github.com/Meland-Inc/game-services/src/common/serviceLog"
)

func InitDaprCallHandle() (err error) {
	serviceLog.Info(" InitDaprCallHandle ------ begin ------")
	if err = daprInvoke.AddServiceInvocationHandler(
		string(grpc.ManagerServiceActionRegister),
		RegisterServiceHandler,
	); err != nil {
		return err
	}

	if err = daprInvoke.AddServiceInvocationHandler(
		string(grpc.ManagerServiceActionSelectService),
		SelectServiceHandler,
	); err != nil {
		return err
	}

	if err = daprInvoke.AddServiceInvocationHandler(
		string(grpc.ManagerServiceActionMultiSelectService),
		MultiSelectServiceHandler,
	); err != nil {
		return err
	}

	serviceLog.Info(" InitDaprCallHandle ------ end ------")
	return nil
}
