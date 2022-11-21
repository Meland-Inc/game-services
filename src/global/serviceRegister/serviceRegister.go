package serviceRegister

import (
	"encoding/json"
	"fmt"
	"game-message-core/grpc"
	base_data "game-message-core/grpc/baseData"
	"game-message-core/grpc/methodData"
	"game-message-core/grpc/pubsubEventData"
	"time"

	"github.com/Meland-Inc/game-services/src/common/daprInvoke"
	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	"github.com/Meland-Inc/game-services/src/common/time_helper"
	"github.com/Meland-Inc/game-services/src/global/serviceCnf"
)

func serviceRealInfo(
	cnf serviceCnf.ServiceConfig, online int32,
) base_data.ServiceData {
	return base_data.ServiceData{
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

func RegisterService(cnf serviceCnf.ServiceConfig, online int32) (offsetMs int64, err error) {
	input := methodData.ServiceRegisterInput{
		Service:    serviceRealInfo(cnf, online),
		RegisterAt: time.Now().UTC().UnixMilli(),
	}
	inputBytes, err := json.Marshal(input)
	if err != nil {
		return 0, err
	}

	bs, err := daprInvoke.InvokeMethod(
		string(grpc.GAME_SERVICE_APPID_MANAGER),
		string(grpc.ManagerServiceActionRegister),
		inputBytes,
	)
	if err != nil {
		return 0, err
	}
	resAt := time.Now().UTC().UnixMilli()

	out := &methodData.ServiceRegisterOutput{}
	err = json.Unmarshal(bs, out)
	if err != nil {
		return 0, err
	}
	if !out.Success {
		return 0, fmt.Errorf("register service failed")
	}
	// serviceLog.Debug("[%v] register mgrAt[%v] selfAt[%v], resAt[%v], offsetMs[%v]",
	// 	cnf.AppId, out.ManagerAt, out.RegisterAt, resAt,
	// 	(out.ManagerAt-out.RegisterAt)-((resAt-input.RegisterAt)/2),
	// )
	return (out.ManagerAt - out.RegisterAt) - ((resAt - input.RegisterAt) / 2), nil
}

func UnRegisterService(cnf serviceCnf.ServiceConfig, online int32) error {
	env := pubsubEventData.ServiceUnregisterEvent{
		MsgVersion: time_helper.NowUTCMill(),
		Service:    serviceRealInfo(cnf, online),
	}

	inputBytes, err := json.Marshal(env)
	if err != nil {
		serviceLog.Error("RPCPubsub Unregister service Marshal failed err: %+v", err)
		return err
	}

	serviceLog.Info("pubsub event service unregister %+v", env.Service)

	return daprInvoke.PubSubEventCall(string(grpc.SubscriptionEventServiceUnregister), string(inputBytes))
}
