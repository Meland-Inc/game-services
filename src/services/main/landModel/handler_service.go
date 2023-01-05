package land_model

import (
	"game-message-core/grpc/methodData"

	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	"github.com/Meland-Inc/game-services/src/global/contract"
	"github.com/Meland-Inc/game-services/src/global/grpcAPI/grpcNetTool"
	"github.com/Meland-Inc/game-services/src/global/module"
)

func (p *LandModel) GRPCGetAllBuildHandler(env contract.IModuleEventReq, curMs int64) {
	inputBs, ok := env.GetMsg().([]byte)
	serviceLog.Debug("received GetAllBuild : %s, [%v]", inputBs, ok)
	if !ok {
		serviceLog.Error("service GetAllBuild to string failed: %s", inputBs)
		return
	}

	output := &methodData.MainServiceActionGetAllBuildOutput{Success: true}
	result := &module.ModuleEventResult{}
	defer func() {
		if output.ErrMsg != "" {
			output.Success = false
		}
		serviceLog.Debug("GetAllBuild output = %+v", output)
		result.SetResult(output)
		env.WriteResult(result)
	}()

	input := &methodData.MainServiceActionGetAllBuildInput{}
	err := grpcNetTool.UnmarshalGrpcData(inputBs, input)
	if err != nil {
		output.ErrMsg = err.Error()
		return
	}

	mapRecord, err := p.GetMapLandRecord(input.MapId)
	if err != nil {
		output.ErrMsg = err.Error()
		return
	}

	nftBuilds := mapRecord.GetAllNftBuild()
	for _, nftBuild := range nftBuilds {
		output.AllBuilds = append(output.AllBuilds, nftBuild.ToGrpcData())
	}
}
