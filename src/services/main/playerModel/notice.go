package playerModel

import (
	"game-message-core/proto"

	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	"github.com/Meland-Inc/game-services/src/global/component"
	"github.com/Meland-Inc/game-services/src/global/serviceCnf"
	"github.com/Meland-Inc/game-services/src/global/userAgent"
)

func (p *PlayerModel) SendToPlayer(userId int64, msg *proto.Envelope) {
	iUserAgentModel, exist := component.GetInstance().GetModel(component.MODEL_NAME_USER_AGENT)
	if !exist {
		return
	}
	agentModel := iUserAgentModel.(*userAgent.UserAgentModel)
	agent, exist := agentModel.GetUserAgent(userId)
	if !exist {
		serviceLog.Warning("user [%d] agent data not found", userId)
		return
	}
	agent.SendToPlayer(serviceCnf.GetInstance().ServerName, msg)
}

func (p *PlayerModel) noticePlayerProfileUpdate(userId int64, profiles []*proto.EntityProfileUpdate) {
	msg := &proto.Envelope{
		Type: proto.EnvelopeType_BroadCastEntityProfileUpdate,
		Payload: &proto.Envelope_BroadCastEntityProfileUpdateResponse{
			BroadCastEntityProfileUpdateResponse: &proto.BroadCastEntityProfileUpdateResponse{
				EntityId: &proto.EntityId{Type: proto.EntityType_EntityTypePlayer, Id: userId},
				Profiles: profiles,
			},
		},
	}

	p.SendToPlayer(userId, msg)
}
