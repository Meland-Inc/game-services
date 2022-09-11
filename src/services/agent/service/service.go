package service

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Meland-Inc/game-services/src/common/daprInvoke"
	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	agentSerCnf "github.com/Meland-Inc/game-services/src/services/agent/config"
	daprService "github.com/Meland-Inc/game-services/src/services/agent/dapr"
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
		return err
	}

	serviceLog.Init(agentSerCnf.GetInstance().ServerId, true)

	signal.Notify(s.osSignal, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)

	if err := daprService.Init(); err != nil {
		return err
	}

	return nil
}

func (s *Service) OnStart() error {
	return nil
}

func (s *Service) OnExit() {
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
				s.onReceivedOsSignal(si)
				errChan <- fmt.Errorf("stop service by os signal")
				return
			}
		}
	}()

	err := <-errChan
	serviceLog.Error(err.Error())

}

func (s *Service) onReceivedOsSignal(si os.Signal) {
	switch si {
	case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
		serviceLog.Info("service[%v], received signal [%v]", agentSerCnf.GetInstance().ServerId, si)
		s.OnExit()
	default:
		serviceLog.Info("close gameServer si[%v]", si)
	}
}
