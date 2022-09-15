package service

import (
	"time"

	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	"github.com/Meland-Inc/game-services/src/global/serviceRegister"
	"github.com/Meland-Inc/game-services/src/services/account/msgChannel"
)

func (s *Service) onStart() error {
	if err := s.modelMgr.StartModel(); err != nil {
		return err
	}

	if err := s.registerService(); err != nil {
		return err
	}

	msgChannel.InitAndRun()

	return nil
}

func (s *Service) registerService() error {
	time.Sleep(time.Millisecond * 300) // 延时 300Ms 等待dapr init 完成
	err := serviceRegister.RegisterService(*s.serviceCnf, 0)
	serviceLog.Info("registerService data: %+v, err: %v", *s.serviceCnf, err)
	return err
}
