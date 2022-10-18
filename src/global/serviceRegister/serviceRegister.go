package serviceRegister

import (
	"encoding/json"
	"fmt"
	"game-message-core/grpc"
	"game-message-core/grpc/methodData"

	"github.com/Meland-Inc/game-services/src/common/daprInvoke"
	"github.com/Meland-Inc/game-services/src/common/time_helper"
	"github.com/Meland-Inc/game-services/src/global/serviceCnf"
)

func serviceRealInfo(
	cnf serviceCnf.ServiceConfig, online int32,
) methodData.ServiceDataInput {
	return methodData.ServiceDataInput{
		MsgVersion:  time_helper.NowUTCMill(),
		AppId:       cnf.AppId,
		ServiceType: cnf.ServiceType,
		Host:        cnf.Host,
		Port:        cnf.Port,
		MapId:       cnf.MapId,
		Online:      online,
		MaxOnline:   cnf.MaxOnline,
		CreatedAt:   cnf.StartMs,
		UpdatedAt:   time_helper.NowUTCMill(),
	}
}

func RegisterService(cnf serviceCnf.ServiceConfig, online int32) error {
	input := serviceRealInfo(cnf, online)
	inputBytes, err := json.Marshal(input)
	if err != nil {
		return err
	}

	bs, err := daprInvoke.InvokeMethod(
		string(grpc.GAME_SERVICE_APPID_MANAGER),
		string(grpc.ManagerServiceActionRegister),
		inputBytes,
	)
	if err != nil {
		return err
	}
	out := &methodData.ServiceDataOutput{}
	err = json.Unmarshal(bs, out)
	if err != nil {
		return err
	}
	if !out.Success {
		return fmt.Errorf("register service failed")
	}
	return err
}

func UnRegisterService(cnf serviceCnf.ServiceConfig, online int32) error {
	input := serviceRealInfo(cnf, online)
	inBytes, err := json.Marshal(input)
	if err != nil {
		return err
	}

	_, err = daprInvoke.InvokeMethod(
		string(grpc.GAME_SERVICE_APPID_MANAGER),
		string(grpc.ManagerServiceActionDestroy),
		inBytes,
	)

	return err
}
