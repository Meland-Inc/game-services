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
	"github.com/Meland-Inc/game-services/src/services/manager/httpSer"

	daprService "github.com/Meland-Inc/game-services/src/services/agent/dapr"
	"github.com/spf13/cast"
)

func (s *Service) init() error {
	if err := s.initServiceCnf(); err != nil {
		return err
	}
	serviceLog.Init(s.serviceCnf.ServerId, true)

	s.initOsSignal()

	if err := s.initHttpService(); err != nil {
		return err
	}

	if err := s.initDapr(); err != nil {
		return err
	}

	return nil
}

func (s *Service) initServiceCnf() error {
	sc := serviceCnf.GetInstance()
	s.serviceCnf = sc

	sc.ServerId = cast.ToInt64(os.Getenv("MELAND_SERVICE_MGR_NODE_ID"))
	sc.ServiceType = proto.ServiceType_ServiceTypeManager
	sc.StartMs = time_helper.NowUTCMill()
	sc.ServerName = os.Getenv("MELAND_SERVICE_MGR_DAPR_APPID")

	fmt.Println(fmt.Sprintf("serviceCnf: [%+v]", *sc))

	if sc.ServerId == 0 {
		return fmt.Errorf("invalid serviceId [%v]", sc.ServerId)
	}
	if sc.ServerName == "" {
		return fmt.Errorf("server app id is empty")
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
	if err := daprService.Init(); err != nil {
		serviceLog.Error("dapr init fail err:%v", err)
		return err
	}
	return nil
}

func (s *Service) initHttpService() error {
	if err := httpSer.Init(); err != nil {
		return err
	}
	return nil
}
