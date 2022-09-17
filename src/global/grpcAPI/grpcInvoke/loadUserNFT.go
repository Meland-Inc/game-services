package grpcInvoke

import (
	"encoding/json"
	"fmt"

	"github.com/Meland-Inc/game-services/src/common/daprInvoke"
	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	message "github.com/Meland-Inc/game-services/src/global/web3Message"
)

func RPCLoadUserNFTS(userId int64) (*message.GetUserNFTsOutput, error) {
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
	}

	nfts := &message.GetUserNFTsOutput{}
	err = nfts.UnmarshalJSON(outBytes)
	if err != nil {
		serviceLog.Error("UserPlaceablesOutput Unmarshal : err : %+v", err)
		return nil, err
	}
	return nfts, err
}
