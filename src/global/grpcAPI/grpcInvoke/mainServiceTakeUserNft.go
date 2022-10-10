package grpcInvoke

import (
	"encoding/json"
	"errors"
	"game-message-core/grpc"
	"game-message-core/grpc/methodData"

	"github.com/Meland-Inc/game-services/src/common/daprInvoke"
	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	"github.com/Meland-Inc/game-services/src/common/time_helper"
	"github.com/Meland-Inc/game-services/src/global/grpcAPI/grpcNetTool"
)

func RPCMainServiceTakeNFT(userId int64, nfts []methodData.TakeNftData) error {
	input := &methodData.MainServiceActionTakeNftInput{
		MsgVersion: time_helper.NowUTCMill(),
		UserId:     userId,
		TakeNfts:   nfts,
	}
	inputBytes, err := json.Marshal(input)
	if err != nil {
		serviceLog.Error("Marshal MainServiceActionTakeNftInput failed err: %+v", err)
		return err
	}

	outBytes, err := daprInvoke.InvokeMethod(
		string(grpc.AppIdMelandServiceMain),
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
