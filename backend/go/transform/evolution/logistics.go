package evolution

import (
	"backend/packages/models"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
)

func LogisticsReset(m *mongo.Database) {
	buildings, err := models.GetAllReadyBuildingByGroup(m, "Logistics")
	if err != nil {
		log.Println(err)
		return
	}
	for _, building := range buildings {
		err = models.LogisticsReset(m, building)
	}
}
