package service

func (s *Service) onStart() error {
	if err := s.modelMgr.StartModel(); err != nil {
		return err
	}
	return nil
}
