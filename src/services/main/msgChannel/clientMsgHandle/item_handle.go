package clientMsgHandle

import (
	"game-message-core/grpc/methodData"
	"game-message-core/proto"
)

func ItemGetHandle(input *methodData.PullClientMessageInput, msg *proto.Envelope) {
	res := &proto.ItemGetResponse{}
	respMsg := makeResponseMsg(msg)
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

	dataModel, err := getPlayerDataModel()
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
