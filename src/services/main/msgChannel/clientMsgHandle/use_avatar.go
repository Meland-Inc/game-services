package clientMsgHandle

import (
	"game-message-core/grpc/methodData"
	"game-message-core/proto"

	"github.com/Meland-Inc/game-services/src/common/serviceLog"
)

func LoadAvatarHandle(input *methodData.PullClientMessageInput, msg *proto.Envelope) {
	res := &proto.UpdateAvatarResponse{}
	respMsg := makeResponseMsg(msg)
	defer func() {
		if respMsg.ErrorMessage != "" {
			respMsg.ErrorCode = 20005 // TODO: USE PROTO ERROR CODE
		}
		respMsg.Payload = &proto.Envelope_UpdateAvatarResponse{UpdateAvatarResponse: res}
		ResponseClientMessage(input, respMsg)
	}()

	if input.UserId < 1 {
		respMsg.ErrorMessage = "Invalid User ID"
		return
	}

	req := msg.GetUpdateAvatarRequest()
	if req == nil {
		serviceLog.Error("main service use item request is nil")
		return
	}

	dataModel, err := getPlayerDataModel()
	if err != nil {
		respMsg.ErrorMessage = err.Error()
		return
	}

	err = dataModel.LoadAvatar(input.UserId, req.ItemId, req.Position)
	if err != nil {
		respMsg.ErrorMessage = err.Error()
		return
	}
}

func UnloadAvatarHandle(input *methodData.PullClientMessageInput, msg *proto.Envelope) {
	res := &proto.UnloadAvatarResponse{}
	respMsg := makeResponseMsg(msg)
	defer func() {
		if respMsg.ErrorMessage != "" {
			respMsg.ErrorCode = 20005 // TODO: USE PROTO ERROR CODE
		}
		respMsg.Payload = &proto.Envelope_UnloadAvatarResponse{UnloadAvatarResponse: res}
		ResponseClientMessage(input, respMsg)
	}()

	if input.UserId < 1 {
		respMsg.ErrorMessage = "Invalid User ID"
		return
	}

	req := msg.GetUnloadAvatarRequest()
	if req == nil {
		serviceLog.Error("main service use item request is nil")
		return
	}

	dataModel, err := getPlayerDataModel()
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
