package serviceHeart

import (
	"time"

	"github.com/Meland-Inc/game-services/src/common/time_helper"
	"github.com/Meland-Inc/game-services/src/global/component"
	"github.com/Meland-Inc/game-services/src/global/serviceCnf"
	"github.com/Meland-Inc/game-services/src/global/serviceRegister"
)

const SERVICE_HEART_CD_MS int64 = 1000 * 3 // 3S

type ServiceHeartModel struct {
	component.ModelBase
	serCnf      *serviceCnf.ServiceConfig
	subModel    ServiceHeartInterface
	nextHeartMs int64
}

func NewServiceHeartModel(cnf *serviceCnf.ServiceConfig) *ServiceHeartModel {
	p := &ServiceHeartModel{serCnf: cnf}
	p.ModelBase.InitBaseModel(p, component.MODEL_NAME_HEART)
	p.subModel = p
	return p
}

func (sh *ServiceHeartModel) Init(
	cnf *serviceCnf.ServiceConfig, sub ServiceHeartInterface,
) {
	sh.serCnf = cnf
	sh.subModel = sub
	sh.ModelBase.InitBaseModel(sh, component.MODEL_NAME_HEART)
}

func (sh *ServiceHeartModel) OnInit(modelMgr *component.ModelManager) error {
	sh.ModelBase.OnInit(modelMgr)
	return nil
}

func (sh *ServiceHeartModel) OnStart() error {
	sh.ModelBase.OnStart()
	sh.updateHeartCD(time_helper.NowUTCMill())
	return nil
}

func (sh *ServiceHeartModel) OnTick(utc time.Time) {
	sh.ModelBase.OnTick(utc)
	if sh.nextHeartMs > utc.UnixMilli() {
		return
	}
	sh.subModel.SendHeart(utc.UnixMilli())
}

func (sh *ServiceHeartModel) EventCall(env *component.ModelEventReq) *component.ModelEventResult {
	return nil
}
func (sh *ServiceHeartModel) EventCallNoReturn(env *component.ModelEventReq)    {}
func (sh *ServiceHeartModel) OnEvent(env *component.ModelEventReq, curMs int64) {}

func (sh *ServiceHeartModel) Secondly(utc time.Time) {}
func (sh *ServiceHeartModel) Minutely(utc time.Time) {}
func (sh *ServiceHeartModel) Hourly(utc time.Time)   {}
func (sh *ServiceHeartModel) Daily(utc time.Time)    {}

func (sh *ServiceHeartModel) Send(online int32, curMs int64) error {
	sh.updateHeartCD(curMs)
	offsetMs, err := serviceRegister.RegisterService(*sh.serCnf, online)
	if err == nil {
		time_helper.SetTimeOffsetMs(offsetMs)
	}
	return err
}

func (sh *ServiceHeartModel) SendHeart(curMs int64) error {
	return sh.Send(0, curMs)
}

func (sh *ServiceHeartModel) updateHeartCD(curMs int64) {
	sh.nextHeartMs = curMs + SERVICE_HEART_CD_MS
}
