package msgChannel

import (
	"game-message-core/grpc/methodData"
	"game-message-core/proto"

	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	"github.com/Meland-Inc/game-services/src/common/time_helper"
	"github.com/Meland-Inc/game-services/src/global/auth"
	gameDb "github.com/Meland-Inc/game-services/src/global/gameDB"
	dbData "github.com/Meland-Inc/game-services/src/global/gameDB/data"
	"github.com/spf13/cast"
	"gorm.io/gorm"
)

type HandleFunc func(*methodData.PullClientMessageInput, *proto.Envelope)

func (ch *MsgChannel) registerHandler() {
	ch.msgHandler[proto.EnvelopeType_CreatePlayer] = ch.CreatePlayerHandler
	ch.msgHandler[proto.EnvelopeType_QueryPlayer] = ch.QueryPlayerHandler

}

func (ch *MsgChannel) QueryPlayerHandler(input *methodData.PullClientMessageInput, msg *proto.Envelope) {
	req := msg.GetQueryPlayerRequest()
	if req == nil {
		serviceLog.Error("account query player request is nil")
		return
	}

	respMsg := &proto.Envelope{
		Type:  msg.Type,
		SeqId: msg.SeqId,
		Payload: &proto.Envelope_QueryPlayerResponse{
			QueryPlayerResponse: &proto.QueryPlayerResponse{},
		},
	}

	defer func() {
		if respMsg.ErrorMessage != "" {
			serviceLog.Error("query player err : %s", respMsg.ErrorMessage)
		}
		ch.SendToPlayer(input.AgentAppId, input.SocketId, 0, respMsg)
	}()
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

	serviceLog.Info(
		"QueryPlayer------agent[%s],userId[%s],socketId[%s]",
		input.AgentAppId, userIdStr, input.SocketId,
	)

	userId := cast.ToInt64(userIdStr)
	player := &dbData.PlayerBaseData{}
	err = gameDb.GetGameDB().Where("user_id = ?", userId).First(player).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		respMsg.ErrorCode = 20003 // TODO: USE PROTO ERROR CODE
		respMsg.ErrorMessage = err.Error()
		return
	}

	respMsg.Payload = &proto.Envelope_QueryPlayerResponse{
		QueryPlayerResponse: &proto.QueryPlayerResponse{
			Player: player.ToNetPlayerBaseData(),
		},
	}
}

func (ch *MsgChannel) CreatePlayerHandler(input *methodData.PullClientMessageInput, msg *proto.Envelope) {
	req := msg.GetCreatePlayerRequest()
	if req == nil {
		serviceLog.Error("account create player request is nil")
		return
	}

	respMsg := &proto.Envelope{
		Type:  msg.Type,
		SeqId: msg.SeqId,
	}
	defer func() {
		if respMsg.ErrorMessage != "" {
			serviceLog.Error("create player err : %s", respMsg.ErrorMessage)
		}
		ch.SendToPlayer(input.AgentAppId, input.SocketId, 0, respMsg)
	}()

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
	err = gameDb.GetGameDB().Where("user_id = ?", userId).First(player).Error
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

	err = gameDb.GetGameDB().Create(player).Error
	if err != nil {
		respMsg.ErrorCode = 20003 // TODO: USE PROTO ERROR CODE
		respMsg.ErrorMessage = err.Error()
		return
	}

	respMsg.Payload = &proto.Envelope_CreatePlayerResponse{
		CreatePlayerResponse: &proto.CreatePlayerResponse{
			Player: player.ToNetPlayerBaseData(),
		},
	}
}
