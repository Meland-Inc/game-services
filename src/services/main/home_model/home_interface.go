package home_model

import (
	"github.com/Meland-Inc/game-services/src/global/gameDB"
	dbData "github.com/Meland-Inc/game-services/src/global/gameDB/data"
	"gorm.io/gorm"
)

func (p *HomeModel) GetUserHomeData(userId int64) (data *dbData.HomeData, err error) {
	data = dbData.NewHomeData(userId, "", "")
	err = gameDB.GetGameDB().Where("user_id = ?", userId).FirstOrCreate(data).Error
	return data, err
}

func (p *HomeModel) UpdateUserHomeData(userId int64, soilJson, livestockJson string) error {
	data, err := p.GetUserHomeData(userId)
	if err != nil {
		return err
	}

	return gameDB.GetGameDB().Transaction(func(tx *gorm.DB) error {
		data.SoilJson = soilJson
		data.LivestockJson = livestockJson
		return tx.Save(data).Error
	})
}
