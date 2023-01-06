package clientHandler

import (
	"game-message-core/grpc/methodData"
	"game-message-core/proto"

	"github.com/Meland-Inc/game-services/src/global/userAgent"
	"github.com/Meland-Inc/game-services/src/services/main/playerModel"
)

func ItemSlotGetHandler(
	agent *userAgent.UserAgentData, input *methodData.PullClientMessageInput, msg *proto.Envelope,
) {
	res := &proto.GetItemSlotResponse{}
	respMsg := userAgent.MakeResponseMsg(msg)
	defer func() {
		if respMsg.ErrorMessage != "" {
			respMsg.ErrorCode = 20006 // TODO: USE PROTO ERROR CODE
		}
		respMsg.Payload = &proto.Envelope_GetItemSlotResponse{GetItemSlotResponse: res}
		userAgent.ResponseClientMessage(agent, input, respMsg)
	}()

	if input.UserId < 1 {
		respMsg.ErrorMessage = "item slot get Invalid User ID"
		return
	}

	playerDataModel, err := playerModel.GetPlayerDataModel()
	if err != nil {
		respMsg.ErrorMessage = err.Error()
		return
	}

	playerSlot, err := playerDataModel.GetPlayerItemSlots(input.UserId)
	if err != nil {
		respMsg.ErrorMessage = err.Error()
		return
	}
	for _, s := range playerSlot.GetSlotList().SlotList {
		res.Slots = append(res.Slots, &proto.ItemSlot{
			Level:    int32(s.Level),
			Position: proto.AvatarPosition(s.Position),
		})
	}
}

func ItemSlotUpgradeHandler(
	agent *userAgent.UserAgentData, input *methodData.PullClientMessageInput, msg *proto.Envelope,
) {
	res := &proto.UpgradePlayerLevelResponse{}
	respMsg := userAgent.MakeResponseMsg(msg)
	defer func() {
		if respMsg.ErrorMessage != "" {
			respMsg.ErrorCode = 20006 // TODO: USE PROTO ERROR CODE
		}
		respMsg.Payload = &proto.Envelope_UpgradePlayerLevelResponse{UpgradePlayerLevelResponse: res}
		userAgent.ResponseClientMessage(agent, input, respMsg)
	}()

	if input.UserId < 1 {
		respMsg.ErrorMessage = "upgrade player level Invalid User ID"
		return
	}

	playerDataModel, err := playerModel.GetPlayerDataModel()
	if err != nil {
		respMsg.ErrorMessage = err.Error()
		return
	}

	lv, exp, err := playerDataModel.UpgradePlayerLevel(input.UserId)
	if err != nil {
		respMsg.ErrorMessage = err.Error()
		return
	}
	res.CurExp = int64(exp)
	res.CurLevel = lv
}
