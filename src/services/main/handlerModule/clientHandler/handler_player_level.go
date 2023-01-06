package clientHandler

import (
	"game-message-core/grpc/methodData"
	"game-message-core/proto"

	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	"github.com/Meland-Inc/game-services/src/global/userAgent"
	"github.com/Meland-Inc/game-services/src/services/main/playerModel"
)

func PlayerLevelUpgradeHandler(
	agent *userAgent.UserAgentData, input *methodData.PullClientMessageInput, msg *proto.Envelope,
) {
	res := &proto.UpgradeItemSlotResponse{}
	respMsg := userAgent.MakeResponseMsg(msg)
	defer func() {
		if respMsg.ErrorMessage != "" {
			respMsg.ErrorCode = 20007 // TODO: USE PROTO ERROR CODE
		}
		respMsg.Payload = &proto.Envelope_UpgradeItemSlotResponse{UpgradeItemSlotResponse: res}
		userAgent.ResponseClientMessage(agent, input, respMsg)
	}()

	if input.UserId < 1 {
		respMsg.ErrorMessage = "item slot upgrade Invalid User ID"
		return
	}

	playerDataModel, err := playerModel.GetPlayerDataModel()
	if err != nil {
		respMsg.ErrorMessage = err.Error()
		return
	}

	req := msg.GetUpgradeItemSlotRequest()
	if req == nil {
		serviceLog.Error("main service upgrade slot request is nil")
		return
	}
	slotData, err := playerDataModel.UpgradeItemSlots(input.UserId, req.Position, false)
	if err != nil {
		respMsg.ErrorMessage = err.Error()
		return
	}
	for _, s := range slotData.GetSlotList().SlotList {
		res.Slots = append(res.Slots, &proto.ItemSlot{
			Level:    int32(s.Level),
			Position: proto.AvatarPosition(s.Position),
		})
	}
}
