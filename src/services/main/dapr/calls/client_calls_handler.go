package daprCalls

import (
	"context"

	"game-message-core/grpc/methodData"
	"game-message-core/proto"

	"github.com/Meland-Inc/game-services/src/common/daprInvoke"
	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	"github.com/Meland-Inc/game-services/src/global/component"
	"github.com/Meland-Inc/game-services/src/global/grpcAPI/grpcNetTool"
	land_model "github.com/Meland-Inc/game-services/src/services/main/landModel"
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

		err = onReceiveClientMessage(name, proto.EnvelopeType(input.MsgId), in)

		out := &methodData.PullClientMessageOutput{Success: true}
		if err != nil {
			out.ErrMsg = err.Error()
		}
		return daprInvoke.MakeOutputContent(in, out)
	}
}

func clientMsgCall(name string, modelName string, in *common.InvocationEvent) error {
	model, err := land_model.GetLandModel()
	if err != nil {
		return err
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
		return clientMsgCall(name, component.MODEL_NAME_LAND, in)

	case proto.EnvelopeType_SigninPlayer,
		proto.EnvelopeType_ItemGet,
		proto.EnvelopeType_ItemUse,
		proto.EnvelopeType_UpdateAvatar,
		proto.EnvelopeType_UnloadAvatar,
		proto.EnvelopeType_GetItemSlot,
		proto.EnvelopeType_UpgradeItemSlot,
		proto.EnvelopeType_UpgradePlayerLevel:
		return clientMsgCall(name, component.MODEL_NAME_PLAYER_DATA, in)
	}
	return nil
}
