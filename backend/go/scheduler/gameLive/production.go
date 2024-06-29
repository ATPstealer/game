package gameLive

import (
	"backend/packages/models"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"time"
)

func Production(m *mongo.Database) {
	buildings, err := models.GetProduction(m)
	if err != nil {
		log.Println("Can't get productions: " + err.Error())
		return
	}

	blueprints, err := models.GetBlueprints(m, 0)
	if err != nil {
		log.Println("Can't get blueprints: " + err.Error())
		return
	}

	now := time.Now()
	for _, building := range buildings {
		if !models.CheckEnoughStorage(m, building.UserId, building.X, building.Y, 0) {
			if err := models.BuildingStatusUpdate(m, building.Id, models.StorageNeededStatus); err != nil {
				log.Println("Can't update building status: " + err.Error())
			}
			if err := models.BuildingSetWorkStarted(m, building.Id, now); err != nil {
				log.Println("Can't update production start time: " + err.Error())
			}
			continue
		}
		workTime := now.Sub(building.WorkStarted).Seconds()
		blueprint := blueprints[building.Production.BlueprintId-1]

		log.Println(workTime)

		// Formula of production pace
		equipmentImpact := getEquipmentImpact(building, blueprint.Id)
		// The level and square are taken into account through workers
		productionCycles := int((workTime / blueprint.ProductionTime.Seconds()) * (float64(building.Workers) + equipmentImpact) / float64(building.BuildingType.Workers))
		blueprintCycles := float64(productionCycles) * float64(building.BuildingType.Workers) / (float64(building.Workers) + equipmentImpact)
		log.Print(productionCycles, blueprintCycles)

		if productionCycles == 0 {
			continue
		}

		if building.OnStrike {
			if err := models.BuildingSetWorkStarted(m, building.Id, now); err != nil {
				log.Println("Can't update production start time: " + err.Error())
			}
			continue
		}

		enoughResources := true
		for _, resource := range blueprint.UsedResources {
			if !models.CheckEnoughResources(m, resource.ResourceId, building.UserId, building.X, building.Y, resource.Amount*float64(productionCycles)) {
				if err := models.BuildingStatusUpdate(m, building.Id, models.ResourcesNeededStatus); err != nil {
					log.Println("Can't update building status: " + err.Error())
				}
				if err := models.BuildingSetWorkStarted(m, building.Id, now); err != nil {
					log.Println("Can't update production start time: " + err.Error())
				}
				enoughResources = false
				break
			}
		}

		if enoughResources {
			for _, resource := range blueprint.UsedResources {
				if err := models.AddResource(m, resource.ResourceId, building.UserId, building.X,
					building.Y, (-1)*resource.Amount*float64(productionCycles)); err != nil {
					log.Println("Can't add resources: " + err.Error())
				}
			}
			for _, resource := range blueprint.ProducedResources {
				if err := models.AddResource(m, resource.ResourceId, building.UserId, building.X,
					building.Y, resource.Amount*float64(productionCycles)); err != nil {
					log.Println("Can't add resources: " + err.Error())
				}
			}
			if err := models.BuildingStatusUpdate(m, building.Id, models.ProductionStatus); err != nil {
				log.Println("Can't update building status: " + err.Error())
			}
			newWorkStarted := building.WorkStarted.Add(time.Duration(blueprintCycles * float64(blueprint.ProductionTime)))
			if err := models.BuildingSetWorkStarted(m, building.Id, newWorkStarted); err != nil {
				log.Println("Can't update production start time: " + err.Error())
			}
		}

	}
}

func getEquipmentImpact(building models.BuildingWithData, blueprintId uint) float64 {
	if building.EquipmentEffect == nil {
		return 0
	}
	equipmentImpact := 0.0
	for _, equipmentEffect := range *building.EquipmentEffect {
		if equipmentEffect.EffectId == 1 && (equipmentEffect.BlueprintId == 0 || equipmentEffect.BlueprintId == blueprintId) {
			equipmentImpact += equipmentEffect.Value
		}
	}
	log.Print(equipmentImpact)
	return equipmentImpact
}

func StopWork(m *mongo.Database) {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(60*time.Second))
	defer cancel()

	filter := bson.D{{"workEnd", bson.D{{"$lt", time.Now()}}}}
	update := bson.D{
		{"$set", bson.D{
			{"status", models.ReadyStatus},
			{"workEnd", nil},
			{"workStarted", nil},
			{"production", nil},
		}},
	}
	_, err := m.Collection("buildings").UpdateMany(ctx, filter, update)
	if err != nil {
		log.Println("Production: " + err.Error())
		return
	}

	_, err = m.Collection("productions").DeleteMany(ctx, filter)
	if err != nil {
		log.Println("Production: " + err.Error())
		return
	}
}
