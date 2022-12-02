package service

import (
	"fmt"
	"os"
	"syscall"
	"time"

	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	"github.com/Meland-Inc/game-services/src/common/time_helper"
	"github.com/Meland-Inc/game-services/src/global/serviceRegister"
	chatDapr "github.com/Meland-Inc/game-services/src/services/chat/dapr"
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
	chatDapr.Run(errChan)
	s.registerService()
	go func() {
		t := time.NewTicker(5 * time.Millisecond)

		for {
			select {
			case <-t.C:
				s.onTick(time_helper.NowUTC())

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

func (s *Service) registerService() {
	offsetMs, err := serviceRegister.RegisterService(*s.serviceCnf, 0)
	serviceLog.Info("registerService ------ end ----------data: %+v, err: %v", *s.serviceCnf, err)
	if err != nil {
		panic(err)
	}
	time_helper.SetTimeOffsetMs(offsetMs)
}

func (s *Service) onTick(utc time.Time) {
	s.modelMgr.TickModel(utc)
}
