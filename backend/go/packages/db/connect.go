package db

import (
	"backend/packages/cfg"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	DB *gorm.DB
)

func ConnectToDatabase(config cfg.Vars) {
	dsn := config.DBUser + ":" + config.DBPassword + "@tcp(" + config.DBHost + ":" +
		config.DBPort + ")/" + config.DBDataBase + "?charset=utf8mb4&parseTime=True&loc=UTC"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		print(err)
	}
	db.Session(&gorm.Session{Logger: logger.Default.LogMode(logger.Info)})
	DB = db
}
