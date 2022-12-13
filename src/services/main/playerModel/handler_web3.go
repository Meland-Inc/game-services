package playerModel

import (
	"fmt"

	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	"github.com/Meland-Inc/game-services/src/global/component"
	"github.com/Meland-Inc/game-services/src/global/grpcAPI/grpcNetTool"
	"github.com/Meland-Inc/game-services/src/global/serviceCnf"
	message "github.com/Meland-Inc/game-services/src/global/web3Message"
	"github.com/dapr/go-sdk/service/common"
	"github.com/spf13/cast"
)

// ---------------------- calls ----------------------------------
func (p *PlayerDataModel) Web3DeductUserExpHandler(env *component.ModelEventReq, curMs int64) {
	inputBs, ok := env.Msg.([]byte)
	serviceLog.Info("service received web3 DeductUserExp : %s, [%v]", inputBs, ok)
	if !ok {
		serviceLog.Error("web3 data to string failed: %s", inputBs)
		return
	}

	output := &message.DeductUserExpOutput{DeductSuccess: true}
	result := &component.ModelEventResult{}
	defer func() {
		result.SetResult(output)
		env.WriteResult(result)
	}()

	input := &message.DeductUserExpInput{}
	err := grpcNetTool.UnmarshalGrpcData(inputBs, input)
	if err != nil {
		output.DeductSuccess = false
		result.Err = err
		return
	}

	deductExp := cast.ToInt64(input.DeductExp)
	userId, err := cast.ToInt64E(input.UserId)
	if err != nil || userId < 1 || deductExp < 1 {
		output.DeductSuccess = false
		result.Err = fmt.Errorf(
			"web3 deduct user exp invalid userId [%s] or expr[%d]", input.UserId, input.DeductExp,
		)
		return
	}

	if err = p.DeductExp(userId, int32(deductExp)); err != nil {
		output.DeductSuccess = false
		result.Err = err
	}
}

func (p *PlayerDataModel) Web3GetPlayerDataHandler(env *component.ModelEventReq, curMs int64) {
	inputBs, ok := env.Msg.([]byte)
	serviceLog.Info("service received web3 GetPlayerData : %s, [%v]", inputBs, ok)
	if !ok {
		serviceLog.Error("web3 data to string failed: %s", inputBs)
		return
	}

	output := &message.GetPlayerInfoByUserIdOutput{}
	result := &component.ModelEventResult{}
	defer func() {
		result.SetResult(output)
		env.WriteResult(result)
	}()

	input := &message.GetPlayerInfoByUserIdInput{}
	err := grpcNetTool.UnmarshalGrpcData(inputBs, input)
	if err != nil {
		result.Err = err
		return
	}
	userId, err := cast.ToInt64E(input.UserId)
	if err != nil || userId < 1 {
		result.Err = fmt.Errorf("web3 getPlayerData invalid userId [%s] or expr[%d]", input.UserId)
		return
	}
	sceneData, err := p.GetPlayerSceneData(userId)
	if err != nil {
		result.Err = err
		return
	}
	baseData, err := p.GetPlayerBaseData(userId)
	if err != nil {
		result.Err = err
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

// -------------------- pubsub event -----------------------

func (p *PlayerDataModel) Web3UpdateUserNftEvent(env *component.ModelEventReq, curMs int64) {
	msg, ok := env.Msg.(*common.TopicEvent)
	serviceLog.Info("Web3UpdateUserNft : %s, [%v]", msg, ok)
	if !ok {
		serviceLog.Error("Web3UpdateUserNft to TopicEvent failed: %v", msg)
		return
	}

	input := &message.UpdateUserNFT{}
	err := grpcNetTool.UnmarshalGrpcTopicEvent(msg, input)
	if err != nil {
		serviceLog.Error("Web3UpdateUserNft UnmarshalEvent fail err: %v ", err)
		return
	}

	if input.Etag < int(serviceCnf.GetInstance().StartMs/1000) {
		return
	}

	serviceLog.Info("Receive Web3UpdateUserNft: %+v", input)

	userId := cast.ToInt64(input.UserId)
	if userId < 1 {
		serviceLog.Error("Web3UpdateUserNft invalid nft Data[%v]", input)
		return
	}

	p.UpdatePlayerNFTs(userId, []message.NFT{input.Nft})
}

func (p *PlayerDataModel) Web3MultiUpdateUserNftEvent(env *component.ModelEventReq, curMs int64) {
	msg, ok := env.Msg.(*common.TopicEvent)
	serviceLog.Info("Web3MultiUpdateUserNft : %s, [%v]", msg, ok)
	if !ok {
		serviceLog.Error("Web3MultiUpdateUserNft to TopicEvent failed: %v", msg)
		return
	}

	input := &message.MultiUpdateUserNFT{}
	err := grpcNetTool.UnmarshalGrpcTopicEvent(msg, input)
	if err != nil {
		serviceLog.Error("Web3MultiUpdateUserNft UnmarshalEvent fail err: %v ", err)
		return
	}

	if input.Etag < int(serviceCnf.GetInstance().StartMs/1000) {
		return
	}

	serviceLog.Info("Receive Web3MultiUpdateUserNft: %+v", input)

	userId := cast.ToInt64(input.UserId)
	if userId < 1 {
		serviceLog.Error("Web3MultiUpdateUserNft invalid nft Data[%v]", input)
		return
	}

	p.UpdatePlayerNFTs(userId, input.Nfts)
}
