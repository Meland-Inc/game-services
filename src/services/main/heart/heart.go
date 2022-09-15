package mainHeart

import (
	"github.com/Meland-Inc/game-services/src/global/component"
	"github.com/Meland-Inc/game-services/src/global/serviceCnf"
	"github.com/Meland-Inc/game-services/src/global/serviceHeart"
)

type MainHeart struct {
	serviceHeart.ServiceHeartModel
	serCnf *serviceCnf.ServiceConfig
}

func NewMainHeart(cnf *serviceCnf.ServiceConfig) *MainHeart {
	return &MainHeart{serCnf: cnf}
}

func (ah *MainHeart) OnInit(modelMgr *component.ModelManager) error {
	ah.ServiceHeartModel.OnInit(modelMgr)
	ah.ServiceHeartModel.SubModel = ah
	return nil
}

func (ah *MainHeart) OnStop() error {
	ah.serCnf = nil
	return ah.ServiceHeartModel.OnStop()
}

func (ah *MainHeart) SendHeart(curMs int64) error {
	var online int32 = 0
	return ah.ServiceHeartModel.Send(*ah.serCnf, online, curMs)
}
