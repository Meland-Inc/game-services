package accountDB

import (
	"os"

	"gorm.io/gorm"

	"github.com/Meland-Inc/game-services/src/common/gormDB"
	"github.com/Meland-Inc/game-services/src/global/dbData"
)

var db *gorm.DB

func GetAccountDB() *gorm.DB {
	return db
}

func getDbTableModels() []interface{} {
	return []interface{}{
		dbData.PlayerRow{},
	}
}

func Init() (err error) {
	host := os.Getenv("MELAND_ACCOUNT_DB_HOST")
	port := os.Getenv("MELAND_ACCOUNT_DB_PORT")
	user := os.Getenv("MELAND_ACCOUNT_DB_USER")
	pass := os.Getenv("MELAND_ACCOUNT_DB_PASS")
	dbName := os.Getenv("MELAND_ACCOUNT_DB_DATABASE")
	db, err = gormDB.InitGormDB(host, port, user, pass, dbName, getDbTableModels())
	return err
}
