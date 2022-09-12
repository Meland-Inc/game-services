package serviceRegister

import (
	"encoding/json"
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
		Id:          cnf.ServerId,
		Name:        cnf.ServerName,
		AppId:       cnf.ServerName,
		ServiceType: cnf.ServiceType,
		Host:        cnf.Host,
		Port:        cnf.Port,
		MapId:       cnf.MapId,
		Online:      online,
		MaxOnline:   cnf.MaxOnline,
		CreatedAt:   cnf.StartMs,
		UpdatedAt:   cnf.StartMs,
	}
}

func RegisterService(cnf serviceCnf.ServiceConfig, online int32) error {
	input := serviceRealInfo(cnf, online)
	inBytes, err := json.Marshal(input)
	if err != nil {
		return err
	}

	_, err = daprInvoke.InvokeMethod(
		string(grpc.AppIdMelandServiceManager),
		string(grpc.ManagerServiceActionRegister),
		inBytes,
	)

	return err
}

func UnRegisterService(cnf serviceCnf.ServiceConfig, online int32) error {
	input := serviceRealInfo(cnf, online)
	inBytes, err := json.Marshal(input)
	if err != nil {
		return err
	}

	_, err = daprInvoke.InvokeMethod(
		string(grpc.AppIdMelandServiceManager),
		string(grpc.ManagerServiceActionDestroy),
		inBytes,
	)

	return err
}
