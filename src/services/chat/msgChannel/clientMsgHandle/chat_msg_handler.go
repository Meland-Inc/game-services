package clientMsgHandle

import (
	"game-message-core/grpc/methodData"
	"game-message-core/proto"

	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	"github.com/Meland-Inc/game-services/src/common/time_helper"
	"github.com/Meland-Inc/game-services/src/global/userAgent"
	"github.com/Meland-Inc/game-services/src/services/chat/chatModel"
)

func chatMsgResponse(
	input *methodData.PullClientMessageInput,
	agent *userAgent.UserAgentData,
	reqMsg *proto.Envelope,
	errMsg string,
	msgId int64,
) {
	if agent == nil {
		serviceLog.Warning("ItemGetResponse user [%d] agent data not found", input.UserId)
		return
	}

	respMsg := makeResponseMsg(reqMsg)
	if respMsg.ErrorMessage != "" {
		respMsg.ErrorCode = 50002 // TODO: USE PROTO ERROR CODE
		serviceLog.Error(respMsg.ErrorMessage)
	}
	respMsg.Payload = &proto.Envelope_SendChatMessageResponse{
		SendChatMessageResponse: &proto.SendChatMessageResponse{
			MsgId: msgId,
		},
	}
	ResponseClientMessage(agent, input, respMsg)
}

func ChatMsgHandle(input *methodData.PullClientMessageInput, msg *proto.Envelope) {
	serviceLog.Info("chat service send chat msg:%+v ", input)
	agent := GetOrStoreUserAgent(input)
	if input.UserId < 1 {
		chatMsgResponse(input, agent, msg, "Invalid User ID", 0)
		return
	}

	model, err := chatModel.GetChatModel()
	if err != nil {
		chatMsgResponse(input, agent, msg, "chat model not found", 0)
		return
	}

	msgId := time_helper.NowUTCMicro()
	err = model.OnReceiveChatPbMsg(input.UserId, msgId, msg)
	if err != nil {
		chatMsgResponse(input, agent, msg, err.Error(), 0)
		return
	}
	chatMsgResponse(input, agent, msg, "", msgId)
}
