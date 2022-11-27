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
		ItemGetResponse: &proto.ItemGetResponse{},
	}
	ResponseClientMessage(agent, input, respMsg)
}

func BroadCaseInitUserItem(
	input *methodData.PullClientMessageInput,
	agent *userAgent.UserAgentData,
	pbItems []*proto.Item,
) {
	if agent == nil {
		serviceLog.Warning("ItemGetGroupingResponse user [%d] agent data not found", input.UserId)
		return
	}

	initRes := &proto.BroadCastInitItemResponse{}
	msg := &proto.Envelope{
		Type: proto.EnvelopeType_BroadCastInitItem,
		Payload: &proto.Envelope_BroadCastInitItemResponse{
			BroadCastInitItemResponse: initRes,
		},
	}

	n := 12 // 单个protoItem 长度350B 20个=7000B
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
		initRes.Items = pbItems[beginIdx:endIdx]
		ResponseClientMessage(agent, input, msg)
	}
}

func ItemGetHandle(input *methodData.PullClientMessageInput, msg *proto.Envelope) {
	serviceLog.Info("main service userId[%v] get items ", input.UserId)
	agent := GetOrStoreUserAgent(input)
	if input.UserId < 1 {
		ItemGetResponse(input, agent, msg, "item Get Invalid User ID")
		return
	}

	dataModel, err := playerModel.GetPlayerDataModel()
	if err != nil {
		ItemGetResponse(input, agent, msg, err.Error())
		return
	}

	playerItems, err := dataModel.GetPlayerItems(input.UserId)
	if err != nil {
		ItemGetResponse(input, agent, msg, err.Error())
		return
	}

	pbItems := []*proto.Item{}
	for _, it := range playerItems.Items {
		pbItems = append(pbItems, it.ToNetItem())
	}

	serviceLog.Info("main service userId[%v] itemLength[%v]", input.UserId, len(pbItems))

	BroadCaseInitUserItem(input, agent, pbItems)
	ItemGetResponse(input, agent, msg, "")
}

func ItemUseHandle(input *methodData.PullClientMessageInput, msg *proto.Envelope) {
	agent := GetOrStoreUserAgent(input)
	res := &proto.ItemUseResponse{}
	respMsg := makeResponseMsg(msg)
	defer func() {
		if respMsg.ErrorMessage != "" {
			respMsg.ErrorCode = 20004 // TODO: USE PROTO ERROR CODE
			serviceLog.Error(respMsg.ErrorMessage)
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
		respMsg.ErrorMessage = "main service use item request is nil"
		return
	}

	dataModel, err := playerModel.GetPlayerDataModel()
	if err != nil {
		respMsg.ErrorMessage = err.Error()
		return
	}

	err = dataModel.UseItem(input.UserId, req.ItemId, req.Args)
	if err != nil {
		respMsg.ErrorMessage = err.Error()
		return
	}
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
