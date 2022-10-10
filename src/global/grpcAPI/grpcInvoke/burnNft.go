package grpcInvoke

import (
	"encoding/json"
	"fmt"

	"github.com/Meland-Inc/game-services/src/common/daprInvoke"
	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	"github.com/Meland-Inc/game-services/src/common/time_helper"
	message "github.com/Meland-Inc/game-services/src/global/web3Message"
)

// web3 Burn player NFT interface
func BurnNFT(userId int64, nftId string, num int32) error {
	beginMs := time_helper.NowMill()
	defer func() {
		serviceLog.Info("BurnNFT used time [%04d]Ms", time_helper.NowMill()-beginMs)
	}()

	if num < 1 || userId == 0 || nftId == "" {
		return fmt.Errorf("BurnNFT invalid data userId[%d] id[%s], num[%d]", userId, nftId, num)
	}

	input := message.BurnNFTInput{
		UserId: fmt.Sprint(userId),
		NftId:  nftId,
		Amount: int(num),
	}

	serviceLog.Info(" dapr BurnNFT info = %+v", input)
	inputBytes, err := json.Marshal(input)
	if err != nil {
		return err
	}

	outBytes, err := daprInvoke.InvokeMethod(
		string(message.AppIdWeb3Service),
		string(message.Web3ServiceActionBurnNFT),
		inputBytes,
	)
	serviceLog.Info("BurnNFT outPut:[%v], err: %v", string(outBytes), err)
	return err
}
