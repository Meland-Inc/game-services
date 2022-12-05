package daprCalls

import (
	"context"
	"fmt"

	"game-message-core/grpc"

	"github.com/Meland-Inc/game-services/src/common/daprInvoke"
	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	"github.com/Meland-Inc/game-services/src/global/component"
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
	if err := daprInvoke.AddServiceInvocationHandler(
		makeServiceCallHandler(
			string(message.GameDataServiceActionDeductUserExp), component.MODEL_NAME_PLAYER_DATA,
		),
	); err != nil {
		return err
	}
	if err := daprInvoke.AddServiceInvocationHandler(
		makeServiceCallHandler(
			string(message.GameDataServiceActionGetPlayerInfoByUserId), component.MODEL_NAME_PLAYER_DATA,
		),
	); err != nil {
		return err
	}
	if err := daprInvoke.AddServiceInvocationHandler(
		makeServiceCallHandler(
			string(grpc.UserActionGetUserData), component.MODEL_NAME_PLAYER_DATA,
		),
	); err != nil {
		return err
	}

	if err := daprInvoke.AddServiceInvocationHandler(
		makeServiceCallHandler(
			string(grpc.MainServiceActionTakeNFT), component.MODEL_NAME_PLAYER_DATA,
		),
	); err != nil {
		return err
	}

	if err := daprInvoke.AddServiceInvocationHandler(
		makeServiceCallHandler(
			string(grpc.MainServiceActionGetAllBuild), component.MODEL_NAME_LAND,
		),
	); err != nil {
		return err
	}

	return nil
}

func makeServiceCallHandler(name string, modelName string) (
	string, func(ctx context.Context, in *common.InvocationEvent) (*common.Content, error),
) {
	return name, func(ctx context.Context, in *common.InvocationEvent) (*common.Content, error) {
		model, exist := component.GetInstance().GetModel(modelName)
		if !exist {
			return nil, fmt.Errorf("model [%s] not found", modelName)
		}

		env := &component.ModelEventReq{
			EventType: name,
			Msg:       in.Data,
		}
		serviceLog.Info("receive [%s] env:%v", name, string(in.Data))
		resCh := model.EventCall(env)
		if resCh.Err != nil {
			return nil, resCh.Err
		}
		return daprInvoke.MakeOutputContent(in, resCh.Result)
	}
}
