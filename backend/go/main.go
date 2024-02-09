package main

import (
	"backend/packages/cfg"
	"backend/packages/db"
	"backend/packages/models"
	"backend/packages/router"
	"log"
)

func main() {
	cfg.LoadConfig()                   // global cfg.Config
	db.MongoConnect(cfg.Config)        // global db.M
	db.ConnectToDatabase(cfg.Config)   // global db.DB
	models.AutoMigrateModel(db.DB)     // migrate database models
	models.InitMongo(db.M, cfg.Config) // Init mongoDB vars
	models.Init(db.DB, cfg.Config)     // Init database during first run
	models.DeleteObsoleteTokens(db.DB)

	r := router.MakeRouter()               // register controllers
	err := r.Run(":" + cfg.Config.AppPort) // run web server
	if err != nil {
		log.Fatal("Server has failed")
	}
}
