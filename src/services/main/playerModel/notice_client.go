package playerModel

import (
	"game-message-core/proto"

	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	dbData "github.com/Meland-Inc/game-services/src/global/gameDB/data"
	"github.com/Meland-Inc/game-services/src/global/serviceCnf"
	"github.com/Meland-Inc/game-services/src/global/userAgent"
)

func (p *PlayerDataModel) SendToPlayer(userId int64, msg *proto.Envelope) {
	agentModel := userAgent.GetUserAgentModel()
	agent, exist := agentModel.GetUserAgent(userId)
	if !exist {
		serviceLog.Warning("user [%d] agent data not found", userId)
		return
	}
	agent.SendToPlayer(serviceCnf.GetInstance().AppId, msg)
}

func (p *PlayerDataModel) noticePlayerProfileUpdate(userId int64, profiles []*proto.EntityProfileUpdate) {
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

func (p *PlayerDataModel) noticePlayerItemMsg(userId int64, nType proto.EnvelopeType, items []*Item) {
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

func (p *PlayerDataModel) noticeUpdatePlayerItemSlot(slot *dbData.ItemSlot) {
	if slot == nil {
		return
	}
	pbSlots := []*proto.ItemSlot{}
	for _, s := range slot.GetSlotList().SlotList {
		pbSlots = append(pbSlots, &proto.ItemSlot{
			Position: proto.AvatarPosition(s.Position),
			Level:    int32(s.Level),
		})
	}
	msg := &proto.Envelope{
		Type: proto.EnvelopeType_BroadCastUpdateItemSlot,
		Payload: &proto.Envelope_BroadCastUpdateItemSlotResponse{
			BroadCastUpdateItemSlotResponse: &proto.BroadCastUpdateItemSlotResponse{
				Slots: pbSlots,
			},
		},
	}
	p.SendToPlayer(slot.UserId, msg)
}
