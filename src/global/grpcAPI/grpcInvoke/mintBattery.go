package grpcInvoke

import (
	"encoding/json"
	"fmt"

	"github.com/Meland-Inc/game-services/src/common/daprInvoke"
	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	"github.com/Meland-Inc/game-services/src/common/time_helper"
	message "github.com/Meland-Inc/game-services/src/global/web3Message"
)

// web3 mint player Battery interface
func MintBattery(userId int64, mintNum, giftNum int32) error {
	beginMs := time_helper.NowMill()
	defer func() {
		serviceLog.Info("MintBattery used time [%04d]Ms", time_helper.NowMill()-beginMs)
	}()

	if userId < 1 || mintNum < 1 {
		return fmt.Errorf("MintBattery invalid data userId[%d], mintNum[%d]", userId, mintNum)
	}

	input := &message.MintBatteryInput{
		UserId:  fmt.Sprint(userId),
		Num:     int(mintNum),
		GiftNum: int(giftNum),
	}

	serviceLog.Info("call web3 MintBatteryInput info = %+v", input)

	inputBytes, err := json.Marshal(input)
	if err != nil {
		serviceLog.Error("MintBatteryInput Marshal err : %+v", err)
		return err
	}

	outBytes, err := daprInvoke.InvokeMethod(
		string(message.AppIdWeb3Service),
		string(message.LandServiceActionMintBattery),
		inputBytes,
	)

	serviceLog.Info("call LandServiceActionMintBattery outPut:[%v], err: %v", string(outBytes), err)

	output := &message.MintBatteryOutput{}
	err = output.UnmarshalJSON(outBytes)
	if err != nil {
		serviceLog.Error("MintBatteryOutput Unmarshal : err : %+v", err)
		return err
	}
	if !output.Success {
		return fmt.Errorf(output.FailedReason)
	}
	return nil
}
