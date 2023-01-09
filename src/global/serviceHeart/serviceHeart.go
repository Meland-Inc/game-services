package serviceHeart

import (
	"time"

	"github.com/Meland-Inc/game-services/src/common/time_helper"
	"github.com/Meland-Inc/game-services/src/global/contract"
	"github.com/Meland-Inc/game-services/src/global/module"
	"github.com/Meland-Inc/game-services/src/global/serviceCnf"
	"github.com/Meland-Inc/game-services/src/global/serviceRegister"
)

const SERVICE_HEART_CD_MS int64 = 1000 * 3 // 3S

type ServiceHeartModel struct {
	module.ModuleBase
	serCnf      *serviceCnf.ServiceConfig
	subModel    contract.IServiceHeartInterface
	nextHeartMs int64
}

func NewServiceHeartModel(cnf *serviceCnf.ServiceConfig) *ServiceHeartModel {
	p := &ServiceHeartModel{serCnf: cnf}
	p.ModuleBase.InitBaseModel(p, module.MODULE_NAME_HEART)
	p.subModel = p
	return p
}

func (sh *ServiceHeartModel) Init(
	cnf *serviceCnf.ServiceConfig, sub contract.IServiceHeartInterface,
) {
	sh.serCnf = cnf
	sh.subModel = sub
	sh.ModuleBase.InitBaseModel(sh, module.MODULE_NAME_HEART)
}

func (sh *ServiceHeartModel) OnInit() error {
	sh.ModuleBase.OnInit()
	return nil
}

func (sh *ServiceHeartModel) OnStart() error {
	sh.ModuleBase.OnStart()
	sh.updateHeartCD(time_helper.NowUTCMill())
	return nil
}

func (sh *ServiceHeartModel) OnTick(utc time.Time) {
	sh.ModuleBase.OnTick(utc)
	if sh.nextHeartMs > utc.UnixMilli() {
		return
	}
	sh.subModel.SendHeart(utc.UnixMilli())
}

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
