package main

import (
	"backend/packages/cfg"
	"backend/packages/db"
	"backend/scheduler/gameLive"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

func main() {
	cfg.LoadConfig()                 // global cfg.Config
	db.ConnectToDatabase(cfg.Config) // global db.DB
	db.MongoConnect(cfg.Config)      // global db.M

	aliveMongo(db.M)
	// alive(db.DB)
}

func alive(db *gorm.DB) {
	gameLive.Production(db)
	gameLive.StopWork(db)
	gameLive.LogisticsDone(db)
	gameLive.StoragesUpdate(db)
	gameLive.StoreSell(db)
}

// mongo
func aliveMongo(m *mongo.Database) {
	gameLive.StoragesUpdateMongo(m)
	gameLive.LogisticsDoneMongo(m)
	gameLive.StopWorkMongo(m)
}
