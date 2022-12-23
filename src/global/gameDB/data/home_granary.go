package dbData

import (
	"game-message-core/proto"
	"time"

	"github.com/Meland-Inc/game-services/src/common/time_helper"
)

type HomeGranary struct {
	UserId           int64     `gorm:"primaryKey;autoIncrement:false" json:"userId"`
	ItemCid          int32     `gorm:"primaryKey;autoIncrement:false" json:"itemCid"`
	Num              int32     `json:"num"`
	Quality          int32     `json:"quality"`
	LastPushUserId   int64     `json:"lastPushUserId"`
	LastPushUserName string    `json:"lastPushUserName"`
	CreatedAt        time.Time `json:"createdAt"`
	UpdateAt         time.Time `json:"updateAt"`
}

func NewHomeGranary(userId int64, itemCid, num, quality int32, lastPushUser int64, lastPushUserName string, upTm time.Time) *HomeGranary {
	row := &HomeGranary{
		UserId:           userId,
		ItemCid:          itemCid,
		Num:              num,
		Quality:          quality,
		LastPushUserId:   lastPushUser,
		LastPushUserName: lastPushUserName,
		UpdateAt:         upTm,
		CreatedAt:        time_helper.NowUTC(),
	}
	return row
}

func (p *HomeGranary) ToProtoData() *proto.ItemBaseInfo {
	return &proto.ItemBaseInfo{
		Cid:     p.ItemCid,
		Num:     p.Num,
		Quality: p.Quality,
	}
}
