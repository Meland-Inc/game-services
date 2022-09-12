package service

import (
	"time"

	"github.com/Meland-Inc/game-services/src/common/daprInvoke"
	"github.com/Meland-Inc/game-services/src/common/serviceLog"
)

func (s *Service) onStop() error {
	if err := s.modelMgr.StopModel(); err != nil {
		serviceLog.Error(
			"agent service [%s] StopModel err: %v", s.serviceCnf.ServerName, err,
		)
	}

	time.Sleep(300 * time.Millisecond)
	daprInvoke.Stop()
	return nil
}

func (s *Service) onExit() error {
	s.modelMgr.ExitModel()
	s.modelMgr = nil
	close(s.stopChan)
	close(s.osSignal)
	return nil
}
