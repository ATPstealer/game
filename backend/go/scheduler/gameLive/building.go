package gameLive

import (
	"backend/packages/models"
	"gorm.io/gorm"
	"log"
	"time"
)

func Production(db *gorm.DB) {
	var buildings []models.Building
	res := db.Model(&models.Building{}).Where("status IN (?)  AND NOW() < work_end",
		[]models.BuildingStatus{models.Production, models.ResourcesNeeded, models.StorageNeeded}).Scan(&buildings)
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
			building.Status = models.StorageNeeded
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
				building.Status = models.ResourcesNeeded
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
			building.Status = models.Production
			newWorkStarted := building.WorkStarted.Add(time.Duration(cycles) * blueprint.ProductionTime)
			building.WorkStarted = &newWorkStarted
			db.Save(building)
		}
	}
}

func StopWork(db *gorm.DB) {
	res := db.Model(&models.Building{}).Where("status IN (?) AND work_end < NOW()",
		[]models.BuildingStatus{models.Construction, models.Production, models.ResourcesNeeded, models.StorageNeeded}).
		Updates(map[string]interface{}{
			"status":        models.Ready,
			"work_end":      nil,
			"work_started":  nil,
			"production_id": 0,
		})
	if res.Error != nil {
		log.Println(res.Error)
	}
}
