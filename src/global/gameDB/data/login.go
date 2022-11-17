package dbData

import (
	"time"

	"github.com/Meland-Inc/game-services/src/common/time_helper"
)

type LoginData struct {
	UserId         int64     `gorm:"primaryKey" json:"userId"`
	AgentAppId     string    `json:"agentAppId"`
	SocketId       string    `json:"socketId"`
	InSceneService string    `json:"inSceneService"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdateAt       time.Time `json:"updateAt"`
}

func NewLoginData(userId int64, agentAppId, socketId, inSceneService string) *LoginData {
	loginData := &LoginData{}
	loginData.UserId = userId
	loginData.AgentAppId = agentAppId
	loginData.SocketId = socketId
	loginData.InSceneService = inSceneService
	loginData.CreatedAt = time_helper.NowUTC()
	loginData.UpdateAt = loginData.CreatedAt
	return loginData
}
