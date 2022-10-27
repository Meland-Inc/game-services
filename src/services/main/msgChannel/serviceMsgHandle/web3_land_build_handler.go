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

	mapRecord, err := getMapLandRecord(int32(input.MapId))
	if err != nil {
		serviceLog.Error("MultiLandDataUpdateEvent error: %v", err)
		return
	}

	upLands := make([]*proto.LandData, 0, len(input.Lands))
	for _, l := range input.Lands {
		upLands = append(upLands, message.ToProtoLandData(l))
	}
	mapRecord.MultiUpdateLandData(upLands)
}

func Web3MultiRecyclingHandler(iMsg interface{}) {
	input, ok := iMsg.(*message.MultiRecyclingEvent)
	if !ok {
		serviceLog.Error("iMsg to MultiRecyclingEvent failed")
		return
	}

	mapRecord, err := getMapLandRecord(int32(input.MapId))
	if err != nil {
		serviceLog.Error("MultiRecyclingEvent error: %v", err)
		return
	}

	for _, buildId := range input.BuildIds {
		err = mapRecord.OnReceiveRecyclingEvent(int64(buildId))
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

	mapRecord, err := getMapLandRecord(int32(input.MapId))
	if err != nil {
		serviceLog.Error("MultiBuildUpdateEvent error: %v", err)
		return
	}

	for _, build := range input.BuildDatas {
		err = mapRecord.UpdateNftBuildWeb3Data(build)
		if err != nil {
			serviceLog.Error("MultiBuildUpdateEvent error: %v", err)
		}
	}

}
