package daprCalls

import (
	"context"
	"fmt"

	"game-message-core/grpc/methodData"
	"game-message-core/proto"

	"github.com/Meland-Inc/game-services/src/common/daprInvoke"
	"github.com/Meland-Inc/game-services/src/common/serviceLog"

	"github.com/Meland-Inc/game-services/src/global/grpcAPI/grpcNetTool"
	"github.com/Meland-Inc/game-services/src/global/module"
	"github.com/dapr/go-sdk/service/common"
)

func makeClientMsgHandler(name string) (string, func(ctx context.Context, in *common.InvocationEvent) (*common.Content, error)) {
	return name, func(ctx context.Context, in *common.InvocationEvent) (*common.Content, error) {
		input := &methodData.PullClientMessageInput{}
		err := grpcNetTool.UnmarshalGrpcData(in.Data, input)
		if err != nil {
			return nil, err
		}

		serviceLog.Info("main service received clientPbMsg data: %v", proto.EnvelopeType(input.MsgId))

		out := &methodData.PullClientMessageOutput{Success: true}
		err = onReceiveClientMessage(name, proto.EnvelopeType(input.MsgId), in)
		if err != nil {
			out.Success = false
			out.ErrMsg = err.Error()
		}
		return daprInvoke.MakeOutputContent(in, out)
	}
}

func clientMsgCall(modelName string, name string, in *common.InvocationEvent) error {
	model, exist := module.GetModel(modelName)
	if !exist {
		return fmt.Errorf("%s  model not found", modelName)
	}

	model.EventCallNoReturn(module.NewModuleEventReq(name, in.Data, false, nil))
	return nil
}

func onReceiveClientMessage(name string, msgType proto.EnvelopeType, in *common.InvocationEvent) error {
	switch msgType {
	case proto.EnvelopeType_QueryLands,
		proto.EnvelopeType_Build,
		proto.EnvelopeType_Recycling,
		proto.EnvelopeType_MintBattery,
		proto.EnvelopeType_Charged,
		proto.EnvelopeType_Harvest,
		proto.EnvelopeType_Collection,
		proto.EnvelopeType_SelfNftBuilds:
		return clientMsgCall(module.MODULE_NAME_LAND, name, in)

	case proto.EnvelopeType_SigninPlayer,
		proto.EnvelopeType_ItemGet,
		proto.EnvelopeType_ItemUse,
		proto.EnvelopeType_UpdateAvatar,
		proto.EnvelopeType_UnloadAvatar,
		proto.EnvelopeType_GetItemSlot,
		proto.EnvelopeType_UpgradeItemSlot,
		proto.EnvelopeType_UpgradePlayerLevel:
		return clientMsgCall(module.MODULE_NAME_PLAYER_DATA, name, in)

	case proto.EnvelopeType_QueryGranary,
		proto.EnvelopeType_GranaryCollect:
		return clientMsgCall(module.MODULE_NAME_HOME, name, in)

	}
	return nil
}
