package taskHeart

import (
	"github.com/Meland-Inc/game-services/src/global/component"
	"github.com/Meland-Inc/game-services/src/global/serviceCnf"
	"github.com/Meland-Inc/game-services/src/global/serviceHeart"
)

type TaskHeart struct {
	serviceHeart.ServiceHeartModel
	serCnf *serviceCnf.ServiceConfig
}

func NewTaskHeart(cnf *serviceCnf.ServiceConfig) *TaskHeart {
	return &TaskHeart{serCnf: cnf}
}

func (ah *TaskHeart) OnInit(modelMgr *component.ModelManager) error {
	ah.ServiceHeartModel.OnInit(modelMgr)
	ah.ServiceHeartModel.SubModel = ah
	return nil
}

func (ah *TaskHeart) OnStart() error {
	return ah.ServiceHeartModel.OnStart()
}

func (ah *TaskHeart) OnTick(curMs int64) error {
	return ah.ServiceHeartModel.OnTick(curMs)
}

func (ah *TaskHeart) OnStop() error {
	ah.serCnf = nil
	return ah.ServiceHeartModel.OnStop()
}

func (ah *TaskHeart) OnExit() error {
	return ah.ServiceHeartModel.OnExit()
}

func (ah *TaskHeart) SendHeart(curMs int64) error {
	var online int32 = 0
	return ah.ServiceHeartModel.Send(*ah.serCnf, online, curMs)
}
