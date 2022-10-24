package grpcInvoke

import (
	"encoding/json"
	"fmt"

	"github.com/Meland-Inc/game-services/src/common/daprInvoke"
	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	"github.com/Meland-Inc/game-services/src/common/time_helper"
	message "github.com/Meland-Inc/game-services/src/global/web3Message"
)

func RPCLoadLandData(mapId int32) ([]message.LandData, error) {
	beginMs := time_helper.NowMill()
	defer func() {
		serviceLog.Info("RPCLoadLandData used time [%04d]Ms", time_helper.NowMill()-beginMs)
	}()

	input := message.GetAllLandDataInput{MapId: int(mapId)}
	inputBytes, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}

	outBytes, err := daprInvoke.InvokeMethod(
		string(message.AppIdWeb3Service),
		string(message.LandServiceActionGetAllLandData),
		inputBytes,
	)
	if err != nil {
		serviceLog.Error("load all land data failed err:%+v", err)
		return nil, err
	}

	output := &message.GetAllLandDataOutput{}
	err = output.UnmarshalJSON(outBytes)
	if err != nil {
		serviceLog.Error("GetAllLandDataOutput Unmarshal : err : %+v", err)
		return nil, err
	}

	if !output.Success {
		return nil, fmt.Errorf(output.FailedReason)
	}

	return output.AllLandData, nil
}

func RPCOccupyLand(userId int64, mapId, landId int32) error {
	beginMs := time_helper.NowMill()
	defer func() {
		serviceLog.Info("RPCOccupyLand used time [%04d]Ms", time_helper.NowMill()-beginMs)
	}()

	input := message.OccupyLandInput{
		UserId: fmt.Sprint(userId),
		LandId: int(landId),
		MapId:  int(mapId),
	}
	inputBytes, err := json.Marshal(input)
	if err != nil {
		return err
	}

	outBytes, err := daprInvoke.InvokeMethod(
		string(message.AppIdWeb3Service),
		string(message.LandServiceActionOccupyLand),
		inputBytes,
	)
	if err != nil {
		serviceLog.Error("RPCOccupyLand failed err:%+v", err)
		return err
	}

	output := &message.OccupyLandOutput{}
	err = output.UnmarshalJSON(outBytes)
	if err != nil {
		serviceLog.Error("RPCOccupyLand Unmarshal : err : %+v", err)
		return err
	}

	if !output.Success {
		return fmt.Errorf(output.FailedReason)
	}

	return nil
}
