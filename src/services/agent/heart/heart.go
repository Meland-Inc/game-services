package agentHeart

import (
	"github.com/Meland-Inc/game-services/src/global/component"
	"github.com/Meland-Inc/game-services/src/global/serviceCnf"
	"github.com/Meland-Inc/game-services/src/global/serviceHeart"
	"github.com/Meland-Inc/game-services/src/services/agent/userChannel"
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
	var online int32 = userChannel.GetInstance().OnlineCount()
	return ah.ServiceHeartModel.Send(*ah.serCnf, online, curMs)
}
