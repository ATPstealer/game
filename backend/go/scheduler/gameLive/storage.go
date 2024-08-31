package gameLive

import (
	"backend/packages/models"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

func StoragesUpdate(m *mongo.Database) {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(60*time.Second))
	defer cancel()

	storages, err := models.GetAllStorages(m)
	if err != nil {
		log.Println(err)
	}
	cleanStorage(&storages)

	resources, err := models.GetAllResources(m)
	if err != nil {
		log.Println(err)
	}
	for _, resource := range resources {
		findResource(&storages, resource)
	}

	buildingStorages, err := models.GetAllReadyBuildingByGroup(m, "Storage")
	if err != nil {
		log.Println(err)
	}

	// Storage size depend on workers count
	storageBuildingType, _ := models.GetBuildingTypeById(m, 1) // 1 = Storage
	for _, buildingStorage := range buildingStorages {
		findBuildingStorage(&storages, buildingStorage, float64(storageBuildingType.Workers))
	}

	limit := 9000000000000000000 // infinite // TODO: сделать че-нибудь нормальное
	sell := true
	orders, err := models.GetOrders(m, models.FindOrderParams{Limit: &limit, Sell: &sell})

	if err != nil {
		log.Println(err)
	}
	for _, order := range orders {
		findOrder(&storages, order)
	}

	for _, storage := range storages {
		if storage.Id == primitive.NilObjectID {
			storage.Id = primitive.NewObjectID()
		}
		_, err := m.Collection("storages").UpdateOne(ctx,
			bson.M{"_id": storage.Id},
			bson.M{"$set": storage},
			options.Update().SetUpsert(true))
		if err != nil {
			log.Println(err)
		}
	}
}

func cleanStorage(storages *[]models.Storage) {
	for i := 0; i < len(*storages); i++ {
		(*storages)[i].VolumeMax = 0
		(*storages)[i].VolumeOccupied = 0
	}
}

func findResource(storages *[]models.Storage, resource models.ResourceWithData) {
	for i, storage := range *storages {
		if storage.UserId == resource.UserId && storage.X == resource.X && storage.Y == resource.Y {
			(*storages)[i].VolumeOccupied += resource.ResourceType.Volume * resource.Amount
			return
		}
	}
	*storages = append(*storages, models.Storage{
		UserId:         resource.UserId,
		VolumeOccupied: resource.ResourceType.Volume * resource.Amount,
		X:              resource.X,
		Y:              resource.Y,
	})
}

func findBuildingStorage(storages *[]models.Storage, buildingStorage models.Building, workersNeeded float64) {
	for i, storage := range *storages {
		if storage.UserId == buildingStorage.UserId && storage.X == buildingStorage.X && storage.Y == buildingStorage.Y {
			(*storages)[i].VolumeMax += float64(buildingStorage.Workers) * 100 * 5 / workersNeeded
			return
		}
	}
	*storages = append(*storages, models.Storage{
		UserId:    buildingStorage.UserId,
		VolumeMax: float64(buildingStorage.Workers) * 100 * 5 / workersNeeded,
		X:         buildingStorage.X,
		Y:         buildingStorage.Y,
	})
}

func findOrder(storages *[]models.Storage, order models.OrderMongoWithData) {
	for i, storage := range *storages {
		if storage.UserId == order.UserId && storage.X == order.X && storage.Y == order.Y {
			(*storages)[i].VolumeOccupied += order.Amount * order.ResourceType.Volume
			return
		}
	}
	*storages = append(*storages, models.Storage{
		UserId:         order.UserId,
		VolumeOccupied: order.ResourceType.Volume * order.Amount,
		X:              order.X,
		Y:              order.Y,
	})
}
