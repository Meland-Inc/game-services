package serviceHandler

import (
	"fmt"
	base_data "game-message-core/grpc/baseData"
	"game-message-core/grpc/methodData"
	"game-message-core/proto"

	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	"github.com/Meland-Inc/game-services/src/global/contract"
	"github.com/Meland-Inc/game-services/src/global/grpcAPI/grpcInvoke"
	"github.com/Meland-Inc/game-services/src/global/module"
	"github.com/Meland-Inc/game-services/src/services/main/home_model"
	land_model "github.com/Meland-Inc/game-services/src/services/main/landModel"
	"github.com/Meland-Inc/game-services/src/services/main/playerModel"
)

func GRPCGetHomeDataHandler(env contract.IModuleEventReq, curMs int64) {
	output := &methodData.MainServiceActionGetHomeDataOutput{Success: true}
	result := &module.ModuleEventResult{}
	defer func() {
		if output.ErrMsg != "" {
			output.Success = false
		}
		serviceLog.Debug("getHomeData output succ = %+v", output.Success)
		result.SetResult(output)
		env.WriteResult(result)
	}()

	input := &methodData.MainServiceActionGetHomeDataInput{}
	err := env.UnmarshalToDaprCallData(input)
	if err != nil {
		output.ErrMsg = err.Error()
		return
	}

	homeModel, err := home_model.GetHomeModel()
	if err != nil {
		output.ErrMsg = err.Error()
		return
	}

	homeData, err := homeModel.GetUserHomeData(input.UserId)
	if err != nil {
		output.ErrMsg = err.Error()
		return
	}
	output.UserId = input.UserId
	output.Data = home_model.ToGrpcHomeData(*homeData)
}

func GRPCGetAllBuildHandler(env contract.IModuleEventReq, curMs int64) {
	output := &methodData.MainServiceActionGetAllBuildOutput{Success: true}
	result := &module.ModuleEventResult{}
	defer func() {
		if output.ErrMsg != "" {
			output.Success = false
		}
		serviceLog.Debug("GetAllBuild output = %+v", output)
		result.SetResult(output)
		env.WriteResult(result)
	}()

	input := &methodData.MainServiceActionGetAllBuildInput{}
	err := env.UnmarshalToDaprCallData(input)
	if err != nil {
		output.ErrMsg = err.Error()
		return
	}

	landModel, err := land_model.GetLandModel()
	if err != nil {
		output.ErrMsg = err.Error()
		return
	}

	mapRecord, err := landModel.GetMapLandRecord(input.MapId)
	if err != nil {
		output.ErrMsg = err.Error()
		return
	}

	nftBuilds := mapRecord.GetAllNftBuild()
	for _, nftBuild := range nftBuilds {
		output.AllBuilds = append(output.AllBuilds, nftBuild.ToGrpcData())
	}
}

func GRPCGetUserDataHandler(env contract.IModuleEventReq, curMs int64) {
	output := &methodData.GetUserDataOutput{Success: true}
	result := &module.ModuleEventResult{}
	defer func() {
		if output.ErrMsg != "" {
			output.Success = false
		}
		result.SetResult(output)
		env.WriteResult(result)
	}()

	input := &methodData.GetUserDataInput{}
	err := env.UnmarshalToDaprCallData(input)
	if err != nil {
		output.ErrMsg = err.Error()
		return
	}

	playerDataModel, err := playerModel.GetPlayerDataModel()
	if err != nil {
		output.ErrMsg = err.Error()
		return
	}

	baseData, sceneData, avatars, profile, err := playerDataModel.PlayerAllData(input.UserId)
	if err != nil {
		output.ErrMsg = err.Error()
		return
	}

	pos := &proto.Vector3{X: sceneData.X, Y: sceneData.Y, Z: sceneData.Z}
	dir := &proto.Vector3{X: sceneData.DirX, Y: sceneData.DirY, Z: sceneData.DirZ}
	pbAvatars := []proto.PlayerAvatar{}
	for _, avatar := range avatars {
		pbAvatars = append(pbAvatars, *avatar.ToNetPlayerAvatar())
	}

	output.BaseData.Set(baseData.ToNetPlayerBaseData())
	output.Profile.Set(profile)
	output.MapId = sceneData.MapId
	output.Pos.Set(pos)
	output.Dir.Set(dir)
	for _, avatar := range pbAvatars {
		grpcAvatar := base_data.GrpcPlayerAvatar{
			Position: avatar.Position,
			ObjectId: avatar.ObjectId,
		}
		grpcAvatar.Attribute = &base_data.GrpcAvatarAttribute{}
		grpcAvatar.Attribute.Set(avatar.Attribute)
		output.Avatars = append(output.Avatars, grpcAvatar)
	}
}

func GRPCMintUserNftHandler(env contract.IModuleEventReq, curMs int64) {
	output := &methodData.MainServiceActionMintNftOutput{Success: true}
	result := &module.ModuleEventResult{}
	defer func() {
		if output.FailedMsg != "" {
			output.Success = false
		}
		result.SetResult(output)
		env.WriteResult(result)
	}()

	input := &methodData.MainServiceActionMintNftInput{}
	err := env.UnmarshalToDaprCallData(input)
	if err != nil {
		output.FailedMsg = err.Error()
		return
	}
	if input.UserId < 1 {
		output.FailedMsg = fmt.Sprintf("invalid user id: %d", input.UserId)
		return
	}

	err = grpcInvoke.Web3MintNFT(input.UserId, input.Item.Cid, input.Item.Num, input.Item.Quality, 0, 0)
	if err != nil {
		output.FailedMsg = err.Error()
	}
}

func GRPCTakeUserNftHandler(env contract.IModuleEventReq, curMs int64) {
	output := &methodData.MainServiceActionTakeNftOutput{Success: true}
	result := &module.ModuleEventResult{}
	defer func() {
		if output.FailedMsg != "" {
			output.Success = false
		}
		result.SetResult(output)
		env.WriteResult(result)
	}()

	input := &methodData.MainServiceActionTakeNftInput{}
	err := env.UnmarshalToDaprCallData(input)
	if err != nil {
		output.FailedMsg = err.Error()
		return
	}
	if input.UserId < 1 {
		output.FailedMsg = fmt.Sprintf("invalid user id: %d", input.UserId)
		return
	}

	playerDataModel, err := playerModel.GetPlayerDataModel()
	if err != nil {
		output.FailedMsg = err.Error()
		return
	}

	playerItem, err := playerDataModel.GetPlayerItems(input.UserId)
	if err != nil {
		output.FailedMsg = err.Error()
		return
	}

	for _, tn := range input.TakeNfts {
		var giveCount = tn.Num
		for _, item := range playerItem.Items {
			if tn.NftId != "" && tn.NftId != item.Id {
				continue
			}
			if tn.ItemCid != 0 && tn.ItemCid != item.Cid {
				continue
			}
			giveCount -= item.Num
			if giveCount <= 0 {
				break
			}
		}
		if giveCount > 0 {
			output.FailedMsg = fmt.Sprintf("not fund NFT %+v", tn)
			return
		}
	}

	for _, takeNft := range input.TakeNfts {
		if takeNft.NftId != "" {
			err = playerDataModel.TakeNftById(input.UserId, takeNft.NftId, takeNft.Num)
		} else {
			err = playerDataModel.TakeNftByItemCid(input.UserId, takeNft.ItemCid, takeNft.Num)
		}
		if err != nil {
			serviceLog.Error(err.Error())
		}
	}
}
