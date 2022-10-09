package clientMsgHandle

import (
	"game-message-core/grpc/methodData"
	"game-message-core/proto"

	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	"github.com/Meland-Inc/game-services/src/services/main/playerModel"
)

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

	err = dataModel.LoadAvatar(input.UserId, req.ItemId)
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

	err = dataModel.UnloadAvatar(input.UserId, req.ItemId, true)
	if err != nil {
		respMsg.ErrorMessage = err.Error()
		return
	}
}
