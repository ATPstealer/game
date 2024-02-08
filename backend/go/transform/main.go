package main

import (
	"backend/packages/cfg"
	"backend/packages/db"
	"backend/transform/evolution"
	"gorm.io/gorm"
)

func main() {
	cfg.LoadConfig()                 // global cfg.Config
	db.ConnectToDatabase(cfg.Config) // global db.DB
	transform(db.DB)
}

func transform(db *gorm.DB) {
	evolution.CellAveragePrices(db)
	evolution.CellSpendMax(db)
	evolution.ResetStats(db)
	evolution.Hiring(db)
	evolution.Payroll(db)
}
