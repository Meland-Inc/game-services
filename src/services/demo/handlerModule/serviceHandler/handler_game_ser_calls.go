package serviceHandler

import (
	"game-message-core/grpc/methodData"

	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	"github.com/Meland-Inc/game-services/src/global/contract"
	"github.com/Meland-Inc/game-services/src/global/module"
)

func GRPCGetHomeDataHandler(env contract.IModuleEventReq, curMs int64) {
	output := &methodData.MainServiceActionGetHomeDataOutput{Success: true}
	result := &module.ModuleEventResult{}
	defer func() {
		if output.ErrMsg != "" {
			output.Success = false
		}
		// serviceLog.Debug("getHomeData output = %+v", output)
		serviceLog.Debug("getHomeData output succ = %+v", output.Success)
		result.SetResult(output)
		env.WriteResult(result)
	}()

	input := &methodData.MainServiceActionGetHomeDataInput{}
	err := env.UnmarshalToDaprCallData(input)
	if err != nil {
		output.ErrMsg = err.Error()
		return
	}

	// TODO logic
}
