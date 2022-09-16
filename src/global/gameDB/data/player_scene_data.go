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
	X           float64   `json:"x"`
	Y           float64   `json:"y"`
	Z           float64   `json:"z"`
	DirX        float64   `json:"dirX"`
	DirY        float64   `json:"dirY"`
	DirZ        float64   `json:"dirZ"`
	BirthMapId  int32     `json:"birthMapId"`
	BirthX      float64   `json:"birthX"`
	BirthY      float64   `json:"birthY"`
	BirthZ      float64   `json:"birthZ"`
	LastLoginAt time.Time `json:"lastLoginAt"`
}
