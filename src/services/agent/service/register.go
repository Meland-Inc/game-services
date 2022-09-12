package service

import (
	"encoding/json"
	"game-message-core/grpc"
	"game-message-core/grpc/methodData"
	"time"

	"github.com/Meland-Inc/game-services/src/common/daprInvoke"
	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	"github.com/Meland-Inc/game-services/src/common/time_helper"
	"github.com/Meland-Inc/game-services/src/services/agent/config"
)

func (s *Service) serviceRealInfo() methodData.ServiceDataInput {
	return methodData.ServiceDataInput{
		MsgVersion:  time_helper.NowUTCMill(),
		Id:          config.GetInstance().ServerId,
		Name:        config.GetInstance().ServerName,
		AppId:       config.GetInstance().ServerName,
		ServiceType: config.GetInstance().ServiceType,
		Host:        config.GetInstance().Host,
		Port:        config.GetInstance().Port,
		MapId:       0,
		Online:      1, // TODO ... DADA FROM USER CHANNEL MANAGER
		MaxOnline:   config.GetInstance().MaxOnline,
		CreatedAt:   config.GetInstance().StartMs,
		UpdatedAt:   config.GetInstance().StartMs,
	}
}

func (s *Service) registerService() error {
	time.Sleep(time.Millisecond * 300) // 延时 300Ms 等待dapr init 完成

	input := s.serviceRealInfo()
	inBytes, err := json.Marshal(input)
	if err != nil {
		return err
	}

	_, err = daprInvoke.InvokeMethod(
		string(grpc.AppIdMelandServiceManager),
		string(grpc.ManagerServiceActionRegister),
		inBytes,
	)

	serviceLog.Info("registerService data: %+v, err: %v", input, err)
	return err
}

func (s *Service) unRegisterService() error {
	input := s.serviceRealInfo()
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
