package grpcInvoke

import (
	"encoding/json"
	"fmt"

	"github.com/Meland-Inc/game-services/src/common/daprInvoke"
	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	message "github.com/Meland-Inc/game-services/src/global/web3Message"
)

func RPCCallUseConsumableToWeb3(userId int64, nftId string, x, y int32, args string) error {
	input := message.UseConsumableInput{
		UserId: fmt.Sprint(userId),
		NftId:  nftId,
		Amount: 1,
		LandId: int(message.XyToTileId(x, y)),
		Args:   &args,
	}
	inputBytes, err := json.Marshal(input)
	if err != nil {
		return err
	}

	outBytes, err := daprInvoke.InvokeMethod(
		string(message.AppIdWeb3Service),
		string(message.Web3ServiceActionUseConsumable),
		inputBytes,
	)
	if err != nil {
		serviceLog.Error("grpc useConsumable NFT failed err:%+v", err)
		return err
	}

	output := &message.UseConsumableOutput{}
	err = output.UnmarshalJSON(outBytes)
	if err != nil {
		serviceLog.Error("UseConsumableOutput Unmarshal : err : %+v", err)
		return err
	}
	if !output.Success {
		return fmt.Errorf("grpc useConsumable failed")
	}
	return nil
}
