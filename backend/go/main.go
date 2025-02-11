package main

import (
	_ "backend/docs"
	"backend/packages/cfg"
	"backend/packages/db"
	"backend/packages/router"
	"log"
)

//	@title		Game API
//	@version	2.0
//	@host		staging.game.k8s.atpstealer.com
//	@BasePath	/api/v2

func main() {
	cfg.LoadConfig()            // global cfg.Config
	db.MongoConnect(cfg.Config) // global db.M

	r := router.MakeRouter() // register controllers

	err := r.Run(":" + cfg.Config.AppPort) // run web server
	if err != nil {
		log.Fatal("Web server has failed")
	}
}
