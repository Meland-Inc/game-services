package agentHeart

import (
	"github.com/Meland-Inc/game-services/src/global/component"
	"github.com/Meland-Inc/game-services/src/global/serviceCnf"
	"github.com/Meland-Inc/game-services/src/global/serviceHeart"
)

type AgentHeart struct {
	serviceHeart.ServiceHeartModel
	serCnf *serviceCnf.ServiceConfig
}

func NewAgentHeart(cnf *serviceCnf.ServiceConfig) *AgentHeart {
	return &AgentHeart{serCnf: cnf}
}

func (ah *AgentHeart) OnInit(modelMgr *component.ModelManager) error {
	ah.ServiceHeartModel.OnInit(modelMgr)
	ah.ServiceHeartModel.SubModel = ah
	return nil
}

func (ah *AgentHeart) SendHeart(curMs int64) error {
	var online int32 = 0 // TODO... get online from user channel mgr
	return ah.ServiceHeartModel.Send(*ah.serCnf, online, curMs)
}
