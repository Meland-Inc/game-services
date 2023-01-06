package clientHandler

import (
	"game-message-core/grpc/methodData"
	"game-message-core/proto"

	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	"github.com/Meland-Inc/game-services/src/global/auth"
	"github.com/Meland-Inc/game-services/src/global/gameDB"
	dbData "github.com/Meland-Inc/game-services/src/global/gameDB/data"
	"github.com/Meland-Inc/game-services/src/global/userAgent"
	"github.com/spf13/cast"
	"gorm.io/gorm"
)

func QueryPlayerHandler(
	agent *userAgent.UserAgentData, input *methodData.PullClientMessageInput, msg *proto.Envelope,
) {
	res := &proto.QueryPlayerResponse{}
	respMsg := userAgent.MakeResponseMsg(msg)
	defer func() {
		if respMsg.ErrorMessage != "" {
			respMsg.ErrorCode = 20000 // TODO: USE PROTO ERROR CODE
			serviceLog.Error("Query Player err: %s", respMsg.ErrorMessage)
		}
		respMsg.Payload = &proto.Envelope_QueryPlayerResponse{QueryPlayerResponse: res}
		userAgent.ResponseClientMessage(agent, input, respMsg)
	}()

	req := msg.GetQueryPlayerRequest()
	if req == nil {
		serviceLog.Error("account query player request is nil")
		return
	}

	serviceLog.Info(
		"QueryPlayer------agent[%s], socketId[%s], token: %s",
		input.AgentAppId, input.SocketId, req.Token,
	)

	userIdStr, err := auth.CheckDefaultAuth(req.Token)
	if err != nil {
		respMsg.ErrorCode = 20001 // TODO: USE PROTO ERROR CODE
		respMsg.ErrorMessage = err.Error()
		return
	}

	userId := cast.ToInt64(userIdStr)
	player := &dbData.PlayerBaseData{}
	err = gameDB.GetGameDB().Where("user_id = ?", userId).First(player).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		respMsg.ErrorCode = 20003 // TODO: USE PROTO ERROR CODE
		respMsg.ErrorMessage = err.Error()
		return
	}

	res.Player = player.ToNetPlayerBaseData()
}
