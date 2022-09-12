package service

import (
	"fmt"
	"game-message-core/proto"
	"os"
	"os/signal"
	"syscall"

	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	"github.com/Meland-Inc/game-services/src/common/time_helper"
	"github.com/Meland-Inc/game-services/src/global/serviceCnf"

	demoService "github.com/Meland-Inc/game-services/src/services/demo/dapr"
	"github.com/spf13/cast"
)

func (s *Service) init() error {
	if err := s.initServiceCnf(); err != nil {
		return err
	}
	serviceLog.Init(s.serviceCnf.ServerId, true)
	s.initOsSignal()
	if err := s.initDapr(); err != nil {
		return err
	}

	return nil
}

func (s *Service) initServiceCnf() error {
	sc := serviceCnf.GetInstance()
	s.serviceCnf = sc

	sc.ServerId = cast.ToInt64(os.Getenv("MELAND_SERVICE_DEMO_NODE_ID"))
	sc.ServiceType = proto.ServiceType_ServiceTypeAgent
	sc.StartMs = time_helper.NowUTCMill()
	sc.ServerName = os.Getenv("MELAND_SERVICE_AGENT_DAPR_APPID")
	sc.Host = os.Getenv("MELAND_SERVICE_DEMO_SOCKET_HOST")
	sc.Port = cast.ToInt32(os.Getenv("MELAND_SERVICE_DEMO_SOCKET_PORT"))
	sc.MaxOnline = cast.ToInt32(os.Getenv("MELAND_SERVICE_DEMO_ONLINE_LIMIT"))
	if sc.MaxOnline == 0 {
		sc.MaxOnline = 5000
	}

	fmt.Println(fmt.Sprintf(
		"serviceId:[%d], serviceName:[%s], serviceType:[%v], Socket:[%s:%d], maxOnline:[%d]",
		sc.ServerId, sc.ServerName, sc.ServiceType, sc.Host, sc.Port, sc.MaxOnline,
	))

	if sc.ServerId == 0 {
		return fmt.Errorf("invalid serviceId [%v]", sc.ServerId)
	}
	if sc.ServerName == "" {
		return fmt.Errorf("server app id is empty")
	}
	if sc.Port == 0 || sc.Host == "" {
		return fmt.Errorf("invalid socket data, host[%v], port[%v]", sc.Host, sc.Port)
	}
	return nil
}

func (s *Service) initOsSignal() {
	signal.Notify(s.osSignal,
		syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT,
		syscall.SIGABRT, syscall.SIGUSR1, syscall.SIGUSR2,
	)
}

func (s *Service) initDapr() error {
	if err := demoService.Init(); err != nil {
		serviceLog.Error("dapr init fail err:%v", err)
		return err
	}
	return nil
}
