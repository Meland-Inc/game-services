package service

import (
	"fmt"
	"os"
	"time"

	"github.com/Meland-Inc/game-services/src/common/daprInvoke"
	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	agentSerCnf "github.com/Meland-Inc/game-services/src/services/agent/config"
	daprService "github.com/Meland-Inc/game-services/src/services/agent/dapr"
	"github.com/Meland-Inc/game-services/src/services/manager/config"
)

type Service struct {
	osSignal chan os.Signal
	stopChan chan chan struct{}
}

func NewAgentService() *Service {
	return &Service{
		osSignal: make(chan os.Signal, 1),
		stopChan: make(chan chan struct{}, 1),
	}
}

func (s *Service) OnInit() error {
	if err := agentSerCnf.GetInstance().Init(); err != nil {
		serviceLog.Error("service init fail err:%v", err)
		return err
	}
	return nil
}

func (s *Service) OnStart() error {
	serviceLog.Init(agentSerCnf.GetInstance().ServerId, true)

	s.initOsSignal()

	if err := daprService.Init(); err != nil {
		serviceLog.Error("dapr init fail err:%v", err)
		return err
	}

	if err := s.registerService(); err != nil {
		serviceLog.Error("service register fail err:%v", err)
		return err
	}

	return nil
}

func (s *Service) OnExit() {
	if err := s.unRegisterService(); err != nil {
		serviceLog.Error(
			"agent service [%s] unRegisterService err: %v",
			config.GetInstance().ServerName, err,
		)
	}
	close(s.stopChan)
	close(s.osSignal)
}

func (s *Service) Run() {
	errChan := make(chan error)
	go func() {
		errChan <- daprInvoke.Start()
	}()

	go func() {
		t := time.NewTicker(1 * time.Second)

		for {

			select {
			case <-t.C:
				// num = s.Tick(timeHelper.NowMill())

			case stopFinished := <-s.stopChan:
				s.OnExit()
				stopFinished <- struct{}{}
				return

			case si := <-s.osSignal:
				s.onReceivedOsSignal(si, agentSerCnf.GetInstance().ServerName)
				errChan <- fmt.Errorf("stop service by os signal")
				return
			}
		}
	}()

	err := <-errChan
	serviceLog.Error(err.Error())
}
