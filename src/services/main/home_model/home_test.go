package home_model

import (
	"os"
	"testing"

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

func Test_FindOrInitHome(t *testing.T) {
	t.Log(makeDb())
	var userId int64 = 699
	data := dbData.NewHomeData(userId, "33333", "超vv")
	err := gameDB.GetGameDB().Where("user_id = ?", userId).FirstOrCreate(data).Error
	t.Log(err)
	t.Log(data)

	err = gameDB.GetGameDB().Transaction(func(tx *gorm.DB) error {
		data.SoilJson = "BBBBBBBBBBBBBBBb"
		data.LivestockJson = "livestockJson"
		return tx.Save(data).Error
	})
	t.Log(err)
	t.Log(data)
}
