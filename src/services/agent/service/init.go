package service

import (
	"fmt"
	"game-message-core/proto"
	"os"
	"os/signal"
	"syscall"

	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	"github.com/Meland-Inc/game-services/src/common/time_helper"
	"github.com/Meland-Inc/game-services/src/global/daprService"
	"github.com/Meland-Inc/game-services/src/global/serviceCnf"
	"github.com/Meland-Inc/game-services/src/global/userAgent"
	agentHandleModule "github.com/Meland-Inc/game-services/src/services/agent/handlerModule"
	agentHeart "github.com/Meland-Inc/game-services/src/services/agent/heart"
	"github.com/spf13/cast"
)

func (s *Service) init() error {
	if err := s.initServiceCnf(); err != nil {
		return err
	}
	serviceLog.Init(s.serviceCnf.AppId, true)
	s.initOsSignal()

	if err := s.initHandlerModel(); err != nil {
		return err
	}

	if err := s.initDapr(); err != nil {
		return err
	}

	if err := s.initServiceModels(); err != nil {
		return err
	}

	return nil
}

func (s *Service) initServiceCnf() error {
	sc := serviceCnf.GetInstance()
	s.serviceCnf = sc

	sc.ServiceType = proto.ServiceType_ServiceTypeAgent
	sc.StartMs = time_helper.NowUTCMill()
	sc.AppId = os.Getenv("APP_ID")
	sc.Host = os.Getenv("SOCKET_HOST")
	sc.Port = cast.ToInt32(os.Getenv("SOCKET_PORT"))
	sc.MaxOnline = cast.ToInt32(os.Getenv("ONLINE_LIMIT"))
	if sc.MaxOnline == 0 {
		sc.MaxOnline = 5000
	}

	fmt.Println(fmt.Sprintf("serviceCnf: [%+v]", sc))

	if sc.AppId == "" {
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
		syscall.SIGABRT,
	)
}

func (s *Service) initHandlerModel() error {
	model := agentHandleModule.NewHandlerModule()
	err := s.modelMgr.AddModel(model)
	if err != nil {
		serviceLog.Error("init service handler model fail, err: %v", err)
	}
	return err
}

func (s *Service) initDapr() error {
	if err := daprService.Init(); err != nil {
		serviceLog.Error("dapr init fail err:%v", err)
		return err
	}
	return nil
}

func (s *Service) initServiceModels() error {
	if err := s.initHeartModel(); err != nil {
		return err
	}
	if err := s.initUserAgentModel(); err != nil {
		return err
	}

	return nil
}

func (s *Service) initHeartModel() error {
	heartModel := agentHeart.NewAgentHeart(s.serviceCnf)
	err := s.modelMgr.AddModel(heartModel)
	if err != nil {
		serviceLog.Error("init agent heart model fail, err: %v", err)
	}
	return err
}

func (s *Service) initUserAgentModel() error {
	m := userAgent.NewUserAgentModel()
	err := s.modelMgr.AddModel(m)
	if err != nil {
		serviceLog.Error("init user agent model fail, err: %v", err)
	}
	return err
}
