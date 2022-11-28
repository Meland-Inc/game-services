package gameDB

import (
	"os"

	"gorm.io/gorm"

	"github.com/Meland-Inc/game-services/src/common/gormDB"
	dbData "github.com/Meland-Inc/game-services/src/global/gameDB/data"
)

var db *gorm.DB

func GetGameDB() *gorm.DB {
	return db
}

func getDbTableModels() []interface{} {
	return []interface{}{
		dbData.PlayerBaseData{},
		dbData.PlayerSceneData{},
		dbData.UsingNft{},
		dbData.ItemSlot{},
		dbData.PlayerTask{},
		dbData.NftBuild{},
		dbData.LoginData{},
	}
}

func Init() (err error) {
	host := os.Getenv("GAME_DB_HOST")
	port := os.Getenv("GAME_DB_PORT")
	user := os.Getenv("GAME_DB_USER")
	pass := os.Getenv("GAME_DB_PASS")
	dbName := os.Getenv("GAME_DB_DATABASE")
	db, err = gormDB.InitGormDB(host, port, user, pass, dbName, getDbTableModels())
	return err
}
