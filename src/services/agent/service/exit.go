package service

import (
	"time"

	"github.com/Meland-Inc/game-services/src/common/daprInvoke"
	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	"github.com/Meland-Inc/game-services/src/global/serviceRegister"
)

func (s *Service) onStop() error {
	s.closed = true
	if err := s.tcpServer.Stop(); err != nil {
		serviceLog.Error(
			"agent service [%s] stop tcp server err: %v", s.serviceCnf.AppId, err,
		)
	}

	s.unRegisterService()
	time.Sleep(100 * time.Millisecond)
	daprInvoke.Stop()

	if err := s.modelMgr.StopModel(); err != nil {
		serviceLog.Error(
			"agent service [%s] StopModel err: %v", s.serviceCnf.AppId, err,
		)
	}

	return nil
}

func (s *Service) unRegisterService() {
	err := serviceRegister.UnRegisterService(*s.serviceCnf, 0)
	serviceLog.Info("UnRegisterService data: %+v, err: %v", *s.serviceCnf, err)
}

func (s *Service) onExit() error {
	s.modelMgr.ExitModel()
	s.modelMgr = nil
	close(s.stopChan)
	close(s.osSignal)
	return nil
}
