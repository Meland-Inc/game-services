package grpcInvoke

import (
	"encoding/json"
	"fmt"

	"github.com/Meland-Inc/game-services/src/common/daprInvoke"
	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	"github.com/Meland-Inc/game-services/src/common/time_helper"
	message "github.com/Meland-Inc/game-services/src/global/web3Message"
)

// web3 mint player NFT interface
func Web3MintNFT(userId int64, itemCid, num, quality, x, y int32) error {
	beginMs := time_helper.NowMill()
	defer func() {
		serviceLog.Info("MintNFT used time [%04d]Ms", time_helper.NowMill()-beginMs)
	}()

	if itemCid == 0 || num < 1 {
		return fmt.Errorf("invalid data cid[%d], num[%d]", itemCid, num)
	}

	qualityStr := fmt.Sprint(quality)
	input := &message.MintNFTWithItemIdInput{
		UserId:     fmt.Sprint(userId),
		ItemId:     fmt.Sprint(itemCid),
		QualityVal: &qualityStr,
		Amount:     int(num),
		LandId:     int(message.XyToTileId(x, y)),
	}

	serviceLog.Info("meland dapr MintNFTWithItemIdInput info = %+v", input)

	inputBytes, err := json.Marshal(input)
	if err != nil {
		serviceLog.Error("MintNFTWithItemIdInput Marshal err : %+v", err)
		return err
	}

	outBytes, err := daprInvoke.InvokeMethod(
		string(message.AppIdWeb3Service),
		string(message.Web3ServiceActionMintNFTWithItemId),
		inputBytes,
	)
	serviceLog.Info("请求 Web3ServiceActionMintNFTWithItemId outPut:[%v], err: %v", string(outBytes), err)
	return err
}
