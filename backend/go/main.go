package main

import (
	"backend/packages/cfg"
	"backend/packages/db"
	"backend/packages/router"
	"log"
)

func main() {
	cfg.LoadConfig()            // global cfg.Config
	db.MongoConnect(cfg.Config) // global db.M

	r := router.MakeRouter()               // register controllers
	err := r.Run(":" + cfg.Config.AppPort) // run web server
	if err != nil {
		log.Fatal("Web server has failed")
	}
}
