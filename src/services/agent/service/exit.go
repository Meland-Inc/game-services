package service

import (
	base_data "game-message-core/grpc/baseData"
	"time"

	"github.com/Meland-Inc/game-services/src/common/daprInvoke"
	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	"github.com/Meland-Inc/game-services/src/common/time_helper"
	"github.com/Meland-Inc/game-services/src/global/grpcAPI/grpcPubsubEvent"
)

func (s *Service) onStop() error {
	s.closed = true
	s.unRegisterService()
	time.Sleep(100 * time.Millisecond)
	daprInvoke.Stop()

	if err := s.tcpServer.Stop(); err != nil {
		serviceLog.Error(
			"agent service [%s] stop tcp server err: %v", s.serviceCnf.AppId, err,
		)
	}

	if err := s.modelMgr.StopModel(); err != nil {
		serviceLog.Error(
			"agent service [%s] StopModel err: %v", s.serviceCnf.AppId, err,
		)
	}

	return nil
}

func (s *Service) unRegisterService() {
	data := base_data.ServiceData{
		AppId:       s.serviceCnf.AppId,
		ServiceType: s.serviceCnf.ServiceType,
		Host:        s.serviceCnf.Host,
		Port:        s.serviceCnf.Port,
		MapId:       s.serviceCnf.MapId,
		MaxOnline:   s.serviceCnf.MaxOnline,
		CreatedAt:   s.serviceCnf.StartMs,
		UpdatedAt:   time_helper.NowUTCMill(),
	}
	grpcPubsubEvent.RPCPubsubEventServiceUnregister(data)
}

func (s *Service) onExit() error {
	s.modelMgr.ExitModel()
	s.modelMgr = nil
	close(s.stopChan)
	close(s.osSignal)
	return nil
}
