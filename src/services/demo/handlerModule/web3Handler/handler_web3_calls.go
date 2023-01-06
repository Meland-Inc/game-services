package web3Handler

import (
	"github.com/Meland-Inc/game-services/src/global/contract"
	"github.com/Meland-Inc/game-services/src/global/module"
	message "github.com/Meland-Inc/game-services/src/global/web3Message"
)

func Web3DeductUserExpHandler(env contract.IModuleEventReq, curMs int64) {
	output := &message.DeductUserExpOutput{DeductSuccess: true}
	result := &module.ModuleEventResult{}
	defer func() {
		result.SetResult(output)
		env.WriteResult(result)
	}()

	input := &message.DeductUserExpInput{}
	err := env.UnmarshalToDaprCallData(input)
	if err != nil {
		output.DeductSuccess = false
		result.SetError(err)
		return
	}
	// TODO logic
}
