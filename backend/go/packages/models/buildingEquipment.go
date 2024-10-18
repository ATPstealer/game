package models

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

type Equipment struct {
	EquipmentTypeId uint `json:"equipmentTypeId" bson:"equipmentTypeId" validate:"required"`
	Amount          int  `json:"amount" bson:"amount" validate:"required"`
	Durability      int  `json:"durability" bson:"durability" validate:"required"`
} // @name equipment

type InstallEquipmentPayload struct {
	BuildingId      primitive.ObjectID `json:"buildingId" validate:"required"`
	EquipmentTypeId uint               `json:"equipmentTypeId" validate:"required"`
	Amount          int                `json:"amount" validate:"required"`
} // @name installEquipmentPayload

// EffectsId https://docs.google.com/spreadsheets/d/18DblwMx9a-YnLh7hfqwSC5a3Sc-_bRybMWMWp-isgO0/edit?pli=1&gid=1818870012#gid=1818870012
const (
	WorkerProductivity = 1
	DecreasePollution  = 2
	LogisticsCapacity  = 3
	LoansLimit         = 4
	BorrowedLimit      = 5
)

func InstallEquipment(m *mongo.Database, userId primitive.ObjectID, installEquipment InstallEquipmentPayload) error {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(3*time.Second))
	defer cancel()

	var building Building
	err := m.Collection("buildings").FindOne(ctx, bson.M{"_id": installEquipment.BuildingId}).Decode(&building)
	if err != nil {
		return err
	}

	if building.UserId != userId {
		err := errors.New("this building don't belong you")
		log.Println(err)
		return err
	}

	equipmentType, err := GetEquipmentTypesByID(m, installEquipment.EquipmentTypeId)
	if err != nil {
		log.Println(err)
		return err
	}

	if equipmentType.Square*float64(installEquipment.Amount)+building.SquareInUse > float64(building.Square*building.Level) {
		return errors.New("not enough space")
	}

	index := getEquipmentPosition(building.Equipment, installEquipment.EquipmentTypeId)

	if installEquipment.Amount >= 0 {
		if !CheckEnoughResources(m, equipmentType.ResourceTypeId, userId, building.X, building.Y, float64(installEquipment.Amount)) {
			return errors.New("not enough resources in this cell")
		}
	} else {
		if index != -1 {
			if (*building.Equipment)[index].Amount < (-1)*installEquipment.Amount {
				return errors.New("not enough equipment here")
			}
		} else {
			return errors.New("not enough equipment here")
		}
	}

	resourceAdd := installEquipment.Amount
	if installEquipment.Amount < 0 {
		resourceAdd++ // If you uninstall the equipment, you lose one
	}
	err = AddResource(m, equipmentType.ResourceTypeId, userId, building.X, building.Y, float64(-resourceAdd))
	if err != nil {
		log.Println(err)
		return err
	}

	if index != -1 {
		err = updateEquipmentAmount(m, ctx, building, index, equipmentType.Id, installEquipment.Amount)
		if err != nil {
			return err
		}
	} else {
		newEquipment := Equipment{EquipmentTypeId: installEquipment.EquipmentTypeId, Amount: installEquipment.Amount, Durability: equipmentType.Durability}
		err = addEquipment(m, ctx, &building, newEquipment)
		if err != nil {
			return err
		}
	}

	err = CountEffects(m, building.Id)
	if err != nil {
		return err
	}

	_, err = m.Collection("buildings").UpdateOne(ctx,
		bson.M{"_id": installEquipment.BuildingId},
		bson.M{"$set": bson.M{"squareInUse": equipmentType.Square*float64(installEquipment.Amount) + building.SquareInUse}})
	if err != nil {
		return err
	}

	return nil
}

func getEquipmentPosition(equipments *[]Equipment, equipmentTypeId uint) int {
	if equipments == nil {
		return -1
	}
	for i, v := range *equipments {
		if v.EquipmentTypeId == equipmentTypeId {
			return i
		}
	}
	return -1
}

func updateEquipmentAmount(m *mongo.Database, ctx context.Context, building Building, index int, equipmentTypeId uint, amount int) error {
	(*building.Equipment)[index].Amount += amount
	if (*building.Equipment)[index].Amount == 0 {
		update := bson.D{
			{"$pull", bson.D{{"equipment", bson.D{{"equipmentTypeId", equipmentTypeId}}}}},
		}
		_, err := m.Collection("buildings").UpdateOne(ctx, bson.M{"_id": building.Id}, update)
		return err
	}

	update := bson.M{
		"$set": bson.M{
			"equipment.$[id].amount": (*building.Equipment)[index].Amount,
		},
	}
	updateOpts := options.Update().SetArrayFilters(options.ArrayFilters{
		Filters: []interface{}{
			bson.D{{"id.equipmentTypeId", equipmentTypeId}},
		},
	})
	_, err := m.Collection("buildings").UpdateOne(ctx, bson.M{"_id": building.Id}, update, updateOpts)
	return err
}

func addEquipment(m *mongo.Database, ctx context.Context, building *Building, newEquipment Equipment) error {
	var err error
	if (*building).Equipment != nil {
		_, err = m.Collection("buildings").UpdateOne(ctx, bson.M{"_id": (*building).Id},
			bson.M{
				"$push": bson.M{
					"equipment": newEquipment,
				},
			},
		)
	} else {
		var newEquipmentArray []Equipment
		newEquipmentArray = append(newEquipmentArray, newEquipment)
		log.Println(newEquipment, newEquipmentArray)
		_, err = m.Collection("buildings").UpdateOne(ctx, bson.M{"_id": (*building).Id},
			bson.M{
				"$set": bson.M{
					"equipment": &newEquipmentArray,
				},
			},
		)
	}
	return err
}

func GetBuildingsForEquipmentRecount(m *mongo.Database) ([]Building, error) {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(3*time.Second))
	defer cancel()

	var buildings []Building
	cursor, err := m.Collection("buildings").Find(ctx, bson.M{"equipment": bson.M{"$ne": nil}})
	if err != nil {
		return buildings, err
	}

	err = cursor.All(ctx, &buildings)
	return buildings, err
}
