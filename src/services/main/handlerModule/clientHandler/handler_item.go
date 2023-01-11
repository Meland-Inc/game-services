package clientHandler

import (
	"game-message-core/grpc/methodData"
	"game-message-core/proto"
	"game-message-core/protoTool"

	"github.com/Meland-Inc/game-services/src/common/net/msgParser"
	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	"github.com/Meland-Inc/game-services/src/global/userAgent"
	"github.com/Meland-Inc/game-services/src/services/main/playerModel"
)

func ItemGetHandler(
	agent *userAgent.UserAgentData, input *methodData.PullClientMessageInput, msg *proto.Envelope,
) {
	respMsg := userAgent.MakeResponseMsg(msg)
	defer func() {
		if respMsg.ErrorMessage != "" {
			respMsg.ErrorCode = 20002 // TODO: USE PROTO ERROR CODE
			serviceLog.Error(respMsg.ErrorMessage)
		}
		respMsg.Payload = &proto.Envelope_ItemGetResponse{
			ItemGetResponse: &proto.ItemGetResponse{},
		}
		userAgent.ResponseClientMessage(agent, input, respMsg)
	}()

	playerDataModel, err := playerModel.GetPlayerDataModel()
	if err != nil {
		respMsg.ErrorMessage = err.Error()
		return
	}

	playerItems, err := playerDataModel.GetPlayerItems(input.UserId)
	if err != nil {
		respMsg.ErrorMessage = err.Error()
		return
	}

	serviceLog.Info("main service userId[%v] itemLength[%v]", input.UserId, len(playerItems.Items))

	initRes := &proto.BroadCastInitItemResponse{Items: []*proto.Item{}}
	dbMsg := &proto.Envelope{
		Type: proto.EnvelopeType_BroadCastInitItem,
		Payload: &proto.Envelope_BroadCastInitItemResponse{
			BroadCastInitItemResponse: initRes,
		},
	}

	maxIdx := len(playerItems.Items) - 1
	for idx, it := range playerItems.Items {
		initRes.Items = append(initRes.Items, it.ToNetItem())
		var msgDataLength int
		if len(initRes.Items) >= 10 {
			msgBody, _ := protoTool.MarshalProto(dbMsg)
			msgDataLength = len(msgBody)
		}
		if idx >= maxIdx || msgDataLength >= msgParser.MSG_LIMIT-1000 {
			userAgent.ResponseClientMessage(agent, input, dbMsg)
			initRes.Items = []*proto.Item{}
		}
	}
}

func ItemUseHandler(
	agent *userAgent.UserAgentData, input *methodData.PullClientMessageInput, msg *proto.Envelope,
) {
	res := &proto.ItemUseResponse{}
	respMsg := userAgent.MakeResponseMsg(msg)
	defer func() {
		if respMsg.ErrorMessage != "" {
			respMsg.ErrorCode = 20004 // TODO: USE PROTO ERROR CODE
			serviceLog.Error(respMsg.ErrorMessage)
		}
		respMsg.Payload = &proto.Envelope_ItemUseResponse{ItemUseResponse: res}
		userAgent.ResponseClientMessage(agent, input, respMsg)
	}()

	if input.UserId < 1 {
		respMsg.ErrorMessage = "item use Invalid User ID"
		return
	}

	playerDataModel, err := playerModel.GetPlayerDataModel()
	if err != nil {
		respMsg.ErrorMessage = err.Error()
		return
	}

	req := msg.GetItemUseRequest()
	if req == nil {
		respMsg.ErrorMessage = "main service use item request is nil"
		return
	}

	err = playerDataModel.UseItem(input.UserId, req.ItemId, req.Args)
	if err != nil {
		respMsg.ErrorMessage = err.Error()
		return
	}
}

func LoadAvatarHandler(
	agent *userAgent.UserAgentData, input *methodData.PullClientMessageInput, msg *proto.Envelope,
) {
	res := &proto.UpdateAvatarResponse{}
	respMsg := userAgent.MakeResponseMsg(msg)
	defer func() {
		if respMsg.ErrorMessage != "" {
			respMsg.ErrorCode = 20004 // TODO: USE PROTO ERROR CODE
			serviceLog.Error(respMsg.ErrorMessage)
		}
		respMsg.Payload = &proto.Envelope_UpdateAvatarResponse{UpdateAvatarResponse: res}
		userAgent.ResponseClientMessage(agent, input, respMsg)
	}()

	if input.UserId < 1 {
		respMsg.ErrorMessage = "load avatar Invalid User ID"
		return
	}

	playerDataModel, err := playerModel.GetPlayerDataModel()
	if err != nil {
		respMsg.ErrorMessage = err.Error()
		return
	}

	req := msg.GetUpdateAvatarRequest()
	if req == nil {
		serviceLog.Error("main service load avatar request is nil")
		return
	}
	err = playerDataModel.LoadAvatar(input.UserId, req.ItemId, req.IsAppearance)
	if err != nil {
		respMsg.ErrorMessage = err.Error()
		return
	}
}

func UnloadAvatarHandler(
	agent *userAgent.UserAgentData, input *methodData.PullClientMessageInput, msg *proto.Envelope,
) {
	res := &proto.UnloadAvatarResponse{}
	respMsg := userAgent.MakeResponseMsg(msg)
	defer func() {
		if respMsg.ErrorMessage != "" {
			respMsg.ErrorCode = 20005 // TODO: USE PROTO ERROR CODE
		}
		respMsg.Payload = &proto.Envelope_UnloadAvatarResponse{UnloadAvatarResponse: res}
		userAgent.ResponseClientMessage(agent, input, respMsg)
	}()

	if input.UserId < 1 {
		respMsg.ErrorMessage = "Unload avatar Invalid User ID"
		return
	}

	playerDataModel, err := playerModel.GetPlayerDataModel()
	if err != nil {
		respMsg.ErrorMessage = err.Error()
		return
	}

	req := msg.GetUnloadAvatarRequest()
	if req == nil {
		serviceLog.Error("main service Unload avatar request is nil")
		return
	}
	err = playerDataModel.UnloadAvatar(input.UserId, req.ItemId, true, false)
	if err != nil {
		respMsg.ErrorMessage = err.Error()
		return
	}
}
