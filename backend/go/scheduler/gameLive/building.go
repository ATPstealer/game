package gameLive

import (
	"backend/packages/models"
	"gorm.io/gorm"
	"log"
	"time"
)

type ProductionResult struct {
	ID          uint
	BuildingID  uint
	BlueprintID uint
	WorkStarted *time.Time
	WorkEnd     *time.Time
	UserID      uint
	X           int
	Y           int
	Square      int
	Level       int
	Status      string
}

func Production(db *gorm.DB) {
	var productions []ProductionResult
	res := db.Model(&models.Production{}).Select("productions.id", "productions.building_id",
		"productions.blueprint_id", "productions.work_started", "productions.work_end", "buildings.user_id",
		"buildings.x", "buildings.y", "buildings.square", "buildings.level", "buildings.status").
		Joins("left join buildings on buildings.id = productions.building_id").
		Where("NOW() < productions.work_end").Scan(&productions)
	if res.Error != nil {
		log.Println(res.Error)
	}

	var resources []models.Resource
	res = db.Model(&models.Resource{}).Find(&resources)
	if res.Error != nil {
		log.Println(res.Error)
	}

	blueprintResults, err := models.GetBlueprints(db, 0)
	if err != nil {
		log.Println(err)
	}

	log.Println(blueprintResults)
	log.Println(productions)

	now := time.Now()
	for _, production := range productions {
		log.Println(production)
		if !models.CheckEnoughStorage(db, production.UserID, production.X, production.Y, 0) {
			db.Model(&models.Building{}).Where("id = ?", production.BuildingID).Update("status", models.StorageNeededStatus)
			db.Model(&models.Production{}).Where("id = ?", production.ID).Update("work_started", &now)
			continue
		}

		workTime := now.Sub(*production.WorkStarted).Seconds()
		blueprint := blueprintResults[production.BlueprintID-1]

		log.Println(blueprint)
		log.Println(workTime)

		// Formula of production pace
		cycles := int(workTime / blueprint.ProductionTime.Seconds())
		productionCycles := cycles * production.Level * production.Square

		if productionCycles == 0 {
			continue
		}

		enoughResources := true
		for _, resource := range blueprint.UsedResources {
			if !models.CheckEnoughResources(db, resource.ResourceID, production.UserID, production.X, production.Y, resource.Amount*float32(productionCycles)) {
				db.Model(&models.Building{}).Where("id = ?", production.BuildingID).Update("status", models.ResourcesNeededStatus)
				db.Model(&models.Production{}).Where("id = ?", production.ID).Update("work_started", &now)
				enoughResources = false
				break
			}
		}

		if enoughResources {
			for _, resource := range blueprint.UsedResources {
				err := models.AddResource(db, resource.ResourceID, production.UserID, production.X, production.Y,
					(-1)*resource.Amount*float32(productionCycles))
				if err != nil {
					log.Println(err)
				}
			}
			for _, resource := range blueprint.ProducedResources {
				err := models.AddResource(db, resource.ResourceID, production.UserID, production.X, production.Y,
					resource.Amount*float32(productionCycles))
				if err != nil {
					log.Println(err)
				}
			}
			db.Model(&models.Building{}).Where("id = ?", production.BuildingID).Update("status", models.ProductionStatus)
			newWorkStarted := production.WorkStarted.Add(time.Duration(cycles) * blueprint.ProductionTime)
			db.Model(&models.Production{}).Where("id = ?", production.ID).Update("work_started", &newWorkStarted)
		}
	}
}

/*
func Production(db *gorm.DB) {
	var buildings []models.Building
	res := db.Model(&models.Building{}).Where("status IN (?)  AND NOW() < work_end",
		[]models.BuildingStatus{models.ProductionStatus, models.ResourcesNeededStatus, models.StorageNeededStatus}).Scan(&buildings)
	if res.Error != nil {
		log.Println(res.Error)
	}

	var resources []models.Resource
	res = db.Model(&models.Resource{}).Find(&resources)
	if res.Error != nil {
		log.Println(res.Error)
	}

	blueprintResults, err := models.GetBlueprints(db, 0)
	if err != nil {
		log.Println(err)
	}

	now := time.Now()
	for _, building := range buildings {
		if !models.CheckEnoughStorage(db, building.UserID, building.X, building.Y, 0) {
			building.WorkStarted = &now
			building.Status = models.StorageNeededStatus
			db.Save(&building)
			continue
		}

		workTime := now.Sub(*building.WorkStarted).Seconds()
		blueprint := blueprintResults[building.ProductionID-1]

		// Formula of production pace
		cycles := int(workTime / blueprint.ProductionTime.Seconds())
		productionCycles := cycles * building.Level * building.Square

		if productionCycles == 0 {
			continue
		}
		enoughResources := true
		for _, resource := range blueprint.UsedResources {
			if !models.CheckEnoughResources(db, resource.ResourceID, building.UserID, building.X, building.Y, resource.Amount*float32(productionCycles)) {
				building.WorkStarted = &now
				building.Status = models.ResourcesNeededStatus
				db.Save(&building)
				enoughResources = false
				break
			}
		}
		if enoughResources {
			for _, resource := range blueprint.UsedResources {
				err := models.AddResource(db, resource.ResourceID, building.UserID, building.X, building.Y,
					(-1)*resource.Amount*float32(productionCycles))
				if err != nil {
					log.Println(err)
				}
			}
			for _, resource := range blueprint.ProducedResources {
				err := models.AddResource(db, resource.ResourceID, building.UserID, building.X, building.Y,
					resource.Amount*float32(productionCycles))
				if err != nil {
					log.Println(err)
				}
			}
			building.Status = models.ProductionStatus
			newWorkStarted := building.WorkStarted.Add(time.Duration(cycles) * blueprint.ProductionTime)
			building.WorkStarted = &newWorkStarted
			db.Save(building)
		}
	}
}*/

func StopWork(db *gorm.DB) {
	res := db.Model(&models.Building{}).Where("status IN (?) AND work_end < NOW()",
		[]models.BuildingStatus{models.ConstructionStatus, models.ProductionStatus, models.ResourcesNeededStatus, models.StorageNeededStatus}).
		Updates(map[string]interface{}{
			"status":       models.ReadyStatus,
			"work_end":     nil,
			"work_started": nil,
		})
	if res.Error != nil {
		log.Println(res.Error)
	}
	db.Model(&models.Production{}).Where("work_end < NOW()").Delete(&models.Production{})
}
