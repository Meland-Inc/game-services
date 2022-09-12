package service

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/Meland-Inc/game-services/src/common/serviceLog"
)

func (s *Service) initOsSignal() {
	signal.Notify(s.osSignal,
		syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT,
		syscall.SIGABRT, syscall.SIGUSR1, syscall.SIGUSR2,
	)
}

func (s *Service) onReceivedOsSignal(si os.Signal, serviceName string) {
	serviceLog.Info("service[%s], received +++++++++ signal [%v]", serviceName, si)
	switch si {
	case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
		serviceLog.Info("service[%s], received signal [%v]", serviceName, si)
		s.OnExit()
	default:
		serviceLog.Info("close gameServer si[%v]", si)
		serviceLog.Info("service[%s], received signal [%v]", serviceName, si)
	}
}
