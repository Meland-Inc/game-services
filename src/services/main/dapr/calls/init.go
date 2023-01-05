package daprCalls

import (
	"context"
	"fmt"

	"game-message-core/grpc"

	"github.com/Meland-Inc/game-services/src/common/daprInvoke"
	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	"github.com/Meland-Inc/game-services/src/global/module"
	message "github.com/Meland-Inc/game-services/src/global/web3Message"
	"github.com/dapr/go-sdk/service/common"
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
	return daprInvoke.AddServiceInvocationHandler(
		makeClientMsgHandler(string(grpc.ProtoMessageActionPullClientMessage)),
	)
}

func initServiceGrpcCallHandle() error {
	if err := daprInvoke.AddServiceInvocationHandler(makeServiceCallHandler(
		string(message.GameDataServiceActionDeductUserExp), module.MODULE_NAME_PLAYER_DATA,
	)); err != nil {
		return err
	}
	if err := daprInvoke.AddServiceInvocationHandler(makeServiceCallHandler(
		string(message.GameDataServiceActionGetPlayerInfoByUserId), module.MODULE_NAME_PLAYER_DATA,
	)); err != nil {
		return err
	}
	if err := daprInvoke.AddServiceInvocationHandler(makeServiceCallHandler(
		string(grpc.UserActionGetUserData), module.MODULE_NAME_PLAYER_DATA,
	)); err != nil {
		return err
	}

	if err := daprInvoke.AddServiceInvocationHandler(makeServiceCallHandler(
		string(grpc.MainServiceActionMintNFT), module.MODULE_NAME_PLAYER_DATA,
	)); err != nil {
		return err
	}
	if err := daprInvoke.AddServiceInvocationHandler(makeServiceCallHandler(
		string(grpc.MainServiceActionTakeNFT), module.MODULE_NAME_PLAYER_DATA,
	)); err != nil {
		return err
	}

	if err := daprInvoke.AddServiceInvocationHandler(makeServiceCallHandler(
		string(grpc.MainServiceActionGetAllBuild), module.MODULE_NAME_LAND,
	)); err != nil {
		return err
	}

	if err := daprInvoke.AddServiceInvocationHandler(makeServiceCallHandler(
		string(grpc.MainServiceActionGetHomeData), module.MODULE_NAME_HOME,
	)); err != nil {
		return err
	}

	return nil
}

func makeServiceCallHandler(name string, modelName string) (
	string, func(ctx context.Context, in *common.InvocationEvent) (*common.Content, error),
) {
	return name, func(ctx context.Context, in *common.InvocationEvent) (*common.Content, error) {
		model, exist := module.GetModel(modelName)
		if !exist {
			return nil, fmt.Errorf("model [%s] not found", modelName)
		}

		env := module.NewModuleEventReq(name, in.Data, false, nil)
		serviceLog.Info("receive [%s] env:%v", name, string(in.Data))
		resCh := model.EventCall(env)
		if resCh.GetError() != nil {
			return nil, resCh.GetError()
		}
		return daprInvoke.MakeOutputContent(in, resCh.GetResult())
	}
}
