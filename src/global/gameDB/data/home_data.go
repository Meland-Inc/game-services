package dbData

import (
	"time"

	"github.com/Meland-Inc/game-services/src/common/time_helper"
)

type HomeData struct {
	UserId        int64     `gorm:"primaryKey" json:"userId"`
	SoilJson      string    `json:"soilJson"`  // 土地
	LivestockJson string    `json:"livestock"` // 家畜
	CreatedAt     time.Time `json:"createdAt"`
	UpdateAt      time.Time `json:"updateAt"`
}

func NewHomeData(userId int64, soilJson, livestockJson string) *HomeData {
	data := &HomeData{
		UserId:        userId,
		SoilJson:      soilJson,
		LivestockJson: livestockJson,
		CreatedAt:     time_helper.NowUTC(),
	}
	data.UpdateAt = data.CreatedAt
	return data
}
