package dbData

import (
	"game-message-core/proto"
	"time"

	"github.com/Meland-Inc/game-services/src/common/time_helper"
)

type NftBuild struct {
	UId       uint      `gorm:"primaryKey;autoIncrement" json:"uid,string"`
	Owner     int64     `json:"owner"`
	NftId     string    `json:"nftId"`
	EntityId  int64     `json:"entityId"`
	Cid       int32     `json:"cid"`
	MapId     int32     `json:"mapId"`
	X         float32   `json:"x"`
	Y         float32   `json:"y"`
	Z         float32   `json:"z"`
	DirX      float32   `json:"dirX"`
	DirY      float32   `json:"dirY"`
	DirZ      float32   `json:"dirZ"`
	CreatedAt time.Time `json:"createdAt"`
	UpdateAt  time.Time `json:"updateAt"`
}

func NewNftBuild(
	userId int64, nftId string, cid, mapId int32, pos *proto.Vector3, lands []int32,
) *NftBuild {
	nowTm := time_helper.NowUTC()
	data := &NftBuild{
		Owner:     userId,
		NftId:     nftId,
		EntityId:  time_helper.NowUTCMicro(),
		Cid:       cid,
		MapId:     mapId,
		X:         pos.X,
		Y:         pos.Y,
		Z:         pos.Z,
		DirX:      0,
		DirY:      0,
		DirZ:      1,
		CreatedAt: nowTm,
		UpdateAt:  nowTm,
	}
	return data
}
