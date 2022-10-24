package clientMsgHandle

import (
	"game-message-core/grpc/methodData"
	"game-message-core/proto"

	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	"github.com/Meland-Inc/game-services/src/global/userAgent"
	"github.com/Meland-Inc/game-services/src/services/main/playerModel"
)

func ItemGetResponse(
	input *methodData.PullClientMessageInput,
	agent *userAgent.UserAgentData,
	reqMsg *proto.Envelope,
	errMsg string,
	pbItems []*proto.Item,
) {
	if agent == nil {
		serviceLog.Warning("ItemGetResponse user [%d] agent data not found", input.UserId)
		return
	}

	respMsg := makeResponseMsg(reqMsg)
	if respMsg.ErrorMessage != "" {
		respMsg.ErrorCode = 20002 // TODO: USE PROTO ERROR CODE
		serviceLog.Error(respMsg.ErrorMessage)
	}
	respMsg.Payload = &proto.Envelope_ItemGetResponse{
		ItemGetResponse: &proto.ItemGetResponse{
			Items: pbItems,
		}}
	ResponseClientMessage(agent, input, respMsg)
}

func ItemGetGroupingResponse(
	input *methodData.PullClientMessageInput,
	agent *userAgent.UserAgentData,
	reqMsg *proto.Envelope,
	pbItems []*proto.Item,
) {
	if agent == nil {
		serviceLog.Warning("ItemGetGroupingResponse user [%d] agent data not found", input.UserId)
		return
	}

	ItemGetResponse(input, agent, reqMsg, "", []*proto.Item{})

	addRes := &proto.BroadCastItemAddResponse{}
	msg := &proto.Envelope{
		Type: proto.EnvelopeType_BroadCastItemAdd,
		Payload: &proto.Envelope_BroadCastItemAddResponse{
			BroadCastItemAddResponse: addRes,
		},
	}

	n := 8
	itemLength := len(pbItems)
	left := itemLength / n
	if itemLength%n > 0 {
		left++
	}
	for i := 0; i < left; i++ {
		beginIdx := i * n
		endIdx := beginIdx + n
		if endIdx > itemLength {
			endIdx = itemLength
		}
		addRes.Items = pbItems[beginIdx:endIdx]
		ResponseClientMessage(agent, input, msg)
	}
}

func ItemGetHandle(input *methodData.PullClientMessageInput, msg *proto.Envelope) {
	serviceLog.Info("main service userId[%v] get items ", input.UserId)
	agent := GetOrStoreUserAgent(input)
	if input.UserId < 1 {
		ItemGetResponse(input, agent, msg, "item Get Invalid User ID", nil)
		return
	}

	dataModel, err := playerModel.GetPlayerDataModel()
	if err != nil {
		ItemGetResponse(input, agent, msg, err.Error(), nil)
		return
	}

	playerItems, err := dataModel.GetPlayerItems(input.UserId)
	if err != nil {
		ItemGetResponse(input, agent, msg, err.Error(), nil)
		return
	}

	pbItems := []*proto.Item{}
	for _, it := range playerItems.Items {
		pbItems = append(pbItems, it.ToNetItem())
	}

	serviceLog.Info("main service userId[%v] itemLength[%v]", input.UserId, len(pbItems))

	if len(pbItems) < 8 {
		ItemGetResponse(input, agent, msg, "", pbItems)
	} else {
		ItemGetGroupingResponse(input, agent, msg, pbItems)
	}
}

func ItemUseHandle(input *methodData.PullClientMessageInput, msg *proto.Envelope) {
	agent := GetOrStoreUserAgent(input)
	res := &proto.ItemUseResponse{}
	respMsg := makeResponseMsg(msg)
	defer func() {
		if respMsg.ErrorMessage != "" {
			respMsg.ErrorCode = 20004 // TODO: USE PROTO ERROR CODE
		}
		respMsg.Payload = &proto.Envelope_ItemUseResponse{ItemUseResponse: res}
		ResponseClientMessage(agent, input, respMsg)
	}()

	if input.UserId < 1 {
		respMsg.ErrorMessage = "item use Invalid User ID"
		return
	}

	req := msg.GetItemUseRequest()
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

func LoadAvatarHandle(input *methodData.PullClientMessageInput, msg *proto.Envelope) {
	agent := GetOrStoreUserAgent(input)
	res := &proto.UpdateAvatarResponse{}
	respMsg := makeResponseMsg(msg)
	defer func() {
		if respMsg.ErrorMessage != "" {
			respMsg.ErrorCode = 20005 // TODO: USE PROTO ERROR CODE
		}
		respMsg.Payload = &proto.Envelope_UpdateAvatarResponse{UpdateAvatarResponse: res}
		ResponseClientMessage(agent, input, respMsg)
	}()

	if input.UserId < 1 {
		respMsg.ErrorMessage = "load avatar Invalid User ID"
		return
	}

	req := msg.GetUpdateAvatarRequest()
	if req == nil {
		serviceLog.Error("main service load avatar request is nil")
		return
	}

	dataModel, err := playerModel.GetPlayerDataModel()
	if err != nil {
		respMsg.ErrorMessage = err.Error()
		return
	}

	err = dataModel.LoadAvatar(input.UserId, req.ItemId, req.IsAppearance)
	if err != nil {
		respMsg.ErrorMessage = err.Error()
		return
	}
}

func UnloadAvatarHandle(input *methodData.PullClientMessageInput, msg *proto.Envelope) {
	agent := GetOrStoreUserAgent(input)
	res := &proto.UnloadAvatarResponse{}
	respMsg := makeResponseMsg(msg)
	defer func() {
		if respMsg.ErrorMessage != "" {
			respMsg.ErrorCode = 20005 // TODO: USE PROTO ERROR CODE
		}
		respMsg.Payload = &proto.Envelope_UnloadAvatarResponse{UnloadAvatarResponse: res}
		ResponseClientMessage(agent, input, respMsg)
	}()

	if input.UserId < 1 {
		respMsg.ErrorMessage = "unload avatar Invalid User ID"
		return
	}

	req := msg.GetUnloadAvatarRequest()
	if req == nil {
		serviceLog.Error("main service unload avatar request is nil")
		return
	}

	dataModel, err := playerModel.GetPlayerDataModel()
	if err != nil {
		respMsg.ErrorMessage = err.Error()
		return
	}

	err = dataModel.UnloadAvatar(input.UserId, req.ItemId, true, false)
	if err != nil {
		respMsg.ErrorMessage = err.Error()
		return
	}
}
