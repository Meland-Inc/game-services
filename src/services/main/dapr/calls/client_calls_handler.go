package daprCalls

import (
	"context"
	"fmt"

	"game-message-core/grpc/methodData"
	"game-message-core/proto"

	"github.com/Meland-Inc/game-services/src/common/daprInvoke"
	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	"github.com/Meland-Inc/game-services/src/global/component"
	"github.com/Meland-Inc/game-services/src/global/grpcAPI/grpcNetTool"
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
	model, exist := component.GetInstance().GetModel(modelName)
	if !exist {
		return fmt.Errorf("%s  model not found", modelName)
	}

	model.EventCallNoReturn(&component.ModelEventReq{
		EventType: name,
		Msg:       in.Data,
	})
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
		return clientMsgCall(component.MODEL_NAME_LAND, name, in)

	case proto.EnvelopeType_SigninPlayer,
		proto.EnvelopeType_ItemGet,
		proto.EnvelopeType_ItemUse,
		proto.EnvelopeType_UpdateAvatar,
		proto.EnvelopeType_UnloadAvatar,
		proto.EnvelopeType_GetItemSlot,
		proto.EnvelopeType_UpgradeItemSlot,
		proto.EnvelopeType_UpgradePlayerLevel:
		return clientMsgCall(component.MODEL_NAME_PLAYER_DATA, name, in)
	}
	return nil
}
