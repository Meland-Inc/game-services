package grpcInvoke

import (
	"encoding/json"
	"fmt"

	"github.com/Meland-Inc/game-services/src/common/daprInvoke"
	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	"github.com/Meland-Inc/game-services/src/common/time_helper"
	message "github.com/Meland-Inc/game-services/src/global/web3Message"
)

func RPCLoadUserNFTS(userId int64) (*message.GetUserNFTsOutput, error) {
	beginMs := time_helper.NowMill()
	defer func() {
		serviceLog.Info("RPCLoadUserNFTS used time [%04d]Ms", time_helper.NowMill()-beginMs)
	}()

	input := message.GetUserNFTsInput{UserId: fmt.Sprint(userId)}
	inputBytes, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}

	outBytes, err := daprInvoke.InvokeMethod(
		string(message.AppIdMelandService),
		string(message.MelandServiceActionGetUserNFTs),
		inputBytes,
	)
	if err != nil {
		serviceLog.Error("load web3 user NFT failed err:%+v", err)
		return nil, err
	}

	nfts := &message.GetUserNFTsOutput{}
	err = nfts.UnmarshalJSON(outBytes)
	if err != nil {
		serviceLog.Error("UserPlaceablesOutput Unmarshal : err : %+v", err)
		return nil, err
	}
	return nfts, err
}
