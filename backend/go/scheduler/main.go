package main

import (
	"backend/packages/cfg"
	"backend/packages/db"
	"backend/scheduler/gameLive"
	"gorm.io/gorm"
)

func main() {
	cfg.LoadConfig()                 // global cfg.Config
	db.ConnectToDatabase(cfg.Config) // global db.DB
	alive(db.DB)
}

func alive(db *gorm.DB) {
	gameLive.Production(db)
	gameLive.StopWork(db)
	gameLive.LogisticsDone(db)
	gameLive.StoragesUpdate(db)
	gameLive.StoreSell(db)
}
