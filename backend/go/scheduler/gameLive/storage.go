package gameLive

import (
	"backend/packages/models"
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

	storageBuildingType, _ := models.GetBuildingTypeByID(db, 1) // 1 = Storage
	for _, buildingStorage := range buildingStorages {
		findBuildingStorage(&storages, buildingStorage, float64(storageBuildingType.Workers))
	}
	log.Println(storages)
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
