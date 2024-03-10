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
	productions, err := models.GetProduction(m)
	if err != nil {
		log.Println("Can't get productions: " + err.Error())
		return
	}

	blueprints, err := models.GetBlueprints(m, 0)
	if err != nil {
		log.Println(err)
	}

	now := time.Now()
	for _, production := range productions {
		if !models.CheckEnoughStorage(m, production.Building.UserId, production.Building.X, production.Building.Y, 0) {
			if err := models.BuildingStatusUpdate(m, production.Building.Id, models.StorageNeededStatus); err != nil {
				log.Println("Can't update building status: " + err.Error())
			}
			if err := models.ProductionSetWorkStarted(m, production.Id, now); err != nil {
				log.Println("Can't update production start time: " + err.Error())
			}
			continue
		}
		workTime := now.Sub(production.WorkStarted).Seconds()
		blueprint := blueprints[production.BlueprintId-1]

		// Formula of production pace
		productionCycles := int((workTime / blueprint.ProductionTime.Seconds()) * float64(production.Building.Workers) / float64(production.BuildingType.Workers)) // here the level and square are taken into account through workers
		blueprintCycles := float64(productionCycles) * float64(production.BuildingType.Workers) / float64(production.Building.Workers)

		if productionCycles == 0 {
			continue
		}

		if production.Building.OnStrike {
			if err := models.ProductionSetWorkStarted(m, production.Id, now); err != nil {
				log.Println("Can't update production start time: " + err.Error())
			}
			continue
		}

		enoughResources := true
		for _, resource := range blueprint.UsedResources {
			if !models.CheckEnoughResources(m, resource.ResourceId, production.Building.UserId, production.Building.X, production.Building.Y, resource.Amount*float64(productionCycles)) {
				if err := models.BuildingStatusUpdate(m, production.Building.Id, models.ResourcesNeededStatus); err != nil {
					log.Println("Can't update building status: " + err.Error())
				}
				if err := models.ProductionSetWorkStarted(m, production.Id, now); err != nil {
					log.Println("Can't update production start time: " + err.Error())
				}
				enoughResources = false
				break
			}
		}

		if enoughResources {
			for _, resource := range blueprint.UsedResources {
				if err := models.AddResource(m, resource.ResourceId, production.Building.UserId, production.Building.X,
					production.Building.Y, (-1)*resource.Amount*float64(productionCycles)); err != nil {
					log.Println("Can't add resources: " + err.Error())
				}
			}
			for _, resource := range blueprint.ProducedResources {
				if err := models.AddResource(m, resource.ResourceId, production.Building.UserId, production.Building.X,
					production.Building.Y, resource.Amount*float64(productionCycles)); err != nil {
					log.Println("Can't add resources: " + err.Error())
				}
			}
			if err := models.BuildingStatusUpdate(m, production.Building.Id, models.ProductionStatus); err != nil {
				log.Println("Can't update building status: " + err.Error())
			}
			newWorkStarted := production.WorkStarted.Add(time.Duration(blueprintCycles) * blueprint.ProductionTime)
			if err := models.ProductionSetWorkStarted(m, production.Id, newWorkStarted); err != nil {
				log.Println("Can't update production start time: " + err.Error())
			}
		}

	}
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
