package land_model

import (
	"game-message-core/grpc/methodData"

	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	"github.com/Meland-Inc/game-services/src/global/component"
	"github.com/Meland-Inc/game-services/src/global/grpcAPI/grpcNetTool"
)

func (p *LandModel) GRPCGetAllBuildHandler(env *component.ModelEventReq, curMs int64) {
	inputBs, ok := env.Msg.([]byte)
	serviceLog.Debug("received GetAllBuild : %s, [%v]", inputBs, ok)
	if !ok {
		serviceLog.Error("service GetAllBuild to string failed: %s", inputBs)
		return
	}

	output := &methodData.MainServiceActionGetAllBuildOutput{Success: true}
	result := &component.ModelEventResult{}
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
