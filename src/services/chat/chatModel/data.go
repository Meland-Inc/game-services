package chatModel

import (
	"game-message-core/proto"

	"github.com/Meland-Inc/game-services/src/common/time_helper"
	"github.com/Meland-Inc/game-services/src/global/configData"
)

type PlayerChatData struct {
	UserId            int64                           `json:"userId"`
	Name              string                          `json:"name"`
	RoleIcon          string                          `json:"roleIcon"`
	MapId             int32                           `json:"mapId"`
	X                 float32                         `json:"x"`
	Y                 float32                         `json:"y"`
	Z                 float32                         `json:"z"`
	SceneServiceAppId string                          `json:"sceneServiceAppId"`
	AgentAppId        string                          `json:"agentAppId"`
	UserSocketId      string                          `json:"userSocketId"`
	InGrid            *ViewGrid                       `json:"-"`
	ChatCDs           map[proto.ChatChannelType]int64 `json:"-"` // map[ChatChannelType]nextSendAt<ms>
}

func NewPlayerChatData(
	userId int64,
	name string,
	roleIcon string,
	mapId int32,
	x float32,
	y float32,
	z float32,
	sceneServiceAppId string,
	agentAppId string,
	userSocketId string,
) *PlayerChatData {
	return &PlayerChatData{
		UserId:            userId,
		Name:              name,
		RoleIcon:          roleIcon,
		MapId:             mapId,
		X:                 x,
		Y:                 y,
		Z:                 z,
		SceneServiceAppId: sceneServiceAppId,
		AgentAppId:        agentAppId,
		UserSocketId:      userSocketId,
		ChatCDs:           make(map[proto.ChatChannelType]int64),
	}
}

func (p *PlayerChatData) UpChatCD(chatType proto.ChatChannelType) {
	cnf := configData.ConfigMgr().ChatCnfByType(chatType)
	if cnf == nil {
		p.ChatCDs[chatType] = time_helper.NowUTCMill() + int64(cnf.Cd)
	}
}
