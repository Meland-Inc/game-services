package dbData

import (
	"time"
)

type PlayerSceneData struct {
	UId         uint      `gorm:"primaryKey;autoIncrement" json:"uid,string"`
	UserId      int64     `json:"userId"`
	Level       int64     `json:"level"`
	Exp         int32     `json:"exp"`
	MapId       int32     `json:"mapId"`
	X           float32   `json:"x"`
	Y           float32   `json:"y"`
	Z           float32   `json:"z"`
	DirX        float32   `json:"dirX"`
	DirY        float32   `json:"dirY"`
	DirZ        float32   `json:"dirZ"`
	BirthMapId  int32     `json:"birthMapId"`
	BirthX      float32   `json:"birthX"`
	BirthY      float32   `json:"birthY"`
	BirthZ      float32   `json:"birthZ"`
	LastLoginAt time.Time `json:"lastLoginAt"`
}
