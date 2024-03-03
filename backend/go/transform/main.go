package main

import (
	"backend/packages/cfg"
	"backend/packages/db"
	"backend/transform/evolution"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
)

func main() {
	cfg.LoadConfig()            // global cfg.Config
	db.MongoConnect(cfg.Config) // global db.M
	transformMongo(db.M)
}

func transformMongo(m *mongo.Database) {
	log.Println("=== Calculating prices in cells and maximum expenses of the population ===")
	evolution.CellAveragePricesMongo(m)
	evolution.CellSpendMaxMongo(m)
	log.Println("=== Reset stats ===")
	evolution.ResetStatsMongo(m)
	log.Println("=== Hiring ===")
	evolution.HiringMongo(m)
	log.Println("=== Payroll ===")
	evolution.PayrollMongo(m)
}
