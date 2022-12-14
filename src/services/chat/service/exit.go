package service

import (
	"time"

	"github.com/Meland-Inc/game-services/src/common/daprInvoke"
	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	"github.com/Meland-Inc/game-services/src/global/serviceRegister"
)

func (s *Service) onStop() error {
	if err := s.unRegisterService(); err != nil {
		serviceLog.Error(
			"service [%s] unRegisterService err: %v", s.serviceCnf.AppId, err,
		)
	}

	if err := s.modelMgr.StopModel(); err != nil {
		serviceLog.Error(
			"service [%s] StopModel err: %v", s.serviceCnf.AppId, err,
		)
	}

	time.Sleep(300 * time.Millisecond)
	daprInvoke.Stop()
	return nil
}

func (s *Service) unRegisterService() error {
	err := serviceRegister.UnRegisterService(*s.serviceCnf, 0)
	serviceLog.Info("UnRegisterService data: %+v, err: %v", *s.serviceCnf, err)
	return err
}

func (s *Service) onExit() error {
	s.modelMgr.ExitModel()
	s.modelMgr = nil
	close(s.stopChan)
	close(s.osSignal)
	return nil
}
