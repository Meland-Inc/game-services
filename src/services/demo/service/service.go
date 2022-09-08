package service

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Meland-Inc/game-services/src/common/daprInvoke"
	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	"github.com/Meland-Inc/game-services/src/global/configModule"
	demoServiceConfig "github.com/Meland-Inc/game-services/src/services/demo/config"
	daprService "github.com/Meland-Inc/game-services/src/services/demo/dapr"
)

type Service struct {
	ServiceId int32
	Name      string

	osSignal chan os.Signal
	stopChan chan chan struct{}
}

func NewDemoService() *Service {
	return &Service{}
}

func (s *Service) OnInit() error {
	s.osSignal = make(chan os.Signal, 1)
	signal.Notify(s.osSignal, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)

	fmt.Println("this is demo ------- OnInit() ----------")

	s.Name = "demo-service"
	s.ServiceId = 301
	s.stopChan = make(chan chan struct{})
	serviceLog.Init(int64(s.ServiceId), true)

	if err := demoServiceConfig.GetInstance().Init(); err != nil {
		return err
	}

	if err := configModule.Init(); err != nil {
		return err
	}

	if err := daprService.Init(); err != nil {
		return err
	}

	return nil
}

func (s *Service) OnStart() error {
	fmt.Println("this is demo ------- OnStart() ----------")
	return nil
}

func (s *Service) OnExit() {
	fmt.Println("this is demo ------- OnExit() ----------")
	close(s.stopChan)
	close(s.osSignal)
}

func (s *Service) onReceivedOsSignal(si os.Signal) {
	switch si {
	case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
		serviceLog.Info("service[%v], received signal [%v]", s.Name, si)
		stopDone := make(chan struct{}, 1)
		s.stopChan <- stopDone
		<-stopDone

	default:
		serviceLog.Info("close gameServer si[%v]", si)
	}
}

func (s *Service) Run() {
	fmt.Println("this is demo ------- Run() ----------")

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
			}
		}
	}()

	<-errChan
}
