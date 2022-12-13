package daprCalls

import (
	"context"
	"game-message-core/grpc"

	"github.com/Meland-Inc/game-services/src/common/daprInvoke"
	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	"github.com/Meland-Inc/game-services/src/global/component"
	"github.com/Meland-Inc/game-services/src/services/manager/controller"
	"github.com/dapr/go-sdk/service/common"
)

func InitDaprCallHandle() (err error) {
	serviceLog.Info(" InitDaprCallHandle ------ begin ------")
	if err = daprInvoke.AddServiceInvocationHandler(
		makeCallHandler(string(grpc.ManagerServiceActionRegister)),
	); err != nil {
		return err
	}
	if err = daprInvoke.AddServiceInvocationHandler(
		makeCallHandler(string(grpc.ManagerServiceActionSelectService)),
	); err != nil {
		return err
	}
	if err = daprInvoke.AddServiceInvocationHandler(
		makeCallHandler(string(grpc.ManagerServiceActionMultiSelectService)),
	); err != nil {
		return err
	}
	if err = daprInvoke.AddServiceInvocationHandler(
		makeCallHandler(string(grpc.ManagerServiceActionStartService)),
	); err != nil {
		return err
	}

	serviceLog.Info(" InitDaprCallHandle ------ end ------")
	return nil
}

func makeCallHandler(name string) (string, func(ctx context.Context, in *common.InvocationEvent) (*common.Content, error)) {
	return name, func(ctx context.Context, in *common.InvocationEvent) (*common.Content, error) {
		ctrlModel, _ := controller.GetControllerModel()
		env := &component.ModelEventReq{
			EventType: name,
			Msg:       in.Data,
		}
		// serviceLog.Info("receive [%s] env:%v", name, string(in.Data))
		resCh := ctrlModel.EventCall(env)

		if resCh.Err != nil {
			return nil, resCh.Err
		}
		return daprInvoke.MakeOutputContent(in, resCh.Result)
	}
}
