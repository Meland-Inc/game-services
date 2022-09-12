package demoHeart

import (
	"github.com/Meland-Inc/game-services/src/component"
	"github.com/Meland-Inc/game-services/src/global/serviceCnf"
	"github.com/Meland-Inc/game-services/src/global/serviceHeart"
)

type DemoHeart struct {
	serviceHeart.ServiceHeartModel
	serCnf *serviceCnf.ServiceConfig
}

func NewDemoHeart(cnf *serviceCnf.ServiceConfig) *DemoHeart {
	return &DemoHeart{serCnf: cnf}
}

func (ah *DemoHeart) OnInit(modelMgr *component.ModelManager) error {
	ah.ServiceHeartModel.OnInit(modelMgr)
	ah.ServiceHeartModel.SubModel = ah
	return nil
}

func (ah *DemoHeart) OnStart() error {
	return ah.ServiceHeartModel.OnStart()
}

func (ah *DemoHeart) OnTick(curMs int64) error {
	return ah.ServiceHeartModel.OnTick(curMs)
}

func (ah *DemoHeart) OnStop() error {
	ah.serCnf = nil
	return ah.ServiceHeartModel.OnStop()
}

func (ah *DemoHeart) OnExit() error {
	return ah.ServiceHeartModel.OnExit()
}

func (ah *DemoHeart) SendHeart(curMs int64) error {
	var online int32 = 0
	return ah.ServiceHeartModel.Send(*ah.serCnf, online, curMs)
}
