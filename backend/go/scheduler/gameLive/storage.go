package gameLive

import (
	"backend/packages/models"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/gorm"
	"log"
)

func StoragesUpdate(db *gorm.DB) {
	storages, err := models.GetAllStorages(db)
	if err != nil {
		log.Println(err)
	}
	cleanStorage(&storages)

	resources, err := models.GetAllResources(db)
	if err != nil {
		log.Println(err)
	}
	for _, resource := range resources {
		findResource(&storages, resource)
	}

	limit := -1
	sell := true
	orders, err := models.GetOrders(db, models.FindOrderParams{Limit: &limit, Sell: &sell})
	if err != nil {
		log.Println(err)
	}
	for _, order := range orders {
		findOrder(&storages, order)
	}

	buildingStorages, err := models.GetAllReadyStorages(db)
	if err != nil {
		log.Println(err)
	}

	// Storage size depend on workers count
	storageBuildingType, _ := models.GetBuildingTypeByID(db, 1) // 1 = Storage
	for _, buildingStorage := range buildingStorages {
		findBuildingStorage(&storages, buildingStorage, float64(storageBuildingType.Workers))
	}
	db.Save(&storages)
}

func cleanStorage(storages *[]models.Storage) {
	for i := 0; i < len(*storages); i++ {
		(*storages)[i].VolumeMax = 0
		(*storages)[i].VolumeOccupied = 0
	}
}

func findResource(storages *[]models.Storage, resource models.ResourceResult) {
	for i, storage := range *storages {
		if storage.UserID == resource.UserID && storage.X == resource.X && storage.Y == resource.Y {
			(*storages)[i].VolumeOccupied += resource.Amount * resource.Volume
			return
		}
	}
	*storages = append(*storages, models.Storage{
		UserID:         resource.UserID,
		VolumeOccupied: resource.Volume * resource.Amount,
		X:              resource.X,
		Y:              resource.Y,
	})
}

func findOrder(storages *[]models.Storage, order models.OrdersResult) {
	for i, storage := range *storages {
		if storage.UserID == order.UserID && storage.X == order.X && storage.Y == order.Y {
			(*storages)[i].VolumeOccupied += order.Amount * order.Volume
			return
		}
	}
	*storages = append(*storages, models.Storage{
		UserID:         order.UserID,
		VolumeOccupied: order.Volume * order.Amount,
		X:              order.X,
		Y:              order.Y,
	})
}

func findBuildingStorage(storages *[]models.Storage, buildingStorage models.Building, workersNeeded float64) {
	for i, storage := range *storages {
		if storage.UserID == buildingStorage.UserID && storage.X == buildingStorage.X && storage.Y == buildingStorage.Y {
			(*storages)[i].VolumeMax += float64(buildingStorage.Workers) * 100 * 5 / workersNeeded
			return
		}
	}
	*storages = append(*storages, models.Storage{
		UserID:    buildingStorage.UserID,
		VolumeMax: float64(buildingStorage.Workers) * 100 * 5 / workersNeeded,
		X:         buildingStorage.X,
		Y:         buildingStorage.Y,
	})
}

// mongo

func StoragesUpdateMongo(m *mongo.Database) {
	storages, err := models.GetAllStoragesMongo(m)
	if err != nil {
		log.Println(err)
	}
	cleanStorageMongo(&storages)

	resources, err := models.GetAllResourcesMongo(m)
	if err != nil {
		log.Println(err)
	}
	for _, resource := range resources {
		findResourceMongo(&storages, resource)
	}

	buildingStorages, err := models.GetAllReadyStoragesMongo(m)
	if err != nil {
		log.Println(err)
	}

	log.Println(buildingStorages)

	// Storage size depend on workers count
	storageBuildingType, _ := models.GetBuildingTypeByIDMongo(m, 1) // 1 = Storage
	for _, buildingStorage := range buildingStorages {
		findBuildingStorageMongo(&storages, buildingStorage, float64(storageBuildingType.Workers))
	}

	for _, storage := range storages {
		log.Println(storage)
	}

	limit := 9000000000000000000 // infinite // TODO: сделать че-нибудь нормальное
	sell := true
	orders, err := models.GetOrdersMongo(m, models.FindOrderParamsMongo{Limit: &limit, Sell: &sell})
	log.Println(orders)
	if err != nil {
		log.Println(err)
	}
	for _, order := range orders {
		findOrderMongo(&storages, order)
	}

	for _, storage := range storages {
		if storage.ID == primitive.NilObjectID {
			storage.ID = primitive.NewObjectID()
		}
		_, err := m.Collection("storages").UpdateOne(context.TODO(),
			bson.M{"_id": storage.ID},
			bson.M{"$set": storage},
			options.Update().SetUpsert(true))
		if err != nil {
			log.Println(err)
		}
	}
}

func cleanStorageMongo(storages *[]models.StorageMongo) {
	for i := 0; i < len(*storages); i++ {
		(*storages)[i].VolumeMax = 0
		(*storages)[i].VolumeOccupied = 0
	}
}

func findResourceMongo(storages *[]models.StorageMongo, resource models.ResourceWithTypeMongo) {
	for i, storage := range *storages {
		if storage.UserID == resource.UserID && storage.X == resource.X && storage.Y == resource.Y {
			(*storages)[i].VolumeOccupied += resource.ResourceType.Volume * resource.Amount
			return
		}
	}
	*storages = append(*storages, models.StorageMongo{
		UserID:         resource.UserID,
		VolumeOccupied: resource.ResourceType.Volume * resource.Amount,
		X:              resource.X,
		Y:              resource.Y,
	})
}

func findBuildingStorageMongo(storages *[]models.StorageMongo, buildingStorage models.BuildingMongo, workersNeeded float64) {
	for i, storage := range *storages {
		if storage.UserID == buildingStorage.UserID && storage.X == buildingStorage.X && storage.Y == buildingStorage.Y {
			(*storages)[i].VolumeMax += float64(buildingStorage.Workers) * 100 * 5 / workersNeeded
			return
		}
	}
	*storages = append(*storages, models.StorageMongo{
		UserID:    buildingStorage.UserID,
		VolumeMax: float64(buildingStorage.Workers) * 100 * 5 / workersNeeded,
		X:         buildingStorage.X,
		Y:         buildingStorage.Y,
	})
}

func findOrderMongo(storages *[]models.StorageMongo, order models.OrderMongoWithData) {
	for i, storage := range *storages {
		if storage.UserID == order.UserID && storage.X == order.X && storage.Y == order.Y {
			(*storages)[i].VolumeOccupied += order.Amount * order.ResourceType.Volume
			return
		}
	}
	*storages = append(*storages, models.StorageMongo{
		UserID:         order.UserID,
		VolumeOccupied: order.ResourceType.Volume * order.Amount,
		X:              order.X,
		Y:              order.Y,
	})
}
