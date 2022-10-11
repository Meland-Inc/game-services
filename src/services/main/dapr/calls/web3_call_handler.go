package daprCalls

import (
	"context"
	"fmt"

	"github.com/Meland-Inc/game-services/src/common/daprInvoke"
	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	"github.com/Meland-Inc/game-services/src/global/grpcAPI/grpcNetTool"
	message "github.com/Meland-Inc/game-services/src/global/web3Message"
	"github.com/Meland-Inc/game-services/src/services/main/msgChannel"
	"github.com/Meland-Inc/game-services/src/services/main/playerModel"
	"github.com/dapr/go-sdk/service/common"
	"github.com/spf13/cast"
)

func Web3DeductUserExpHandler(ctx context.Context, in *common.InvocationEvent) (*common.Content, error) {
	resFunc := func(success bool, err error) (*common.Content, error) {
		out := &message.DeductUserExpOutput{}
		out.DeductSuccess = success
		if err != nil {
			out.FailedReason = err.Error()
			serviceLog.Error("web3 deduct user exp err: %v", err)
		}
		content, _ := daprInvoke.MakeOutputContent(in, out)
		return content, err
	}

	serviceLog.Info("web3 deduct user exp received data: %v", string(in.Data))

	input := &message.DeductUserExpInput{}
	err := grpcNetTool.UnmarshalGrpcData(in.Data, input)
	if err != nil {
		return nil, err
	}

	deductExp := cast.ToInt64(input.DeductExp)
	userId, err := cast.ToInt64E(input.UserId)
	if err != nil || userId < 1 || deductExp < 1 {
		return resFunc(
			false,
			fmt.Errorf("web3 deduct user exp invalid userId [%s] or expr[%d]",
				input.UserId, input.DeductExp,
			),
		)
	}

	dataModel, err := playerModel.GetPlayerDataModel()
	if err != nil {
		return resFunc(false, fmt.Errorf("player data model not found"))
	}

	sceneData, err := dataModel.GetPlayerSceneData(userId)
	if sceneData.Exp < int32(deductExp) {
		out, err := resFunc(false, fmt.Errorf("insufficient experience"))
		return out, err
	}

	msgChannel.GetInstance().CallServiceMsg(&msgChannel.ServiceMsgData{
		MsgId:   string(message.GameDataServiceActionDeductUserExp),
		MsgBody: input,
	})

	return resFunc(true, nil)
}

func Web3GetPlayerDataHandler(ctx context.Context, in *common.InvocationEvent) (*common.Content, error) {
	resFunc := func(success bool, err error, data message.PlayerInfo) (*common.Content, error) {
		out := &message.GetPlayerInfoByUserIdOutput{
			PlayerData: data,
		}
		if err != nil {
			serviceLog.Error("web3 get user data err: %v", err)
		}
		content, _ := daprInvoke.MakeOutputContent(in, out)
		return content, err
	}

	serviceLog.Info("web3 get user data received data: %v", string(in.Data))

	input := &message.GetPlayerInfoByUserIdInput{}
	err := grpcNetTool.UnmarshalGrpcData(in.Data, input)
	if err != nil {
		return nil, err
	}

	userId, err := cast.ToInt64E(input.UserId)
	if err != nil || userId < 1 {
		return resFunc(
			false,
			fmt.Errorf("web3 get user data invalid userId [%s]", input.UserId),
			message.PlayerInfo{},
		)
	}

	dataModel, err := playerModel.GetPlayerDataModel()
	if err != nil {
		return resFunc(false, fmt.Errorf("player data model not found"), message.PlayerInfo{})
	}

	sceneData, err := dataModel.GetPlayerSceneData(userId)
	if err != nil {
		return resFunc(false, err, message.PlayerInfo{})
	}
	baseData, err := dataModel.GetPlayerBaseData(userId)
	if err != nil {
		return resFunc(false, err, message.PlayerInfo{})
	}

	data := message.PlayerInfo{
		UserId:     input.UserId,
		PlayerName: baseData.Name,
		RoleCId:    int(baseData.RoleId),
		Icon:       baseData.RoleIcon,
		Feature:    baseData.FeatureJson,
		Level:      int(sceneData.Level),
		CurExp:     int(sceneData.Exp),
		CurHp:      int(sceneData.Hp),
	}

	return resFunc(true, nil, data)
}
