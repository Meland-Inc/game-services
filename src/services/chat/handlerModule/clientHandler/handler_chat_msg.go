package clientHandler

import (
	"game-message-core/grpc/methodData"
	"game-message-core/proto"

	"github.com/Meland-Inc/game-services/src/common/time_helper"
	"github.com/Meland-Inc/game-services/src/global/userAgent"
	"github.com/Meland-Inc/game-services/src/services/chat/chatModel"
)

func ChatMsgHandle(
	agent *userAgent.UserAgentData, input *methodData.PullClientMessageInput, msg *proto.Envelope,
) {
	res := &proto.SendChatMessageResponse{}
	respMsg := userAgent.MakeResponseMsg(msg)
	defer func() {
		if respMsg.ErrorMessage != "" {
			respMsg.ErrorCode = 20003 // TODO: USE PROTO ERROR CODE
		}
		respMsg.Payload = &proto.Envelope_SendChatMessageResponse{SendChatMessageResponse: res}
		userAgent.ResponseClientMessage(agent, input, respMsg)
	}()

	if input.UserId < 1 {
		respMsg.ErrorMessage = "Invalid User ID"
		return
	}

	model, err := chatModel.GetChatModel()
	res.MsgId = time_helper.NowUTCMicro()
	err = model.OnReceiveChatPbMsg(input.UserId, res.MsgId, msg)
	if err != nil {
		respMsg.ErrorMessage = err.Error()
	}
}
