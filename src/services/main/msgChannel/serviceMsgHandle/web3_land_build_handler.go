package serviceMsgHandle

import (
	"game-message-core/proto"

	message "github.com/Meland-Inc/game-services/src/global/web3Message"

	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	message "github.com/Meland-Inc/game-services/src/global/web3Message"
	land_model "github.com/Meland-Inc/game-services/src/services/main/landModel"
)

func getMapLandRecord(mapId int32) (*land_model.MapLandDataRecord, error) {
	dataModel, err := land_model.GetLandModel()
	if err != nil {
		return nil, err
	}
	return dataModel.GetMapLandRecord(mapId)
}

func Web3MultiLandDataUpdateEventHandler(iMsg interface{}) {
	input, ok := iMsg.(*message.MultiLandDataUpdateEvent)
	if !ok {
		serviceLog.Error("iMsg to MultiLandDataUpdateEvent failed")
		return
	}

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
		mapRecord, err := getMapLandRecord(mapId)
		if err != nil {
			serviceLog.Error("MultiLandDataUpdateEvent error: %v", err)
			continue
		}
		mapRecord.MultiUpdateLandData(upLands)
	}
}

func Web3MultiRecyclingHandler(iMsg interface{}) {
	input, ok := iMsg.(*message.MultiRecyclingEvent)
	if !ok {
		serviceLog.Error("iMsg to MultiRecyclingEvent failed")
		return
	}

	for _, info := range input.RecyclingInfos {
		mapRecord, err := getMapLandRecord(int32(info.MapId))
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

func Web3MultiBuildUpdateHandler(iMsg interface{}) {
	input, ok := iMsg.(*message.MultiBuildUpdateEvent)
	if !ok {
		serviceLog.Error("iMsg to MultiBuildUpdateEvent failed")
		return
	}

	for _, build := range input.BuildDatas {
		mapRecord, err := getMapLandRecord(int32(build.MapId))
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
