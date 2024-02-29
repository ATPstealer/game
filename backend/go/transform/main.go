package main

import (
	"backend/packages/cfg"
	"backend/packages/db"
	"backend/transform/evolution"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

func main() {
	cfg.LoadConfig()                 // global cfg.Config
	db.ConnectToDatabase(cfg.Config) // global db.DB
	db.MongoConnect(cfg.Config)      // global db.M

	//transform(db.DB)
	transformMongo(db.M)

}

func transform(db *gorm.DB) {
	evolution.CellAveragePrices(db)
	evolution.CellSpendMax(db)
	evolution.ResetStats(db)
	evolution.Hiring(db)
	evolution.Payroll(db)
}

func transformMongo(m *mongo.Database) {
	evolution.HiringMongo(m)
}
