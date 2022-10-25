package serviceMsgHandle

import (
	"game-message-core/proto"

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

func Web3RecyclingHandler(iMsg interface{}) {
	// input, ok := iMsg.(*message.RecyclingEvent)
	// if !ok {
	// 	serviceLog.Error("iMsg to RecyclingEvent failed")
	// 	return
	// }

	// mapRecord, err := getMapLandRecord(int32(input.MapId))
	// if err != nil {
	// 	serviceLog.Error("RecyclingEvent error: %v", err)
	// 	return
	// }

	// err = mapRecord.Recycling(cast.ToInt64(input.UserId), int64(input.BuildId))
	// if err != nil {
	// 	serviceLog.Error("RecyclingEvent error: %v", err)
	// 	return
	// }
}

func Web3BuildUpdateHandler(iMsg interface{}) {
	// input, ok := iMsg.(*message.BuildUpdateEvent)
	// if !ok {
	// 	serviceLog.Error("iMsg to BuildUpdateEvent failed")
	// 	return
	// }

	// mapRecord, err := getMapLandRecord(int32(input.MapId))
	// if err != nil {
	// 	serviceLog.Error("BuildUpdateEvent error: %v", err)
	// 	return
	// }

	// err = mapRecord.UpdateNftBuildWeb3Data(input.BuildData)
	// if err != nil {
	// 	serviceLog.Error("BuildUpdateEvent error: %v", err)
	// 	return
	// }
}
