package gormDB

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func openMysql(dsn string) (*gorm.DB, error) {
	return gorm.Open(mysql.Open(dsn), &gorm.Config{
		PrepareStmt: false,
	})
}

func InitGormDB(host, port, user, password, dbName string, models []interface{}) (db *gorm.DB, err error) {
	// "mysql://user:password@tcp(127.0.0.1:3306)/DBName?charset=utf8mb4&parseTime=true&loc=UTC"
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=true&loc=UTC",
		user, password, host, port, dbName,
	)
	db, err = openMysql(dsn)
	if err != nil {
		return nil, err
	}
	err = db.AutoMigrate(models...)
	return db, err
}
