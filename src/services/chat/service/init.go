package service

import (
	"fmt"
	"game-message-core/proto"
	"os"
	"os/signal"
	"syscall"

	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	"github.com/Meland-Inc/game-services/src/common/time_helper"
	"github.com/Meland-Inc/game-services/src/global/configData"
	"github.com/Meland-Inc/game-services/src/global/daprService"
	"github.com/Meland-Inc/game-services/src/global/gameDB"
	"github.com/Meland-Inc/game-services/src/global/serviceCnf"
	"github.com/Meland-Inc/game-services/src/global/serviceHeart"
	"github.com/Meland-Inc/game-services/src/global/userAgent"
	"github.com/Meland-Inc/game-services/src/services/chat/chatModel"
	chatHandleModule "github.com/Meland-Inc/game-services/src/services/chat/handlerModule"
)

func (s *Service) init() error {
	if err := s.initServiceCnf(); err != nil {
		return err
	}
	serviceLog.Init(s.serviceCnf.AppId, true)
	s.initOsSignal()

	if err := gameDB.Init(); err != nil {
		return err
	}

	if err := configData.Init(); err != nil {
		return err
	}

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

	sc.ServiceType = proto.ServiceType_ServiceTypeChat
	sc.StartMs = time_helper.NowUTCMill()
	sc.AppId = os.Getenv("APP_ID")

	fmt.Println(fmt.Sprintf("serviceCnf: [%+v]", sc))

	if sc.AppId == "" {
		return fmt.Errorf("server app id is empty")
	}
	return nil
}

func (s *Service) initOsSignal() {
	signal.Notify(s.osSignal,
		syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT, syscall.SIGABRT,
	)
}

func (s *Service) initHandlerModel() error {
	model := chatHandleModule.NewHandlerModule()
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

	if err := s.initChatModel(); err != nil {
		return err
	}

	return nil
}

func (s *Service) initHeartModel() error {
	heartModel := serviceHeart.NewServiceHeartModel(s.serviceCnf)
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

func (s *Service) initChatModel() error {
	m := chatModel.NewChatModel()
	err := s.modelMgr.AddModel(m)
	if err != nil {
		serviceLog.Error("init chat model fail, err: %v", err)
	}
	return err
}
