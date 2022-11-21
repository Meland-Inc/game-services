package serviceHeart

import (
	"fmt"

	"github.com/Meland-Inc/game-services/src/common/time_helper"
	"github.com/Meland-Inc/game-services/src/global/component"
	"github.com/Meland-Inc/game-services/src/global/serviceCnf"
	"github.com/Meland-Inc/game-services/src/global/serviceRegister"
)

const SERVICE_HEART_CD_MS int64 = 1000 * 3 // 3S

type ServiceHeartModel struct {
	SubModel    ServiceHeartInterface
	modelMgr    *component.ModelManager
	modelName   string
	nextHeartMs int64
}

func (sh *ServiceHeartModel) Name() string {
	return sh.modelName
}

func (sh *ServiceHeartModel) ModelMgr() *component.ModelManager {
	return sh.modelMgr
}

func (sh *ServiceHeartModel) OnInit(modelMgr *component.ModelManager) error {
	if modelMgr == nil {
		return fmt.Errorf("service model manager is nil")
	}
	sh.modelMgr = modelMgr
	sh.modelName = component.MODEL_NAME_HEART
	return nil
}

func (sh *ServiceHeartModel) OnStart() error {
	sh.updateHeartCD(time_helper.NowUTCMill())
	return nil
}

func (sh *ServiceHeartModel) OnTick(curMs int64) error {
	if sh.nextHeartMs > curMs {
		return nil
	}
	return sh.SubModel.SendHeart(curMs)
}

func (sh *ServiceHeartModel) OnStop() error {
	sh.SubModel = nil
	sh.modelMgr = nil
	return nil
}

func (sh *ServiceHeartModel) OnExit() error {
	return nil
}

func (sh *ServiceHeartModel) Send(cnf serviceCnf.ServiceConfig, online int32, curMs int64) error {
	sh.updateHeartCD(curMs)
	offsetMs, err := serviceRegister.RegisterService(cnf, online)
	if err == nil {
		time_helper.SetTimeOffsetMs(offsetMs)
	}
	return err
}

func (sh *ServiceHeartModel) SendHeart(curMs int64) error {
	return nil
}

func (sh *ServiceHeartModel) updateHeartCD(curMs int64) {
	sh.nextHeartMs = curMs + SERVICE_HEART_CD_MS
}
