package userAgent

import (
	"game-message-core/grpc"
	"game-message-core/proto"
	"game-message-core/protoTool"

	"github.com/Meland-Inc/game-services/src/common/daprInvoke"
	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	"github.com/Meland-Inc/game-services/src/common/time_helper"
)

type UserAgentData struct {
	AgentAppId          string `json:"agentAppId"`
	InSceneServiceAppId string `json:"inSceneServiceAppId"`
	InMapId             int32  `json:"inMapId"`
	SocketId            string `json:"socketId"`
	UserId              int64  `json:"userId"`
	LoginAt             int64  `json:"loginAt"`
}

func (p *UserAgentData) TryUpdate(userId int64, agentAppId, socketId string) {
	if p.UserId == 0 && userId > 0 {
		p.UserId = userId
	}

	if socketId != "" && socketId != p.SocketId {
		p.SocketId = socketId
	}
	if agentAppId != "" && agentAppId != p.AgentAppId {
		p.AgentAppId = agentAppId
	}
}

func (p *UserAgentData) SendToPlayer(serviceAppId string, msg *proto.Envelope) error {
	input := &proto.BroadCastToClientInput{
		MsgVersion:   time_helper.NowUTCMill(),
		ServiceAppId: serviceAppId,
		UserId:       p.UserId,
		SocketId:     p.SocketId,
		MsgId:        int32(msg.Type),
		Msg:          msg,
	}

	inputBytes, err := protoTool.MarshalProto(input)
	if err != nil {
		serviceLog.Error("SendToPlayer Marshal BroadCastInput failed err: %+v", err)
		return err
	}

	_, err = daprInvoke.InvokeMethod(
		p.AgentAppId,
		string(grpc.ProtoMessageActionBroadCastToClient),
		inputBytes,
	)
	if err != nil {
		serviceLog.Error("UserAgentData SendToPlayer InvokeMethod  failed err:%+v", err)
	}
	return err
}
