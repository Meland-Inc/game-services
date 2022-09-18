package msgChannel

import (
	"game-message-core/grpc/methodData"
	"game-message-core/proto"

	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	"github.com/Meland-Inc/game-services/src/global/auth"
	gameDb "github.com/Meland-Inc/game-services/src/global/gameDB"
	dbData "github.com/Meland-Inc/game-services/src/global/gameDB/data"
	"github.com/spf13/cast"
	"gorm.io/gorm"
)

type HandleFunc func(*methodData.PullClientMessageInput, *proto.Envelope)

func (ch *MsgChannel) registerHandler() {
	ch.msgHandler[proto.EnvelopeType_CreatePlayer] = ch.createPlayerHandle
}

func (ch *MsgChannel) createPlayerHandle(
	input *methodData.PullClientMessageInput,
	msg *proto.Envelope,
) {
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

	userIdStr, err := auth.CheckDefaultAuth(req.Token)
	if err != nil {
		respMsg.ErrorCode = 20001 // TODO: USE PROTO ERROR CODE
		respMsg.ErrorMessage = err.Error()
		return
	}
	userId := cast.ToInt64(userIdStr)

	player := &dbData.PlayerBaseData{}
	err = gameDb.GetGameDB().Where("user_id = ?", userId).First(player).Error
	if err != nil {
		respMsg.ErrorCode = 20003 // TODO: USE PROTO ERROR CODE
		respMsg.ErrorMessage = err.Error()
		if err == gorm.ErrRecordNotFound {
			respMsg.ErrorCode = 20004 // TODO: USE PROTO ERROR CODE
			respMsg.ErrorMessage = "user not found"
		}
		return
	}

	respMsg.Payload = &proto.Envelope_QueryPlayerResponse{
		QueryPlayerResponse: &proto.QueryPlayerResponse{
			Player: player.ToNetPlayerBaseData(),
		},
	}
}
