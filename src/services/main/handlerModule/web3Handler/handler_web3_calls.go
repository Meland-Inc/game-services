package web3Handler

import (
	"fmt"

	"github.com/Meland-Inc/game-services/src/global/contract"
	"github.com/Meland-Inc/game-services/src/global/module"
	message "github.com/Meland-Inc/game-services/src/global/web3Message"
	"github.com/Meland-Inc/game-services/src/services/main/playerModel"
	"github.com/spf13/cast"
)

func Web3DeductUserExpHandler(env contract.IModuleEventReq, curMs int64) {
	output := &message.DeductUserExpOutput{DeductSuccess: true}
	result := &module.ModuleEventResult{}
	defer func() {
		result.SetResult(output)
		env.WriteResult(result)
	}()

	input := &message.DeductUserExpInput{}
	err := env.UnmarshalToDaprCallData(input)
	if err != nil {
		output.DeductSuccess = false
		result.SetError(err)
		return
	}

	playerDataModel, err := playerModel.GetPlayerDataModel()
	if err != nil {
		output.DeductSuccess = false
		result.SetError(err)
		return
	}

	deductExp := cast.ToInt64(input.DeductExp)
	userId, err := cast.ToInt64E(input.UserId)
	if err != nil || userId < 1 || deductExp < 1 {
		output.DeductSuccess = false
		result.SetError(fmt.Errorf(
			"web3 deduct user exp invalid userId [%s] or expr[%d]", input.UserId, input.DeductExp,
		))
		return
	}

	if err = playerDataModel.DeductExp(userId, int32(deductExp)); err != nil {
		output.DeductSuccess = false
		result.SetError(err)
	}
}

func Web3GetPlayerDataHandler(env contract.IModuleEventReq, curMs int64) {
	output := &message.GetPlayerInfoByUserIdOutput{}
	result := &module.ModuleEventResult{}
	defer func() {
		result.SetResult(output)
		env.WriteResult(result)
	}()

	input := &message.GetPlayerInfoByUserIdInput{}
	err := env.UnmarshalToDaprCallData(input)
	if err != nil {
		result.SetError(err)
		return
	}
	userId, err := cast.ToInt64E(input.UserId)
	if err != nil || userId < 1 {
		result.SetError(fmt.Errorf("web3 getPlayerData invalid userId [%s]", input.UserId))
		return
	}
	playerDataModel, err := playerModel.GetPlayerDataModel()
	if err != nil {
		result.SetError(err)
		return
	}

	sceneData, err := playerDataModel.GetPlayerSceneData(userId)
	if err != nil {
		result.SetError(err)
		return
	}
	baseData, err := playerDataModel.GetPlayerBaseData(userId)
	if err != nil {
		result.SetError(err)
		return
	}

	output.PlayerData = message.PlayerInfo{
		UserId:     input.UserId,
		PlayerName: baseData.Name,
		RoleCId:    int(baseData.RoleId),
		Icon:       baseData.RoleIcon,
		Feature:    baseData.FeatureJson,
		Level:      int(sceneData.Level),
		CurExp:     int(sceneData.Exp),
		CurHp:      int(sceneData.Hp),
	}
}
