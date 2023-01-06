package grpcInvoke

import (
	"encoding/json"
	"game-message-core/grpc"
	"game-message-core/grpc/methodData"

	"game-message-core/proto"

	"github.com/Meland-Inc/game-services/src/common/daprInvoke"
	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	"github.com/Meland-Inc/game-services/src/global/grpcAPI/grpcNetTool"
)

func RPCSelectService(
	serviceType proto.ServiceType,
	sceneSubType proto.SceneServiceSubType,
	ownerId int64,
	mapId int32,
) (*methodData.ManagerActionSelectServiceOutput, error) {
	input := &methodData.ManagerActionSelectServiceInput{
		ServiceType:     serviceType,
		SceneSerSubType: sceneSubType,
		OwnerId:         ownerId,
		MapId:           mapId,
	}
	inputBytes, err := json.Marshal(input)
	if err != nil {
		serviceLog.Error("Marshal ManagerActionSelectServiceInput failed err: %+v", err)
		return nil, err
	}

	outBytes, err := daprInvoke.InvokeMethod(
		string(grpc.GAME_SERVICE_APPID_MANAGER),
		string(grpc.ManagerServiceActionSelectService),
		inputBytes,
	)
	serviceLog.Info("Select Service output = %+v, err:%+v", string(outBytes), err)
	if err != nil {
		serviceLog.Error("select service[%v][%v][%v][%v] failed err:%+v", serviceType, sceneSubType, mapId, ownerId, err)
		return nil, err
	}

	output := &methodData.ManagerActionSelectServiceOutput{}
	err = grpcNetTool.UnmarshalGrpcData(outBytes, output)
	if err != nil {
		serviceLog.Error("select service Output Unmarshal : err : %+v", err)
		return nil, err
	}
	return output, err
}

func RPCMultiSelectService(
	serviceType proto.ServiceType,
	sceneSubType proto.SceneServiceSubType,
	ownerId int64,
	mapId int32,
) (*methodData.MultiSelectServiceOutput, error) {
	input := &methodData.MultiSelectServiceInput{
		ServiceType:     serviceType,
		SceneSerSubType: sceneSubType,
		OwnerId:         ownerId,
		MapId:           mapId,
	}
	inputBytes, err := json.Marshal(input)
	if err != nil {
		serviceLog.Error("Marshal MultiSelectServiceInput failed err: %+v", err)
		return nil, err
	}

	outBytes, err := daprInvoke.InvokeMethod(
		string(grpc.GAME_SERVICE_APPID_MANAGER),
		string(grpc.ManagerServiceActionMultiSelectService),
		inputBytes,
	)
	serviceLog.Info("MultiSelectServiceOutput outBytes = %+v, err:%+v", string(outBytes), err)
	if err != nil {
		serviceLog.Error("multi select service[%v][%v][%v][%v] failed err:%+v", serviceType, sceneSubType, mapId, ownerId, err)
		return nil, err
	}

	output := &methodData.MultiSelectServiceOutput{}
	err = grpcNetTool.UnmarshalGrpcData(outBytes, output)
	if err != nil {
		serviceLog.Error("multi select service Output Unmarshal : err : %+v", err)
		return nil, err
	}
	return output, err
}
