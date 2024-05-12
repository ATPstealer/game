package main

import (
	"backend/packages/cfg"
	"backend/packages/db"
	"backend/packages/models"
	"backend/scheduler/gameLive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"time"
)

func main() {
	cfg.LoadConfig()            // global cfg.Config
	db.MongoConnect(cfg.Config) // global db.M

	log.Printf("Run scheduler and init: " + time.Now().String())
	alive(db.M)
	models.Init(db.M, cfg.Config)
	log.Printf("Finish of scheduler and init: " + time.Now().String())

	everyMinute := time.NewTicker(1 * time.Minute)
	every10Minutes := time.NewTicker(10 * time.Minute)

	defer everyMinute.Stop()
	defer every10Minutes.Stop()

	for {
		select {
		case <-everyMinute.C:
			log.Printf("Scheduler start: " + time.Now().String())
			alive(db.M)
			log.Printf("Scheduler finish: " + time.Now().String())
		case <-every10Minutes.C:
			log.Printf("DB serve has been started: " + time.Now().String())
			models.Init(db.M, cfg.Config)     // Init database if Config.Init is True
			models.MongoIndex(db.M)           // Indexes for find nickname and email addresses due to registration
			models.DeleteObsoleteTokens(db.M) // clean old user's tokens
			log.Printf("DB serve has been finished: " + time.Now().String())
		}
	}
}

func alive(m *mongo.Database) {
	gameLive.Production(m)
	gameLive.StopWork(m)
	gameLive.LogisticsDone(m)
	gameLive.StoragesUpdate(m)
	gameLive.StoreSell(m)
}
