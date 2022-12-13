package home_model

import (
	"fmt"
	base_data "game-message-core/grpc/baseData"
	"os"
	"testing"

	"github.com/Meland-Inc/game-services/src/common/time_helper"
	"github.com/Meland-Inc/game-services/src/global/gameDB"
	dbData "github.com/Meland-Inc/game-services/src/global/gameDB/data"
	"gorm.io/gorm"
)

func makeDb() error {
	os.Setenv("GAME_DB_HOST", "127.0.0.1")
	os.Setenv("GAME_DB_PORT", "3306")
	os.Setenv("GAME_DB_USER", "root")
	os.Setenv("GAME_DB_PASS", "root")
	os.Setenv("GAME_DB_DATABASE", "meland_game_data_dev")
	return gameDB.Init()
}

func GetUserHomeData(userId int64) (data *dbData.HomeData, err error) {
	data = dbData.NewHomeData(userId, "", "", "")
	err = gameDB.GetGameDB().Where("user_id = ?", userId).FirstOrCreate(data).Error
	return data, err
}

func UpdateUserHomeData(userId int64, data base_data.GrpcHomeData) error {
	home, err := GetUserHomeData(userId)
	if err != nil {
		return err
	}

	return gameDB.GetGameDB().Transaction(func(tx *gorm.DB) error {
		home.SoilJson = data.SoilJson
		home.ResourceJson = data.ResourceJson
		home.UpdateAt = time_helper.NowUTC()
		return tx.Save(home).Error
	})

	// return gameDB.GetGameDB().Transaction(func(tx *gorm.DB) error {
	// 	return tx.Model(&dbData.HomeData{}).Where("user_id=?", userId).Updates(
	// 		map[string]interface{}{
	// 			"soil_json":      data.SoilJson,
	// 			"livestock_json": "",
	// 			"resource_json":  data.ResourceJson,
	// 		}).Error
	// })
}

func Test_FindOrInitHome(t *testing.T) {
	t.Log(makeDb())
	var userId int64 = 696
	// homeData, err := GetUserHomeData(userId)
	// t.Log(err)
	// t.Log(fmt.Sprintf("- %+v", homeData))

	// for i := 0; i < 10000; i++ {
	UpdateUserHomeData(userId, base_data.GrpcHomeData{
		SoilJson:     "xxxxxxxxx",
		ResourceJson: fmt.Sprintf("%d", time_helper.NowUTCMill()),
	})
	// t.Log(err1)
	// }
	// homeData, err = GetUserHomeData(userId)
	// t.Log(err)
	// t.Log(fmt.Sprintf("- %+v", homeData))
}
