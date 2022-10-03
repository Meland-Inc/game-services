package userAgent

import (
	"game-message-core/proto"

	"github.com/Meland-Inc/game-services/src/common/serviceLog"
)

// "github.com/Meland-Inc/game-services/src/global/grpcAPI/grpcInvoke"

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
	err := BroadCastToClient(p.AgentAppId, serviceAppId, p.UserId, p.SocketId, msg)
	if err != nil {
		serviceLog.Error("UserAgentData SendToPlayer InvokeMethod  failed err:%+v", err)
	}
	return err
}
