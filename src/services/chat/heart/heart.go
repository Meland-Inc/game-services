package ChatHeart

import (
	"github.com/Meland-Inc/game-services/src/global/component"
	"github.com/Meland-Inc/game-services/src/global/serviceCnf"
	"github.com/Meland-Inc/game-services/src/global/serviceHeart"
)

type ChatHeart struct {
	serviceHeart.ServiceHeartModel
	serCnf *serviceCnf.ServiceConfig
}

func NewChatHeart(cnf *serviceCnf.ServiceConfig) *ChatHeart {
	return &ChatHeart{serCnf: cnf}
}

func (ah *ChatHeart) OnInit(modelMgr *component.ModelManager) error {
	ah.ServiceHeartModel.OnInit(modelMgr)
	ah.ServiceHeartModel.SubModel = ah
	return nil
}

func (ah *ChatHeart) OnStart() error {
	return ah.ServiceHeartModel.OnStart()
}

func (ah *ChatHeart) OnTick(curMs int64) error {
	return ah.ServiceHeartModel.OnTick(curMs)
}

func (ah *ChatHeart) OnStop() error {
	ah.serCnf = nil
	return ah.ServiceHeartModel.OnStop()
}

func (ah *ChatHeart) OnExit() error {
	return ah.ServiceHeartModel.OnExit()
}

func (ah *ChatHeart) SendHeart(curMs int64) error {
	var online int32 = 0
	return ah.ServiceHeartModel.Send(*ah.serCnf, online, curMs)
}
