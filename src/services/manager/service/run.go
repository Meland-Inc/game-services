package service

import (
	"fmt"
	"os"
	"syscall"
	"time"

	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	"github.com/Meland-Inc/game-services/src/common/time_helper"
	managerDapr "github.com/Meland-Inc/game-services/src/services/manager/dapr"
	"github.com/Meland-Inc/game-services/src/services/manager/httpSer"
)

func (s *Service) onReceivedOsSignal(si os.Signal) {
	serviceLog.Info("service[%s], received   signal [%v]", s.serviceCnf.AppId, si)
	switch si {
	case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
		serviceLog.Info("service[%s], received signal [%v]", s.serviceCnf.AppId, si)
		s.OnExit()
	default:
		serviceLog.Info("close gameServer si[%v]", si)
		serviceLog.Info("service[%s], received signal [%v]", s.serviceCnf.AppId, si)
	}
}

func (s *Service) run() {
	errChan := make(chan error)
	managerDapr.Run(errChan)

	go func() {
		errChan <- httpSer.Run()
	}()

	go func() {
		t := time.NewTicker(5 * time.Millisecond)

		for {
			select {
			case <-t.C:
				s.onTick(time_helper.NowUTCMill())

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

func (s *Service) onTick(curMs int64) {
	s.modelMgr.TickModel(curMs)
}
