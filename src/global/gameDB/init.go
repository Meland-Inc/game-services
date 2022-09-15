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
	}
}

func Init() (err error) {
	host := os.Getenv("MELAND_GAME_DB_HOST")
	port := os.Getenv("MELAND_GAME_DB_PORT")
	user := os.Getenv("MELAND_GAME_DB_USER")
	pass := os.Getenv("MELAND_GAME_DB_PASS")
	dbName := os.Getenv("MELAND_GAME_DB_DATABASE")
	db, err = gormDB.InitGormDB(host, port, user, pass, dbName, getDbTableModels())
	return err
}
