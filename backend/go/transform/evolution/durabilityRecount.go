package evolution

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

func DurabilityRecount(m *mongo.Database) {
	buildings, err := models.GetBuildingsForEquipmentRecount(m)
	if err != nil {
		log.Fatal(err)
	}
	equipmentTypes, err := models.GetAllEquipmentTypes(m)
	if err != nil {
		log.Fatal(err)
	}

	for _, building := range buildings {
		for _, equipment := range *building.Equipment {
			maxDurability := getMaxDurability(equipment.EquipmentTypeId, &equipmentTypes)
			if maxDurability == 0 {
				log.Fatalf("max durability not found: %v", building)
			}
			err = countAndSaveDurability(m, building.Id, equipment, maxDurability)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}

func getMaxDurability(equipmentId uint, equipmentTypes *[]models.EquipmentType) int {
	for _, equipmentType := range *equipmentTypes {
		if equipmentId == equipmentType.Id {
			return equipmentType.Durability
		}
	}
	return 0
}

func countAndSaveDurability(m *mongo.Database, buildingId primitive.ObjectID, equipment models.Equipment, maxDurability int) error {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(3*time.Second))
	defer cancel()

	newDurability := equipment.Durability - equipment.Amount
	newAmount := equipment.Amount

	if newDurability <= 0 {
		destroyedAmount := (-1)*int(newDurability/maxDurability) + 1
		newAmount = equipment.Amount - destroyedAmount
		newDurability += destroyedAmount * maxDurability
	}

	if newAmount <= 0 {
		update := bson.D{
			{"$pull", bson.D{{"equipment", bson.D{{"equipmentTypeId", equipment.EquipmentTypeId}}}}},
		}
		_, err := m.Collection("buildings").UpdateOne(ctx, bson.M{"_id": buildingId}, update)
		return err
	} else {
		update := bson.D{
			{"$set", bson.D{
				{"equipment.$[elem].durability", newDurability},
				{"equipment.$[elem].amount", newAmount},
			}},
		}
		filter := bson.D{{"_id", buildingId}}
		arrayFilters := options.ArrayFilters{
			Filters: []interface{}{
				bson.D{{"elem.equipmentTypeId", equipment.EquipmentTypeId}},
			},
		}
		updateOptions := options.UpdateOptions{
			ArrayFilters: &arrayFilters,
		}
		_, err := m.Collection("buildings").UpdateOne(ctx, filter, update, &updateOptions)
		return err
	}
}
