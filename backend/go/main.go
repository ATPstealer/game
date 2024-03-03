package main

import (
	"backend/packages/cfg"
	"backend/packages/db"
	"backend/packages/models"
	"backend/packages/router"
	"log"
)

func main() {
	cfg.LoadConfig()              // global cfg.Config
	db.MongoConnect(cfg.Config)   // global db.M
	models.Init(db.M, cfg.Config) // Init mongoDB vars
	models.MongoIndex(db.M)
	// models.DeleteObsoleteTokens(db.DB) // TODO: make it

	r := router.MakeRouter()               // register controllers
	err := r.Run(":" + cfg.Config.AppPort) // run web server
	if err != nil {
		log.Fatal("Server has failed")
	}
}
