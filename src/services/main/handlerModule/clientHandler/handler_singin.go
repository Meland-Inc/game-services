package clientHandler

import (
	"game-message-core/grpc/methodData"
	"game-message-core/proto"

	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	"github.com/Meland-Inc/game-services/src/common/time_helper"
	"github.com/Meland-Inc/game-services/src/global/auth"
	"github.com/Meland-Inc/game-services/src/global/userAgent"
	login_model "github.com/Meland-Inc/game-services/src/services/main/loginModel"
	"github.com/Meland-Inc/game-services/src/services/main/playerModel"
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

	userId, err := auth.GetUserIdByToken(req.Token)
	if err != nil {
		respMsg.ErrorMessage = err.Error()
		return
	}

	playerDataModel, err := playerModel.GetPlayerDataModel()
	if err != nil {
		respMsg.ErrorMessage = err.Error()
		return
	}

	loginModel, _ := login_model.GetLoginModel()
	sceneAppId, err := loginModel.GetUserLoginData(userId, input.AgentAppId, input.SocketId)
	if err != nil {
		respMsg.ErrorMessage = err.Error()
		return
	}
	if req.SceneServiceAppId != "" {
		sceneAppId = req.SceneServiceAppId
	}

	playerData, err := playerDataModel.PlayerProtoData(userId)
	if err != nil {
		respMsg.ErrorMessage = err.Error()
		return
	}

	// 登录时 input.UserId == 0 所以此处需要重新init userAgent
	input.UserId = userId
	agent = userAgent.GetOrStoreUserAgent(input)
	agent.InMapId = playerData.MapId

	res.SceneServiceAppId = sceneAppId
	res.ClientTime = req.ClientTime
	res.ServerTime = time_helper.NowUTCMill()
	res.Player = playerData
}
