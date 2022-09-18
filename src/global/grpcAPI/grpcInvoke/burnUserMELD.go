package grpcInvoke

import (
	"encoding/json"
	"fmt"

	"github.com/Meland-Inc/game-services/src/common/daprInvoke"
	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	"github.com/Meland-Inc/game-services/src/common/time_helper"
	message "github.com/Meland-Inc/game-services/src/global/web3Message"
)

// Burn player MELD interface
func BurnUserMELD(userId int64, num int) error {
	beginMs := time_helper.NowMill()
	defer func() {
		serviceLog.Info("BurnUser MELD, used time [%04d]Ms", time_helper.NowMill()-beginMs)
	}()

	input := message.UseMELDInput{UserId: fmt.Sprint(userId), Amount: num}
	serviceLog.Info("UseMELDInput = %+v", input)
	inputBytes, err := json.Marshal(input)
	if err != nil {
		return err
	}

	outBytes, err := daprInvoke.InvokeMethod(
		string(message.AppIdWeb3Service),
		string(message.Web3ServiceActionUseMELD),
		inputBytes,
	)

	serviceLog.Info("UseMELDOutput bs = %+v", string(outBytes))

	output := &message.UseMELDOutput{}
	err = output.UnmarshalJSON(outBytes)
	if err != nil {
		serviceLog.Error("GetUseMELDOutput Unmarshal : err : %+v", err)
		return err
	}
	if !output.Success {
		return fmt.Errorf("use user MELD fail err: %s", output.FailedReason)
	}
	return nil
}
