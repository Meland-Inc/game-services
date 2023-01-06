package daprService

import (
	"context"

	"game-message-core/grpc"
	"game-message-core/grpc/methodData"
	"game-message-core/proto"

	"github.com/Meland-Inc/game-services/src/common/daprInvoke"
	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	"github.com/Meland-Inc/game-services/src/global/contract"
	"github.com/Meland-Inc/game-services/src/global/globalModule"
	"github.com/Meland-Inc/game-services/src/global/grpcAPI/grpcNetTool"
	"github.com/Meland-Inc/game-services/src/global/module"
	"github.com/dapr/go-sdk/service/common"
)

func InitDaprCallHandle() (err error) {
	serEventModel, err := globalModule.GetServiceEventModel()
	if err != nil {
		return err
	}

	if err = daprInvoke.AddServiceInvocationHandler(
		makeClientMsgHandler(serEventModel, string(grpc.ProtoMessageActionPullClientMessage)),
	); err != nil {
		return err
	}

	for _, eventName := range serEventModel.GetWeb3DaprCallTypes() {
		err = daprInvoke.AddServiceInvocationHandler(makeServiceCallHandle(serEventModel, eventName))
		if err != nil {
			return err
		}
	}

	for _, eventName := range serEventModel.GetGameServiceDaprCallTypes() {
		err = daprInvoke.AddServiceInvocationHandler(makeServiceCallHandle(serEventModel, eventName))
		if err != nil {
			return err
		}
	}

	return nil
}

func makeServiceCallHandle(serEventModel contract.IServiceEvent, name string) (
	string, func(ctx context.Context, in *common.InvocationEvent) (*common.Content, error),
) {
	serviceLog.Info("listen dapr calls [ %s ]", name)

	handler := func(ctx context.Context, in *common.InvocationEvent) (*common.Content, error) {
		env := module.NewModuleEventReq(name, in.Data, false, nil)
		serviceLog.Debug("receive dapr call [%s] env:%v", name, string(in.Data))
		resCh := serEventModel.EventCall(env)
		if resCh.GetError() != nil {
			return nil, resCh.GetError()
		}
		return daprInvoke.MakeOutputContent(in, resCh.GetResult())
	}
	return name, handler
}

func makeClientMsgHandler(serEventModel contract.IServiceEvent, name string) (
	string, func(ctx context.Context, in *common.InvocationEvent) (*common.Content, error),
) {
	serviceLog.Info("listen client dapr calls [ %s ]", name)

	return name, func(ctx context.Context, in *common.InvocationEvent) (*common.Content, error) {
		input := &methodData.PullClientMessageInput{}
		err := grpcNetTool.UnmarshalGrpcData(in.Data, input)
		if err != nil {
			return nil, err
		}

		serviceLog.Debug("received client proto Msg: %v", proto.EnvelopeType(input.MsgId))

		serEventModel.EventCallNoReturn(module.NewModuleEventReq(name, in.Data, false, nil))
		out := &methodData.PullClientMessageOutput{Success: true}
		return daprInvoke.MakeOutputContent(in, out)
	}
}
