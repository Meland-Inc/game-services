package grpcInvoke

import (
	"encoding/json"
	"fmt"

	"github.com/Meland-Inc/game-services/src/common/daprInvoke"
	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	"github.com/Meland-Inc/game-services/src/common/time_helper"
	message "github.com/Meland-Inc/game-services/src/global/web3Message"
)

func RPCLoadNftBuildData(mapId int32) ([]message.BuildData, error) {
	beginMs := time_helper.NowMill()
	defer func() {
		serviceLog.Info("RPCLoadNftBuildData used time [%04d]Ms", time_helper.NowMill()-beginMs)
	}()

	input := message.GetAllBuildDataInput{MapId: int(mapId)}
	inputBytes, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}

	outBytes, err := daprInvoke.InvokeMethod(
		string(message.AppIdWeb3Service),
		string(message.LandServiceActionGetAllBuildData),
		inputBytes,
	)
	if err != nil {
		serviceLog.Error("load all nft build data failed err:%+v", err)
		return nil, err
	}

	output := &message.GetAllBuildDataOutput{}
	err = output.UnmarshalJSON(outBytes)
	if err != nil {
		serviceLog.Error("GetAllLandDataOutput Unmarshal : err : %+v", err)
		return nil, err
	}

	if !output.Success {
		return nil, fmt.Errorf(output.FailedReason)
	}

	return output.AllBuild, nil
}

func RPCBuild(userId int64, nftId string, mapId int32, lands []int32) (*message.BuildData, error) {
	beginMs := time_helper.NowMill()
	defer func() {
		serviceLog.Info("RPCBuild used time [%04d]Ms", time_helper.NowMill()-beginMs)
	}()

	input := message.BuildInput{
		UserId: fmt.Sprint(userId),
		NftId:  nftId,
		MapId:  int(mapId),
	}
	for _, land := range lands {
		input.LandIds = append(input.LandIds, int(land))
	}
	inputBytes, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}

	outBytes, err := daprInvoke.InvokeMethod(
		string(message.AppIdWeb3Service),
		string(message.LandServiceActionBuild),
		inputBytes,
	)
	if err != nil {
		serviceLog.Error("RPCBuild failed err:%+v", err)
		return nil, err
	}

	output := &message.BuildOutput{}
	err = output.UnmarshalJSON(outBytes)
	if err != nil {
		serviceLog.Error("BuildOutput Unmarshal : err : %+v", err)
		return nil, err
	}

	if !output.Success {
		return nil, fmt.Errorf(output.FailedReason)
	}

	return output.BuildData, nil
}

func RPCRecyclingBuild(userId int64, buildId int64, mapId int32) error {
	beginMs := time_helper.NowMill()
	defer func() {
		serviceLog.Info("RPCRecyclingBuild used time [%04d]Ms", time_helper.NowMill()-beginMs)
	}()

	input := message.RecyclingInput{
		UserId:  fmt.Sprint(userId),
		BuildId: int(buildId),
		MapId:   int(mapId),
	}
	inputBytes, err := json.Marshal(input)
	if err != nil {
		return err
	}

	outBytes, err := daprInvoke.InvokeMethod(
		string(message.AppIdWeb3Service),
		string(message.LandServiceActionRecycling),
		inputBytes,
	)
	if err != nil {
		serviceLog.Error("RPCRecyclingBuild failed err:%+v", err)
		return err
	}

	output := &message.RecyclingOutput{}
	err = output.UnmarshalJSON(outBytes)
	if err != nil {
		serviceLog.Error("RPCRecyclingBuild Unmarshal : err : %+v", err)
		return err
	}

	if !output.Success {
		return fmt.Errorf(output.FailedReason)
	}

	return nil
}

func RPCBuildCharged(userId int64, buildId int64, mapId, num, nativeTokenNum int32) error {
	beginMs := time_helper.NowMill()
	defer func() {
		serviceLog.Info("RPCBuildCharged used time [%04d]Ms", time_helper.NowMill()-beginMs)
	}()

	input := message.ChargedInput{
		UserId:            fmt.Sprint(userId),
		BuildId:           int(buildId),
		MapId:             int(mapId),
		Num:               int(num),
		UseNativeTokenNum: int(nativeTokenNum),
	}
	inputBytes, err := json.Marshal(input)
	if err != nil {
		return err
	}

	outBytes, err := daprInvoke.InvokeMethod(
		string(message.AppIdWeb3Service),
		string(message.LandServiceActionCharged),
		inputBytes,
	)
	if err != nil {
		serviceLog.Error("RPCBuildCharged failed err:%+v", err)
		return err
	}

	output := &message.ChargedOutput{}
	err = output.UnmarshalJSON(outBytes)
	if err != nil {
		serviceLog.Error("RPCBuildCharged Unmarshal : err : %+v", err)
		return err
	}

	if !output.Success {
		return fmt.Errorf(output.FailedReason)
	}

	return nil
}

func RPCHarvest(userId int64, buildId int64, mapId int32) error {
	beginMs := time_helper.NowMill()
	defer func() {
		serviceLog.Info("RPCHarvest used time [%04d]Ms", time_helper.NowMill()-beginMs)
	}()

	input := message.HarvestInput{
		UserId:  fmt.Sprint(userId),
		BuildId: int(buildId),
		MapId:   int(mapId),
	}
	inputBytes, err := json.Marshal(input)
	if err != nil {
		return err
	}

	outBytes, err := daprInvoke.InvokeMethod(
		string(message.AppIdWeb3Service),
		string(message.LandServiceActionHarvest),
		inputBytes,
	)
	if err != nil {
		serviceLog.Error("RPCHarvest failed err:%+v", err)
		return err
	}

	output := &message.HarvestOutput{}
	err = output.UnmarshalJSON(outBytes)
	if err != nil {
		serviceLog.Error("RPCHarvest Unmarshal : err : %+v", err)
		return err
	}

	if !output.Success {
		return fmt.Errorf(output.FailedReason)
	}

	return nil
}

func RPCCollection(userId int64, buildId int64, mapId int32) error {
	beginMs := time_helper.NowMill()
	defer func() {
		serviceLog.Info("RPCCollection used time [%04d]Ms", time_helper.NowMill()-beginMs)
	}()

	input := message.CollectionInput{
		UserId:  fmt.Sprint(userId),
		BuildId: int(buildId),
		MapId:   int(mapId),
	}
	inputBytes, err := json.Marshal(input)
	if err != nil {
		return err
	}

	outBytes, err := daprInvoke.InvokeMethod(
		string(message.AppIdWeb3Service),
		string(message.LandServiceActionCollection),
		inputBytes,
	)
	if err != nil {
		serviceLog.Error("RPCCollection failed err:%+v", err)
		return err
	}

	output := &message.CollectionOutput{}
	err = output.UnmarshalJSON(outBytes)
	if err != nil {
		serviceLog.Error("RPCCollection Unmarshal : err : %+v", err)
		return err
	}

	if !output.Success {
		return fmt.Errorf(output.FailedReason)
	}

	return nil
}
