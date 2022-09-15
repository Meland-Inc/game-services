package service

import (
	"fmt"
	"game-message-core/proto"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cast"

	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	"github.com/Meland-Inc/game-services/src/common/time_helper"
	"github.com/Meland-Inc/game-services/src/global/configData"
	"github.com/Meland-Inc/game-services/src/global/serviceCnf"
	accountDb "github.com/Meland-Inc/game-services/src/services/account/accountDB"
	accountDaprService "github.com/Meland-Inc/game-services/src/services/account/dapr"
	accountHeart "github.com/Meland-Inc/game-services/src/services/account/heart"
)

func (s *Service) init() error {
	if err := s.initServiceCnf(); err != nil {
		return err
	}
	serviceLog.Init(s.serviceCnf.ServerId, true)
	s.initOsSignal()

	if err := accountDb.Init(); err != nil {
		return err
	}

	if err := configData.Init(); err != nil {
		return err
	}

	return s.initServiceModels()
}

func (s *Service) initServiceCnf() error {
	sc := serviceCnf.GetInstance()
	s.serviceCnf = sc

	sc.StartMs = time_helper.NowUTCMill()
	sc.ServiceType = proto.ServiceType_ServiceTypeAccount
	sc.ServerId = cast.ToInt64(os.Getenv("MELAND_SERVICE_ACCOUNT_NODE_ID"))
	sc.ServerName = os.Getenv("MELAND_SERVICE_ACCOUNT_DAPR_APPID")
	if sc.MaxOnline == 0 {
		sc.MaxOnline = 5000
	}

	fmt.Println(fmt.Sprintf("serviceCnf: [%+v]", sc))

	if sc.ServerId == 0 {
		return fmt.Errorf("invalid serviceId [%v]", sc.ServerId)
	}
	if sc.ServerName == "" {
		return fmt.Errorf("server app id is empty")
	}
	return nil
}

func (s *Service) initServiceModels() error {
	if err := s.initHeartModel(); err != nil {
		return err
	}

	if err := s.initDapr(); err != nil {
		return err
	}
	return nil
}

func (s *Service) initHeartModel() error {
	heartModel := accountHeart.NewAccountHeart(s.serviceCnf)
	err := s.modelMgr.AddModel(heartModel)
	if err != nil {
		serviceLog.Error("init agent heart model fail, err: %v", err)
	}
	return err
}

func (s *Service) initOsSignal() {
	signal.Notify(s.osSignal,
		syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT,
		syscall.SIGABRT, syscall.SIGUSR1, syscall.SIGUSR2,
	)
}

func (s *Service) initDapr() error {
	if err := accountDaprService.Init(); err != nil {
		serviceLog.Error("dapr init fail err:%v", err)
		return err
	}
	return nil
}
