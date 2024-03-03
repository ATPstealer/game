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
	aliveMongo(db.M)
}

func aliveMongo(m *mongo.Database) {
	gameLive.ProductionMongo(m)
	gameLive.StopWorkMongo(m)
	gameLive.LogisticsDoneMongo(m)
	gameLive.StoragesUpdateMongo(m)
	gameLive.StoreSellMongo(m)
}
