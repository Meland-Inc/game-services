package service

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Meland-Inc/game-services/src/common/serviceLog"

	mgrSerCnf "github.com/Meland-Inc/game-services/src/services/manager/config"
	daprService "github.com/Meland-Inc/game-services/src/services/manager/dapr"
	"github.com/Meland-Inc/game-services/src/services/manager/httpSer"
)

type Service struct {
	osSignal chan os.Signal
}

func NewManagerService() *Service {
	return &Service{
		osSignal: make(chan os.Signal, 1),
	}
}

func (s *Service) OnInit() error {
	if err := mgrSerCnf.GetInstance().Init(); err != nil {
		return err
	}
	serviceLog.Init(mgrSerCnf.GetInstance().ServerId, true)
	signal.Notify(s.osSignal, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	return nil
}

func (s *Service) OnStart() error {
	if err := httpSer.Init(); err != nil {
		return err
	}

	if err := daprService.Init(); err != nil {
		return err
	}

	return nil
}

func (s *Service) OnExit() {
	close(s.osSignal)
	daprService.Stop()
}

func (s *Service) Run() {
	errChan := make(chan error)

	go func() {
		errChan <- httpSer.Run()
	}()

	go func() {
		errChan <- daprService.Run()
	}()

	go func() {
		t := time.NewTicker(1 * time.Second)

		for {
			select {
			case <-t.C:
				// num = s.Tick(timeHelper.NowMill())

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
		serviceLog.Info("service[%v], received signal [%v]", mgrSerCnf.GetInstance().ServerId, si)
		s.OnExit()
	default:
		serviceLog.Info("close gameServer si[%v]", si)
	}
}
