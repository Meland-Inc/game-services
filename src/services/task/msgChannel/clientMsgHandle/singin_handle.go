package clientMsgHandle

import (
	"game-message-core/grpc/methodData"
	"game-message-core/proto"

	"github.com/Meland-Inc/game-services/src/common/serviceLog"
)

func SelfTasksHandler(input *methodData.PullClientMessageInput, msg *proto.Envelope) {
	res := &proto.SigninPlayerResponse{}
	respMsg := makeResponseMsg(msg)
	defer func() {
		if respMsg.ErrorMessage != "" {
			respMsg.ErrorCode = 20001 // TODO: USE PROTO ERROR CODE
			serviceLog.Error("main service SingIn Player err: %s", respMsg.ErrorMessage)
		}
		respMsg.Payload = &proto.Envelope_SigninPlayerResponse{SigninPlayerResponse: res}
		ResponseClientMessage(input, respMsg)
	}()

	req := msg.GetSigninPlayerRequest()
	if req == nil {
		serviceLog.Error("main service singIn player request is nil")
		return
	}

}
