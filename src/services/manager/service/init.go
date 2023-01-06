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
	"github.com/Meland-Inc/game-services/src/services/manager/controller"
	managerHandleModule "github.com/Meland-Inc/game-services/src/services/manager/handlerModule"
	"github.com/Meland-Inc/game-services/src/services/manager/httpSer"
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

	if err := s.initHttpService(); err != nil {
		return err
	}

	return nil
}

func (s *Service) initServiceCnf() error {
	sc := serviceCnf.GetInstance()
	s.serviceCnf = sc

	sc.ServiceType = proto.ServiceType_ServiceTypeManager
	sc.StartMs = time_helper.NowUTCMill()
	sc.AppId = os.Getenv("APP_ID")

	fmt.Println(fmt.Sprintf("serviceCnf: [%+v]", *sc))

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
	model := managerHandleModule.NewHandlerModule()
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

func (s *Service) initHttpService() error {
	if err := httpSer.Init(); err != nil {
		return err
	}
	return nil
}

func (s *Service) initServiceModels() error {
	if err := s.initControllerModel(); err != nil {
		return err
	}

	return nil
}

func (s *Service) initControllerModel() error {
	ctrlModel := controller.NewControllerModel()
	err := s.modelMgr.AddModel(ctrlModel)
	if err != nil {
		serviceLog.Error("init controller model fail, err: %v", err)
	}
	return err
}
