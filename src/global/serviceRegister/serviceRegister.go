package serviceRegister

import (
	"fmt"
	"game-message-core/grpc"
	"game-message-core/proto"
	"game-message-core/protoTool"

	"github.com/Meland-Inc/game-services/src/common/daprInvoke"
	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	"github.com/Meland-Inc/game-services/src/common/time_helper"
	"github.com/Meland-Inc/game-services/src/global/serviceCnf"
)

func serviceRealInfo(cnf serviceCnf.ServiceConfig, online int32) *proto.ServiceRegisterInput {
	return &proto.ServiceRegisterInput{
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
		CreateAt:    cnf.StartMs,
		UpdateAt:    time_helper.NowUTCMill(),
	}
}

func RegisterService(cnf serviceCnf.ServiceConfig, online int32) error {
	input := serviceRealInfo(cnf, online)

	inputBytes, err := protoTool.MarshalProto(input)
	if err != nil {
		serviceLog.Error("Marshal RegisterService failed err: %+v", err)
		return err
	}

	bs, err := daprInvoke.InvokeMethod(
		string(grpc.AppIdMelandServiceManager),
		string(grpc.ManagerServiceActionRegister),
		inputBytes,
	)
	if err != nil {
		return err
	}
	out := &proto.ServiceRegisterOutput{}
	err = protoTool.UnmarshalProto(bs, out)
	if err != nil {
		return err
	}
	if !out.Success {
		return fmt.Errorf(out.ErrMsg)
	}
	return err
}

func UnRegisterService(cnf serviceCnf.ServiceConfig, online int32) error {
	input := serviceRealInfo(cnf, online)
	inputBytes, err := protoTool.MarshalProto(input)
	if err != nil {
		serviceLog.Error("Marshal UnRegisterService failed err: %+v", err)
		return err
	}

	_, err = daprInvoke.InvokeMethod(
		string(grpc.AppIdMelandServiceManager),
		string(grpc.ManagerServiceActionDestroy),
		inputBytes,
	)

	return err
}
