package gameLive

import (
	"backend/packages/models"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
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
	OnStrike    bool
	Workers     int
	WorkersNeed int
}

func Production(db *gorm.DB) {
	var productions []ProductionResult
	res := db.Model(&models.Production{}).Select("productions.id", "productions.building_id",
		"productions.blueprint_id", "productions.work_started", "productions.work_end", "buildings.user_id",
		"buildings.x", "buildings.y", "buildings.square", "buildings.level", "buildings.status", "buildings.on_strike",
		"buildings.workers AS workers", "building_types.workers AS workers_need").
		Joins("left join buildings on buildings.id = productions.building_id").
		Joins("left join building_types on buildings.type_id = building_types.id").
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

	now := time.Now()
	for _, production := range productions {
		if !models.CheckEnoughStorage(db, production.UserID, production.X, production.Y, 0) {
			db.Model(&models.Building{}).Where("id = ?", production.BuildingID).Update("status", models.StorageNeededStatus)
			db.Model(&models.Production{}).Where("id = ?", production.ID).Update("work_started", &now)
			continue
		}

		workTime := now.Sub(*production.WorkStarted).Seconds()
		blueprint := blueprintResults[production.BlueprintID-1]

		// Formula of production pace
		productionCycles := int((workTime / blueprint.ProductionTime.Seconds()) * float64(production.Workers) / float64(production.WorkersNeed)) // here the level and square are taken into account through workers
		blueprintCycles := float64(productionCycles) * float64(production.WorkersNeed) / float64(production.Workers)

		if productionCycles == 0 {
			continue
		}

		if production.OnStrike {
			db.Model(&models.Production{}).Where("id = ?", production.ID).Update("work_started", &now)
			continue
		}

		enoughResources := true
		for _, resource := range blueprint.UsedResources {
			if !models.CheckEnoughResources(db, resource.ResourceID, production.UserID, production.X, production.Y, resource.Amount*float64(productionCycles)) {
				db.Model(&models.Building{}).Where("id = ?", production.BuildingID).Update("status", models.ResourcesNeededStatus)
				db.Model(&models.Production{}).Where("id = ?", production.ID).Update("work_started", &now)
				enoughResources = false
				break
			}
		}

		if enoughResources {
			for _, resource := range blueprint.UsedResources {
				err := models.AddResource(db, resource.ResourceID, production.UserID, production.X, production.Y,
					(-1)*resource.Amount*float64(productionCycles))
				if err != nil {
					log.Println(err)
				}
			}
			for _, resource := range blueprint.ProducedResources {
				err := models.AddResource(db, resource.ResourceID, production.UserID, production.X, production.Y,
					resource.Amount*float64(productionCycles))
				if err != nil {
					log.Println(err)
				}
			}
			db.Model(&models.Building{}).Where("id = ?", production.BuildingID).Update("status", models.ProductionStatus)
			newWorkStarted := production.WorkStarted.Add(time.Duration(blueprintCycles) * blueprint.ProductionTime)
			db.Model(&models.Production{}).Where("id = ?", production.ID).Update("work_started", &newWorkStarted)
		}
	}
}

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

// mongo

func ProductionMongo(m *mongo.Database) {
	productions, err := models.GetProductionMongo(m)
	if err != nil {
		log.Println("Can't get productions: " + err.Error())
		return
	}
	log.Println(productions)

	resources, err := models.GetAllResourcesMongo(m)
	if err != nil {
		log.Println("Can't get resources: " + err.Error())
		return
	}
	log.Println(resources)

	blueprints, err := models.GetBlueprintsMongo(m, 0)
	if err != nil {
		log.Println(err)
	}

	log.Println(blueprints)

	now := time.Now()
	for _, production := range productions {
		if !models.CheckEnoughStorageMongo(m, production.Building.UserID, production.Building.X, production.Building.Y, 0) {
			if err := models.BuildingStatusUpdate(m, production.Building.ID, models.StorageNeededStatus); err != nil {
				log.Println("Can't update building status: " + err.Error())
			}
			if err := models.ProductionSetWorkStarted(m, production.ID, &now); err != nil {
				log.Println("Can't update production start time: " + err.Error())
			}
			continue
		}
		workTime := now.Sub(*production.WorkStarted).Seconds()
		blueprint := blueprints[production.BlueprintID-1]
		log.Println(workTime, blueprint)

		// Formula of production pace
		productionCycles := int((workTime / blueprint.ProductionTime.Seconds()) * float64(production.Building.Workers) / float64(production.BuildingType.Workers)) // here the level and square are taken into account through workers
		blueprintCycles := float64(productionCycles) * float64(production.BuildingType.Workers) / float64(production.Building.Workers)
		log.Println(productionCycles, blueprintCycles)

		if productionCycles == 0 {
			continue
		}

		if production.Building.OnStrike {
			if err := models.ProductionSetWorkStarted(m, production.ID, &now); err != nil {
				log.Println("Can't update production start time: " + err.Error())
			}
			continue
		}

		enoughResources := true
		for _, resource := range blueprint.UsedResources {
			if !models.CheckEnoughResourcesMongo(m, resource.ResourceID, production.Building.UserID, production.Building.X, production.Building.Y, resource.Amount*float64(productionCycles)) {
				if err := models.BuildingStatusUpdate(m, production.Building.ID, models.ResourcesNeededStatus); err != nil {
					log.Println("Can't update building status: " + err.Error())
				}
				if err := models.ProductionSetWorkStarted(m, production.ID, &now); err != nil {
					log.Println("Can't update production start time: " + err.Error())
				}
				enoughResources = false
				break
			}
		}

		if enoughResources {
			for _, resource := range blueprint.UsedResources {
				if err := models.AddResourceMongo(m, resource.ResourceID, production.Building.UserID, production.Building.X,
					production.Building.Y, (-1)*resource.Amount*float64(productionCycles)); err != nil {
					log.Println("Can't add resources: " + err.Error())
				}
			}
			for _, resource := range blueprint.ProducedResources {
				if err := models.AddResourceMongo(m, resource.ResourceID, production.Building.UserID, production.Building.X,
					production.Building.Y, resource.Amount*float64(productionCycles)); err != nil {
					log.Println("Can't add resources: " + err.Error())
				}
			}
			if err := models.BuildingStatusUpdate(m, production.Building.ID, models.ProductionStatus); err != nil {
				log.Println("Can't update building status: " + err.Error())
			}
			newWorkStarted := production.WorkStarted.Add(time.Duration(blueprintCycles) * blueprint.ProductionTime)
			if err := models.ProductionSetWorkStarted(m, production.ID, &newWorkStarted); err != nil {
				log.Println("Can't update production start time: " + err.Error())
			}
		}

	}
}

func StopWorkMongo(m *mongo.Database) {
	filter := bson.D{{"workEnd", bson.D{{"$lt", time.Now()}}}}
	update := bson.D{
		{"$set", bson.D{
			{"status", models.ReadyStatus},
			{"workEnd", nil},
			{"workStarted", nil},
		}},
	}
	_, err := m.Collection("buildings").UpdateMany(context.TODO(), filter, update)
	if err != nil {
		log.Println("Production: " + err.Error())
		return
	}

	_, err = m.Collection("productions").DeleteMany(context.TODO(), filter)
	if err != nil {
		log.Println("Production: " + err.Error())
		return
	}
}
