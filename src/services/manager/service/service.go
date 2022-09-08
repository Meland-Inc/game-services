package service

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Meland-Inc/game-services/src/common/daprInvoke"
	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	"github.com/Meland-Inc/game-services/src/global/configData"
	mgrSerCnf "github.com/Meland-Inc/game-services/src/services/manager/config"
	daprService "github.com/Meland-Inc/game-services/src/services/manager/dapr"
)

type Service struct {
	osSignal chan os.Signal
	stopChan chan chan struct{}
}

func NewManagerService() *Service {
	return &Service{
		stopChan: make(chan chan struct{}),
		osSignal: make(chan os.Signal, 1),
	}
}

func (s *Service) OnInit() error {
	fmt.Println("manager service init ------- begin ----------")
	if err := mgrSerCnf.GetInstance().Init(); err != nil {
		return err
	}

	serviceLog.Init(mgrSerCnf.GetInstance().ServerId, true)
	signal.Notify(s.osSignal, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)

	fmt.Println("manager service init ------- end ----------")
	return nil
}

func (s *Service) OnStart() error {
	serviceLog.Info("this is demo ------- OnStart() ----------")
	if err := configData.Init(); err != nil {
		return err
	}

	if err := daprService.Init(); err != nil {
		return err
	}
	serviceLog.Info("manager service start ------- end ----------")
	return nil
}

func (s *Service) OnExit() {
	serviceLog.Info("manager service  exit ------- begin ----------")
	close(s.stopChan)
	close(s.osSignal)
	daprInvoke.Stop()
	serviceLog.Info("manager service  exit ------- end ----------")
}

func (s *Service) onReceivedOsSignal(si os.Signal) {
	switch si {
	case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
		serviceLog.Info("service[%v], received signal [%v]", mgrSerCnf.GetInstance().ServerId, si)
		stopDone := make(chan struct{}, 1)
		s.stopChan <- stopDone
		<-stopDone

	default:
		serviceLog.Info("close gameServer si[%v]", si)
	}
}

func (s *Service) Run() {
	serviceLog.Info("manager service  run ------- begin ----------")

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
				errChan <- fmt.Errorf("stop service by finish")
				return

			case si := <-s.osSignal:
				s.onReceivedOsSignal(si)
			}
		}
	}()

	err := <-errChan
	serviceLog.Error(err.Error())
	serviceLog.Info("manager service  run ------- end ----------")
}
