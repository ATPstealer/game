package main

import (
	"backend/packages/cfg"
	"backend/packages/db"
	"backend/packages/models"
	"backend/packages/router"
	"log"
)

func main() {
	cfg.LoadConfig()                  // global cfg.Config
	db.MongoConnect(cfg.Config)       // global db.M
	models.Init(db.M, cfg.Config)     // Init mongoDB vars if variable Init is true
	models.MongoIndex(db.M)           // Indexes for find nickname and email addresses due to registration
	models.DeleteObsoleteTokens(db.M) // clean old user's tokens

	r := router.MakeRouter()               // register controllers
	err := r.Run(":" + cfg.Config.AppPort) // run web server
	if err != nil {
		log.Fatal("Web server has failed")
	}
}
