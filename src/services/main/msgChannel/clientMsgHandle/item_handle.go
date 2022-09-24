package clientMsgHandle

import (
	"game-message-core/proto"

	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	"github.com/Meland-Inc/game-services/src/services/main/playerModel"
)

func ItemGetHandle(input *proto.PullClientMessageInput) {
	res := &proto.ItemGetResponse{}
	respMsg := makeResponseMsg(input.Msg)
	defer func() {
		if respMsg.ErrorMessage != "" {
			respMsg.ErrorCode = 20002 // TODO: USE PROTO ERROR CODE
		}
		respMsg.Payload = &proto.Envelope_ItemGetResponse{ItemGetResponse: res}
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

	playerItems, err := dataModel.GetPlayerItems(input.UserId)
	if err != nil {
		respMsg.ErrorMessage = err.Error()
		return
	}
	for _, it := range playerItems.Items {
		res.Items = append(res.Items, it.ToNetItem())
	}
}

func ItemUseHandle(input *proto.PullClientMessageInput) {
	res := &proto.ItemUseResponse{}
	respMsg := makeResponseMsg(input.Msg)
	defer func() {
		if respMsg.ErrorMessage != "" {
			respMsg.ErrorCode = 20004 // TODO: USE PROTO ERROR CODE
		}
		respMsg.Payload = &proto.Envelope_ItemUseResponse{ItemUseResponse: res}
		ResponseClientMessage(input, respMsg)
	}()

	if input.UserId < 1 {
		respMsg.ErrorMessage = "Invalid User ID"
		return
	}

	req := input.Msg.GetItemUseRequest()
	if req == nil {
		serviceLog.Error("main service use item request is nil")
		return
	}

	dataModel, err := playerModel.GetPlayerDataModel()
	if err != nil {
		respMsg.ErrorMessage = err.Error()
		return
	}

	dataModel.UseItem(input.UserId, req.ItemId)
}
