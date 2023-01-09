package web3Handler

import (
	"game-message-core/proto"

	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	"github.com/Meland-Inc/game-services/src/global/contract"
	"github.com/Meland-Inc/game-services/src/global/serviceCnf"
	message "github.com/Meland-Inc/game-services/src/global/web3Message"
	land_model "github.com/Meland-Inc/game-services/src/services/main/landModel"
	"github.com/Meland-Inc/game-services/src/services/main/playerModel"
	"github.com/spf13/cast"
)

func Web3MultiLandDataUpdateEvent(env contract.IModuleEventReq, curMs int64) {
	input := &message.MultiLandDataUpdateEvent{}
	err := env.UnmarshalToDaprEventData(input)
	if err != nil {
		serviceLog.Error("Web3MultiLandDataUpdateEvent UnmarshalEvent fail err: %v ", err)
		return
	}

	if input.Etag < int(serviceCnf.GetInstance().StartMs/1000) {
		return
	}

	serviceLog.Info("Receive Web3MultiLandDataUpdateEvent: %+v", input)

	landModel, err := land_model.GetLandModel()
	landGroup := make(map[int32][]*proto.LandData)
	for _, land := range input.Lands {
		mapId := int32(land.MapId)
		pbLandData := message.ToProtoLandData(land)
		if pbLandData == nil {
			continue
		}
		if _, exist := landGroup[mapId]; exist {
			landGroup[mapId] = append(landGroup[mapId], pbLandData)
		} else {
			landGroup[mapId] = []*proto.LandData{pbLandData}
		}
	}

	for mapId, upLands := range landGroup {
		mapRecord, err := landModel.GetMapLandRecord(mapId)
		if err != nil {
			serviceLog.Error("MultiLandDataUpdateEvent error: %v", err)
			continue
		}
		mapRecord.MultiUpdateLandData(upLands)
	}
}

func Web3MultiRecyclingEvent(env contract.IModuleEventReq, curMs int64) {
	input := &message.MultiRecyclingEvent{}
	err := env.UnmarshalToDaprEventData(input)
	if err != nil {
		serviceLog.Error("Web3MultiRecyclingEvent UnmarshalEvent fail err: %v ", err)
		return
	}

	if input.Etag < int(serviceCnf.GetInstance().StartMs/1000) {
		return
	}

	serviceLog.Info("Receive Web3MultiRecyclingEvent: %+v", input)
	landModel, err := land_model.GetLandModel()
	if err != nil {
		serviceLog.Error(err.Error())
		return
	}
	for _, info := range input.RecyclingInfos {
		mapRecord, err := landModel.GetMapLandRecord(int32(info.MapId))
		if err != nil {
			serviceLog.Error("MultiRecyclingEvent error: %v", err)
			return
		}

		err = mapRecord.OnReceiveRecyclingEvent(int64(info.BuildId))
		if err != nil {
			serviceLog.Error("MultiRecyclingEvent error: %v", err)
		}
	}
}

func Web3MultiBuildUpdateEvent(env contract.IModuleEventReq, curMs int64) {
	input := &message.MultiBuildUpdateEvent{}
	err := env.UnmarshalToDaprEventData(input)
	if err != nil {
		serviceLog.Error("Web3MultiBuildUpdateEvent UnmarshalEvent fail err: %v ", err)
		return
	}

	if input.Etag < int(serviceCnf.GetInstance().StartMs/1000) {
		return
	}

	serviceLog.Info("Receive Web3MultiBuildUpdateEvent: %+v", input)
	landModel, err := land_model.GetLandModel()
	if err != nil {
		serviceLog.Error(err.Error())
		return
	}
	for _, build := range input.BuildDatas {
		mapRecord, err := landModel.GetMapLandRecord(int32(build.MapId))
		if err != nil {
			serviceLog.Error("MultiBuildUpdateEvent error: %v", err)
			continue
		}

		err = mapRecord.UpdateNftBuildWeb3Data(build)
		if err != nil {
			serviceLog.Error("MultiBuildUpdateEvent error: %v", err)
		}
	}
}

func Web3UpdateUserNftEvent(env contract.IModuleEventReq, curMs int64) {
	input := &message.UpdateUserNFT{}
	err := env.UnmarshalToDaprEventData(input)
	if err != nil {
		serviceLog.Error("Web3UpdateUserNft UnmarshalEvent fail err: %v ", err)
		return
	}

	if input.Etag < int(serviceCnf.GetInstance().StartMs/1000) {
		return
	}

	serviceLog.Info("Receive Web3UpdateUserNft: %+v", input)

	playerDataModel, err := playerModel.GetPlayerDataModel()
	if err != nil {
		serviceLog.Error(err.Error())
		return
	}

	userId := cast.ToInt64(input.UserId)
	if userId < 1 {
		serviceLog.Error("Web3UpdateUserNft invalid nft Data[%v]", input)
		return
	}

	playerDataModel.UpdatePlayerNFTs(userId, []message.NFT{input.Nft})
}

func Web3MultiUpdateUserNftEvent(env contract.IModuleEventReq, curMs int64) {
	input := &message.MultiUpdateUserNFT{}
	err := env.UnmarshalToDaprEventData(input)
	if err != nil {
		serviceLog.Error("Web3MultiUpdateUserNft UnmarshalEvent fail err: %v ", err)
		return
	}

	if input.Etag < int(serviceCnf.GetInstance().StartMs/1000) {
		return
	}

	serviceLog.Info("Receive Web3MultiUpdateUserNft: %+v", input)

	playerDataModel, err := playerModel.GetPlayerDataModel()
	if err != nil {
		serviceLog.Error(err.Error())
		return
	}

	userId := cast.ToInt64(input.UserId)
	if userId < 1 {
		serviceLog.Error("Web3MultiUpdateUserNft invalid nft Data[%v]", input)
		return
	}

	playerDataModel.UpdatePlayerNFTs(userId, input.Nfts)
}
