package service

import (
	"github.com/Meland-Inc/game-services/src/services/account/msgChannel"
)

func (s *Service) onStart() error {
	if err := s.modelMgr.StartModel(); err != nil {
		return err
	}

	msgChannel.InitAndRun()
	return nil
}
