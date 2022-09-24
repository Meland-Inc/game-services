package grpcInvoke

import (
	"game-message-core/grpc"
	"game-message-core/protoTool"

	"game-message-core/proto"

	"github.com/Meland-Inc/game-services/src/common/daprInvoke"
	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	"github.com/Meland-Inc/game-services/src/common/time_helper"
)

func RPCSelectService(
	serviceType proto.ServiceType, mapId int32,
) (*proto.ManagerActionSelectServiceOutput, error) {
	input := &proto.ManagerActionSelectServiceInput{
		MsgVersion:  time_helper.NowUTCMill(),
		ServiceType: serviceType,
		MapId:       mapId,
	}
	inputBytes, err := protoTool.MarshalProto(input)
	if err != nil {
		serviceLog.Error("Marshal ManagerActionSelectServiceInput failed err: %+v", err)
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

	output := &proto.ManagerActionSelectServiceOutput{}
	err = protoTool.UnmarshalProto(outBytes, output)
	if err != nil {
		serviceLog.Error("select service Output Unmarshal : err : %+v", err)
		return nil, err
	}
	return output, err
}
