package service

import (
	"fmt"
	"game-message-core/proto"
	"os"
	"os/signal"
	"syscall"

	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	"github.com/Meland-Inc/game-services/src/common/time_helper"
	configData "github.com/Meland-Inc/game-services/src/global/configData"
	gameDb "github.com/Meland-Inc/game-services/src/global/gameDB"
	"github.com/Meland-Inc/game-services/src/global/serviceCnf"
	"github.com/Meland-Inc/game-services/src/global/serviceHeart"
	"github.com/Meland-Inc/game-services/src/global/userAgent"
	mainDaprService "github.com/Meland-Inc/game-services/src/services/main/dapr"
	"github.com/Meland-Inc/game-services/src/services/main/home_model"
	land_model "github.com/Meland-Inc/game-services/src/services/main/landModel"
	login_model "github.com/Meland-Inc/game-services/src/services/main/loginModel"
	"github.com/Meland-Inc/game-services/src/services/main/playerModel"
)

func (s *Service) init() error {
	if err := s.initServiceCnf(); err != nil {
		return err
	}
	serviceLog.Init(s.serviceCnf.AppId, true)
	s.initOsSignal()

	if err := gameDb.Init(); err != nil {
		return err
	}

	if err := configData.Init(); err != nil {
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

	sc.StartMs = time_helper.NowUTCMill()
	sc.ServiceType = proto.ServiceType_ServiceTypeMain
	sc.AppId = os.Getenv("APP_ID")
	sc.IsDevelop = os.Getenv("DEVELOP_MODEL") == "true"

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

func (s *Service) initDapr() error {
	if err := mainDaprService.Init(); err != nil {
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

	if err := s.initLoginModel(); err != nil {
		return err
	}

	if err := s.initPlayerDataModel(); err != nil {
		return err
	}

	if err := s.initHomeModel(); err != nil {
		return err
	}

	if err := s.initLandModel(); err != nil {
		return err
	}

	return nil
}

func (s *Service) initHeartModel() error {
	heartModel := serviceHeart.NewServiceHeartModel(s.serviceCnf)
	err := s.modelMgr.AddModel(heartModel)
	if err != nil {
		serviceLog.Error("init service heart model fail, err: %v", err)
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

func (s *Service) initPlayerDataModel() error {
	m := playerModel.NewPlayerModel()
	err := s.modelMgr.AddModel(m)
	if err != nil {
		serviceLog.Error("init player data model fail, err: %v", err)
	}
	return err
}

func (s *Service) initLandModel() error {
	m := land_model.NewLandModel()
	err := s.modelMgr.AddModel(m)
	if err != nil {
		serviceLog.Error("init land  model fail, err: %v", err)
	}
	return err
}

func (s *Service) initLoginModel() error {
	m := login_model.NewLoginModel()
	err := s.modelMgr.AddModel(m)
	if err != nil {
		serviceLog.Error("init login  model fail, err: %v", err)
	}
	return err
}

func (s *Service) initHomeModel() error {
	m := home_model.NewHomeModel()
	err := s.modelMgr.AddModel(m)
	if err != nil {
		serviceLog.Error("init home data model fail, err: %v", err)
	}
	return err
}
