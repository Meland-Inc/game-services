package land_model

import (
	"game-message-core/proto"

	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	"github.com/Meland-Inc/game-services/src/global/serviceCnf"
	"github.com/Meland-Inc/game-services/src/global/userAgent"
)

func (p *MapLandDataRecord) SendToPlayer(userId int64, msg *proto.Envelope) {
	agentModel := userAgent.GetUserAgentModel()
	agent, exist := agentModel.GetUserAgent(userId)
	if !exist {
		serviceLog.Warning("user [%d] agent data not found", userId)
		return
	}
	agent.SendToPlayer(serviceCnf.GetInstance().AppId, msg)
}

func (p *MapLandDataRecord) BroadcastLandDataUpdate(upLands []*proto.LandData) {
	if len(upLands) < 1 {
		return
	}
	agentModel := userAgent.GetUserAgentModel()
	if agentModel == nil {
		serviceLog.Error("broadcast land multi update agent model not found")
		return
	}

	onlinePlayers := agentModel.AllOnlineUserIds()
	msg := &proto.Envelope{
		Type: proto.EnvelopeType_BroadCastMultiUpLand,
		Payload: &proto.Envelope_BroadCastMultiUpLandResponse{
			BroadCastMultiUpLandResponse: &proto.BroadCastMultiUpLandResponse{
				Lands: upLands,
			},
		},
	}
	userAgent.MultipleBroadCastToClient(serviceCnf.GetInstance().AppId, onlinePlayers, msg)
}
