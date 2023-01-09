package clientHandler

import (
	"game-message-core/grpc/methodData"
	"game-message-core/proto"

	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	"github.com/Meland-Inc/game-services/src/common/time_helper"
	"github.com/Meland-Inc/game-services/src/global/auth"
	"github.com/Meland-Inc/game-services/src/global/gameDB"
	dbData "github.com/Meland-Inc/game-services/src/global/gameDB/data"
	"github.com/Meland-Inc/game-services/src/global/userAgent"
	"github.com/spf13/cast"
	"gorm.io/gorm"
)

func CreatePlayerHandler(
	agent *userAgent.UserAgentData, input *methodData.PullClientMessageInput, msg *proto.Envelope,
) {
	res := &proto.CreatePlayerResponse{}
	respMsg := userAgent.MakeResponseMsg(msg)
	defer func() {
		if respMsg.ErrorMessage != "" {
			respMsg.ErrorCode = 20000 // TODO: USE PROTO ERROR CODE
			serviceLog.Error("create Player err: %s", respMsg.ErrorMessage)
		}
		respMsg.Payload = &proto.Envelope_CreatePlayerResponse{CreatePlayerResponse: res}
		userAgent.ResponseClientMessage(agent, input, respMsg)
	}()

	req := msg.GetCreatePlayerRequest()
	if req == nil {
		serviceLog.Error("account create player request is nil")
		return
	}

	serviceLog.Info(
		"CreatePlayer--agent[%s], socketId[%s], token: %s",
		input.AgentAppId, input.SocketId, req.Token,
	)

	userIdStr, err := auth.CheckDefaultAuth(req.Token)
	if err != nil {
		respMsg.ErrorCode = 20001 // TODO: USE PROTO ERROR CODE
		respMsg.ErrorMessage = err.Error()
		return
	}

	serviceLog.Info(
		"CreatePlayer--agent[%s],userId[%s],socketId[%s]",
		input.AgentAppId, userIdStr, input.SocketId,
	)

	userId := cast.ToInt64(userIdStr)
	player := &dbData.PlayerBaseData{}
	err = gameDB.GetGameDB().Where("user_id = ?", userId).First(player).Error
	if err == nil {
		respMsg.ErrorCode = 20002 // TODO: USE PROTO ERROR CODE
		respMsg.ErrorMessage = "user already in the database"
		return
	}
	if err != gorm.ErrRecordNotFound {
		respMsg.ErrorCode = 20003 // TODO: USE PROTO ERROR CODE
		respMsg.ErrorMessage = err.Error()
		return
	}

	player.UserId = userId
	player.Name = req.NickName
	player.RoleId = req.RoleId
	player.RoleIcon = req.Icon
	player.SetFeature(req.Feature)
	player.CreatedAt = time_helper.NowUTC()
	player.UpdateAt = time_helper.NowUTC()

	err = gameDB.GetGameDB().Create(player).Error
	if err != nil {
		respMsg.ErrorCode = 20003 // TODO: USE PROTO ERROR CODE
		respMsg.ErrorMessage = err.Error()
		return
	}
	res.Player = player.ToNetPlayerBaseData()
}
