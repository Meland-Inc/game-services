package itemModel

import (
	"game-message-core/proto"

	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	"github.com/Meland-Inc/game-services/src/global/component"
	"github.com/Meland-Inc/game-services/src/global/serviceCnf"
	"github.com/Meland-Inc/game-services/src/global/userAgent"
)

func (p *ItemModel) SendToPlayer(userId int64, msg *proto.Envelope) {
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

func (p *ItemModel) NoticePlayer(userId int64, nType proto.EnvelopeType, items []*Item) {
	pbItems := []*proto.Item{}
	for _, item := range items {
		if pbIt := item.ToNetItem(); pbIt != nil {
			pbItems = append(pbItems, pbIt)
		}
	}

	msg := &proto.Envelope{Type: nType}
	switch nType {
	case proto.EnvelopeType_BroadCastItemUpdate:
		msg.Payload = &proto.Envelope_BroadCastItemUpdateResponse{
			BroadCastItemUpdateResponse: &proto.BroadCastItemUpdateResponse{Items: pbItems},
		}

	case proto.EnvelopeType_BroadCastItemAdd:
		msg.Payload = &proto.Envelope_BroadCastItemAddResponse{
			BroadCastItemAddResponse: &proto.BroadCastItemAddResponse{Items: pbItems},
		}

	case proto.EnvelopeType_BroadCastItemDel:
		msg.Payload = &proto.Envelope_BroadCastItemDelResponse{
			BroadCastItemDelResponse: &proto.BroadCastItemDelResponse{Items: pbItems},
		}
	default:
		return
	}

	p.SendToPlayer(userId, msg)
}
