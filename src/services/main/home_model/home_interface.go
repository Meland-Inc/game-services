package home_model

import (
	base_data "game-message-core/grpc/baseData"

	"github.com/Meland-Inc/game-services/src/common/time_helper"
	"github.com/Meland-Inc/game-services/src/global/gameDB"
	dbData "github.com/Meland-Inc/game-services/src/global/gameDB/data"
	"gorm.io/gorm"
)

func ToGrpcHomeData(data dbData.HomeData) base_data.GrpcHomeData {
	return base_data.GrpcHomeData{
		SoilJson:     data.SoilJson,
		ResourceJson: data.ResourceJson,
	}
}

func (p *HomeModel) GetUserHomeData(userId int64) (data *dbData.HomeData, err error) {
	data = dbData.NewHomeData(userId, "", "", "")
	err = gameDB.GetGameDB().Where("user_id = ?", userId).FirstOrCreate(data).Error
	return data, err
}

func (p *HomeModel) UpdateUserHomeData(userId int64, data base_data.GrpcHomeData) error {
	home, err := p.GetUserHomeData(userId)
	if err != nil {
		return err
	}

	return gameDB.GetGameDB().Transaction(func(tx *gorm.DB) error {
		home.SoilJson = data.SoilJson
		home.ResourceJson = data.ResourceJson
		home.UpdateAt = time_helper.NowUTC()
		return tx.Save(home).Error
	})
}
