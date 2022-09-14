package demoHeart

import (
	"github.com/Meland-Inc/game-services/src/global/component"
	"github.com/Meland-Inc/game-services/src/global/serviceCnf"
	"github.com/Meland-Inc/game-services/src/global/serviceHeart"
)

type AccountHeart struct {
	serviceHeart.ServiceHeartModel
	serCnf *serviceCnf.ServiceConfig
}

func NewAccountHeart(cnf *serviceCnf.ServiceConfig) *AccountHeart {
	return &AccountHeart{serCnf: cnf}
}

func (ah *AccountHeart) OnInit(modelMgr *component.ModelManager) error {
	ah.ServiceHeartModel.OnInit(modelMgr)
	ah.ServiceHeartModel.SubModel = ah
	return nil
}

func (ah *AccountHeart) SendHeart(curMs int64) error {
	var online int32 = 0
	return ah.ServiceHeartModel.Send(*ah.serCnf, online, curMs)
}
