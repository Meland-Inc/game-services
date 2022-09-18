package grpcInvoke

import (
	"encoding/json"
	"game-message-core/grpc"
	"game-message-core/grpc/methodData"
	"game-message-core/proto"

	"github.com/Meland-Inc/game-services/src/common/daprInvoke"
	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	"github.com/Meland-Inc/game-services/src/common/time_helper"
)

func RPCSelectService(
	serviceType proto.ServiceType, mapId int32,
) (*methodData.ManagerActionSelectServiceOutput, error) {
	input := methodData.ManagerActionSelectServiceInput{
		MsgVersion:  time_helper.NowUTCMill(),
		ServiceType: serviceType,
		MapId:       mapId,
	}
	inputBytes, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}

	outBytes, err := daprInvoke.InvokeMethod(
		string(grpc.AppIdMelandServiceManager),
		string(grpc.ManagerServiceActionSelectService),
		inputBytes,
	)
	if err != nil {
		serviceLog.Error("select service[%v][%d] failed err:%+v", serviceType, mapId, err)
	}

	output := &methodData.ManagerActionSelectServiceOutput{}
	err = json.Unmarshal(outBytes, output)
	if err != nil {
		serviceLog.Error("select service Output Unmarshal : err : %+v", err)
		return nil, err
	}
	return output, err
}
