package main

import (
	"backend/packages/cfg"
	"backend/packages/db"
	"backend/scheduler/gameLive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"time"
)

func main() {
	cfg.LoadConfig()            // global cfg.Config
	db.MongoConnect(cfg.Config) // global db.M
	for {
		start := time.Now()
		log.Printf(start.String())
		alive(db.M)

		select {
		case <-time.After(time.Until(start.Add(1 * time.Minute))):
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
