package main

import (
	"backend/packages/cfg"
	"backend/packages/db"
	"backend/scheduler/gameLive"
	"go.mongodb.org/mongo-driver/mongo"
)

func main() {
	cfg.LoadConfig()            // global cfg.Config
	db.MongoConnect(cfg.Config) // global db.M
	alive(db.M)
}

func alive(m *mongo.Database) {
	gameLive.Production(m)
	gameLive.StopWork(m)
	gameLive.LogisticsDone(m)
	gameLive.StoragesUpdate(m)
	gameLive.StoreSell(m)
}
