package clientHandler

import (
	"game-message-core/grpc/methodData"
	"game-message-core/proto"

	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	"github.com/Meland-Inc/game-services/src/global/userAgent"
)

func SingInHandler(
	agent *userAgent.UserAgentData, input *methodData.PullClientMessageInput, msg *proto.Envelope,
) {
	res := &proto.SigninPlayerResponse{}
	respMsg := userAgent.MakeResponseMsg(msg)
	defer func() {
		if respMsg.ErrorMessage != "" {
			respMsg.ErrorCode = 20001 // TODO: USE PROTO ERROR CODE
			serviceLog.Error("main service SingIn Player err: %s", respMsg.ErrorMessage)
		}
		respMsg.Payload = &proto.Envelope_SigninPlayerResponse{SigninPlayerResponse: res}
		userAgent.ResponseClientMessage(agent, input, respMsg)
	}()

	req := msg.GetSigninPlayerRequest()
	if req == nil {
		respMsg.ErrorMessage = "singIn player request is nil"
		serviceLog.Error(respMsg.ErrorMessage)
		return
	}
	// TODO logic
}
