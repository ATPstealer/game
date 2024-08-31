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
	transform(db.M)
}

func transform(m *mongo.Database) {
	log.Println("=== Calculating prices in cells and maximum expenses of the population ===")
	evolution.CellAveragePrices(m)
	evolution.CellSpendMax(m)
	log.Println("=== Reset stats ===")
	evolution.ResetStats(m)
	log.Println("=== Hiring ===")
	evolution.Hiring(m)
	log.Println("=== Payroll ===")
	evolution.Payroll(m)
	log.Println("=== Durability Recount ===")
	evolution.DurabilityRecount(m)
	log.Println("=== Logistics Reset ===")
	evolution.LogisticsReset(m)
}
