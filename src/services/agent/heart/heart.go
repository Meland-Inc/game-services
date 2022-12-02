package agentHeart

import (
	"github.com/Meland-Inc/game-services/src/global/serviceCnf"
	"github.com/Meland-Inc/game-services/src/global/serviceHeart"
	"github.com/Meland-Inc/game-services/src/services/agent/userChannel"
)

type AgentHeart struct {
	serviceHeart.ServiceHeartModel
}

func NewAgentHeart(cnf *serviceCnf.ServiceConfig) *AgentHeart {
	p := &AgentHeart{}
	p.ServiceHeartModel.Init(cnf, p)
	return p
}

func (ah *AgentHeart) SendHeart(curMs int64) error {
	var online int32 = userChannel.GetInstance().OnlineCount()
	return ah.ServiceHeartModel.Send(online, curMs)
}
