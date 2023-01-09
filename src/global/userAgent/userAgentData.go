package userAgent

import (
	"game-message-core/proto"

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

func NewUserAgentData(
	userId int64, agentAppId, socketId, sceneAppId string,
) *UserAgentData {
	return &UserAgentData{
		AgentAppId:          agentAppId,
		SocketId:            socketId,
		InSceneServiceAppId: sceneAppId,
		UserId:              userId,
		LoginAt:             time_helper.NowUTCMill(),
	}
}

func (p *UserAgentData) TryUpdate(userId int64, agentAppId, socketId, sceneAppId string) {
	if p.UserId == 0 && userId > 0 {
		p.UserId = userId
	}

	if socketId != "" && socketId != p.SocketId {
		p.SocketId = socketId
	}
	if agentAppId != "" && agentAppId != p.AgentAppId {
		p.AgentAppId = agentAppId
	}
	if sceneAppId != "" && sceneAppId != p.InSceneServiceAppId {
		p.InSceneServiceAppId = sceneAppId
	}
}

func (p *UserAgentData) SendToPlayer(serviceAppId string, msg *proto.Envelope) error {
	err := BroadCastToClient(p.AgentAppId, serviceAppId, p.UserId, p.SocketId, msg)
	if err != nil {
		serviceLog.Error("UserAgentData SendToPlayer InvokeMethod  failed err:%+v", err)
	}
	return err
}
