package clientMsgHandle

import (
	"game-message-core/proto"

	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	"github.com/Meland-Inc/game-services/src/services/main/playerModel"
)

func ItemSlotGetHandle(input *proto.PullClientMessageInput) {
	res := &proto.GetItemSlotResponse{}
	respMsg := makeResponseMsg(input.Msg)
	defer func() {
		if respMsg.ErrorMessage != "" {
			respMsg.ErrorCode = 20006 // TODO: USE PROTO ERROR CODE
		}
		respMsg.Payload = &proto.Envelope_GetItemSlotResponse{GetItemSlotResponse: res}
		ResponseClientMessage(input, respMsg)
	}()

	if input.UserId < 1 {
		respMsg.ErrorMessage = "Invalid User ID"
		return
	}

	dataModel, err := playerModel.GetPlayerDataModel()
	if err != nil {
		respMsg.ErrorMessage = err.Error()
		return
	}

	playerSlot, err := dataModel.GetPlayerItemSlots(input.UserId)
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

func ItemSlotUpgradeHandle(input *proto.PullClientMessageInput) {
	res := &proto.UpgradeItemSlotResponse{}
	respMsg := makeResponseMsg(input.Msg)
	defer func() {
		if respMsg.ErrorMessage != "" {
			respMsg.ErrorCode = 20007 // TODO: USE PROTO ERROR CODE
		}
		respMsg.Payload = &proto.Envelope_UpgradeItemSlotResponse{UpgradeItemSlotResponse: res}
		ResponseClientMessage(input, respMsg)
	}()

	if input.UserId < 1 {
		respMsg.ErrorMessage = "Invalid User ID"
		return
	}

	req := input.Msg.GetUpgradeItemSlotRequest()
	if req == nil {
		serviceLog.Error("main service upgrade slot request is nil")
		return
	}

	dataModel, err := playerModel.GetPlayerDataModel()
	if err != nil {
		respMsg.ErrorMessage = err.Error()
		return
	}

	_, err = dataModel.UpgradeItemSlots(input.UserId, req.Position, false)
	if err != nil {
		respMsg.ErrorMessage = err.Error()
		return
	}

	playerSlot, err := dataModel.GetPlayerItemSlots(input.UserId)
	for _, s := range playerSlot.GetSlotList().SlotList {
		res.Slots = append(res.Slots, &proto.ItemSlot{
			Level:    int32(s.Level),
			Position: proto.AvatarPosition(s.Position),
		})
	}
}
