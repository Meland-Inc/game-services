package grpcInvoke

import (
	"encoding/json"
	"fmt"
	"game-message-core/proto"

	"github.com/Meland-Inc/game-services/src/common/daprInvoke"
	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	"github.com/Meland-Inc/game-services/src/common/time_helper"
	message "github.com/Meland-Inc/game-services/src/global/web3Message"
)

func makeDynamicServiceEnvs(
	subType proto.SceneServiceSubType, serviceOwner int64, mapId, maxOnline int32,
) (string, error) {
	/*
		export SERVICE_SUB_TYPE=world   #(world | home | dungeon)
		export SERVICE_OWNER=0             #( 0 |  home owner id)
		export GAME_MAP_ID=10001
		export ONLINE_LIMIT=3000
	*/
	if maxOnline < 1 {
		return "", fmt.Errorf("invalid maxOnline")
	}
	if mapId < 1 {
		return "", fmt.Errorf("invalid mapId")
	}

	subTypeStr := ""
	switch subType {
	case proto.SceneServiceSubType_World:
		subTypeStr = "world"

	case proto.SceneServiceSubType_Home:
		subTypeStr = "home"
		if serviceOwner < 1 {
			return "", fmt.Errorf("invalid home owner")
		}

	case proto.SceneServiceSubType_Dungeon:
		subTypeStr = "dungeon"
		if serviceOwner < 1 {
			return "", fmt.Errorf("invalid dungeon owner")
		}

	default:
		return "", fmt.Errorf("invalid scene service sub type")
	}

	envs := fmt.Sprintf(
		"SERVICE_SUB_TYPE=%s SERVICE_OWNER=%d GAME_MAP_ID=%d ONLINE_LIMIT=%d",
		subTypeStr, serviceOwner, mapId, maxOnline,
	)
	return envs, nil
}

// Scene Dynamic Service start
func GRPCDynamicStartSceneService(
	subType proto.SceneServiceSubType, serviceOwner int64, mapId, maxOnline int32,
) (serAppId string, err error) {
	beginMs := time_helper.NowMill()
	defer func() {
		serviceLog.Info("start private scene Ser used time [%04d]Ms", time_helper.NowMill()-beginMs)
	}()

	envs, err := makeDynamicServiceEnvs(subType, serviceOwner, mapId, maxOnline)
	if err != nil {
		return "", err
	}

	input := message.StartServerInput{
		Args: "",
		Envs: envs,
	}

	serviceLog.Info("dapr Dynamic StartServerInput = %+v", input)
	inputBytes, err := json.Marshal(input)
	if err != nil {
		return "", err
	}

	outBytes, err := daprInvoke.InvokeMethod(
		string(message.AppIdSceneDynamicService),
		string(message.SceneDynamicServiceActionStartServer),
		inputBytes,
	)
	serviceLog.Info("Dynamic StartServer output:[%v], err: %v", string(outBytes), err)
	if err != nil {
		return "", err
	}

	output := &message.StartServerOutput{}
	err = output.UnmarshalJSON(outBytes)
	if err != nil {
		serviceLog.Error("StartServerOutput Unmarshal : err : %+v", err)
		return "", err
	}

	return output.ServerAppId, nil
}
