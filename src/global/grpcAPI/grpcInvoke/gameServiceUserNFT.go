package grpcInvoke

import (
	"encoding/json"
	"errors"
	"game-message-core/grpc"
	base_data "game-message-core/grpc/baseData"
	"game-message-core/grpc/methodData"

	"github.com/Meland-Inc/game-services/src/common/daprInvoke"
	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	"github.com/Meland-Inc/game-services/src/global/grpcAPI/grpcNetTool"
)

func RPCMainServiceTakeNFT(userId int64, nfts []methodData.TakeNftData) error {
	input := &methodData.MainServiceActionTakeNftInput{
		UserId:   userId,
		TakeNfts: nfts,
	}
	inputBytes, err := json.Marshal(input)
	if err != nil {
		serviceLog.Error("Marshal MainServiceActionTakeNftInput failed err: %+v", err)
		return err
	}

	outBytes, err := daprInvoke.InvokeMethod(
		string(grpc.GAME_SERVICE_APPID_MAIN),
		string(grpc.MainServiceActionTakeNFT),
		inputBytes,
	)
	serviceLog.Info("call MainService TakeNft outBytes = %+v, err:%+v", string(outBytes), err)
	if err != nil {
		serviceLog.Error("call main service take nft failed err:%+v", err)
		return err
	}

	output := &methodData.MainServiceActionTakeNftOutput{}
	err = grpcNetTool.UnmarshalGrpcData(outBytes, output)
	if err != nil {
		return err
	}
	if !output.Success {
		return errors.New(output.FailedMsg)
	}
	return nil
}

func RPCMainServiceMintNFT(userId int64, item base_data.GrpcItemBaseInfo) error {
	input := &methodData.MainServiceActionMintNftInput{
		UserId: userId,
		Item:   item,
	}
	inputBytes, err := json.Marshal(input)
	if err != nil {
		serviceLog.Error("Marshal MainServiceActionMintNftInput failed err: %+v", err)
		return err
	}

	outBytes, err := daprInvoke.InvokeMethod(
		string(grpc.GAME_SERVICE_APPID_MAIN),
		string(grpc.MainServiceActionMintNFT),
		inputBytes,
	)
	serviceLog.Info("call MainService MintNFT outBytes = %+v, err:%+v", string(outBytes), err)
	if err != nil {
		serviceLog.Error("call MainService MintNFT failed err:%+v", err)
		return err
	}

	output := &methodData.MainServiceActionMintNftOutput{}
	err = grpcNetTool.UnmarshalGrpcData(outBytes, output)
	if err != nil {
		return err
	}
	if !output.Success {
		return errors.New(output.FailedMsg)
	}
	return nil
}
