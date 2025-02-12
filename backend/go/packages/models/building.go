package models

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"time"
)

type BuildingStatus string // @name buildingStatus

const (
	ConstructionStatus    BuildingStatus = "Construction"
	ReadyStatus           BuildingStatus = "Ready"
	ProductionStatus      BuildingStatus = "Production"
	ResourcesNeededStatus BuildingStatus = "ResourcesNeeded"
	StorageNeededStatus   BuildingStatus = "StorageNeeded"
)

type Building struct {
	Id              primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	TypeId          uint               `json:"typeId" bson:"typeId"`
	UserId          primitive.ObjectID `json:"userId" bson:"userId"`
	X               int                `json:"x" bson:"x"`
	Y               int                `json:"y" bson:"y"`
	Square          int                `json:"square" bson:"square"`
	Level           int                `json:"level" bson:"level"`
	SquareInUse     float64            `json:"squareInUse" bson:"squareInUse"`
	Status          BuildingStatus     `json:"status" bson:"status"`
	WorkStarted     time.Time          `json:"workStarted" bson:"workStarted"`
	WorkEnd         time.Time          `json:"workEnd" bson:"workEnd"`
	HiringNeeds     int                `json:"hiringNeeds" bson:"hiringNeeds"`
	Salary          float64            `json:"salary" bson:"salary"`
	Workers         int                `json:"workers" bson:"workers"`
	OnStrike        bool               `json:"onStrike" bson:"onStrike"`
	Production      *Production        `json:"production" bson:"production"`
	Goods           *[]Goods           `json:"goods" bson:"goods"`
	Logistics       *Logistics         `json:"logistics" bson:"logistics"`
	Bank            *Bank              `json:"bank" bson:"bank"`
	CreditTerms     *[]CreditTerms     `json:"creditTerms" bson:"creditTerms"`
	Equipment       *[]Equipment       `json:"equipment" bson:"equipment"`
	EquipmentEffect *[]EquipmentEffect `json:"equipmentEffect"  bson:"equipmentEffect"`
} // @name building

type ConstructBuildingPayload struct {
	TypeId uint `json:"typeId" validate:"required"`
	X      int  `json:"x" validate:"required"`
	Y      int  `json:"y" validate:"required"`
	Square int  `json:"square" validate:"required"`
} // @name constructBuildingPayload

func ConstructBuilding(m *mongo.Database, userId primitive.ObjectID, payload ConstructBuildingPayload) error {
	enoughLand, err := CheckEnoughLandForBuilding(m, userId, payload.Square, payload.X, payload.Y)
	if !enoughLand {
		return errors.New("not enough land in this cell")
	}
	if err != nil {
		return err
	}

	buildingType, err := GetBuildingTypeById(m, payload.TypeId)
	if err != nil {
		return err
	}
	if !CheckEnoughMoney(m, userId, buildingType.Cost*float64(payload.Square)) {
		return errors.New("not enough money")
	}

	var requirements []ResourceAmount
	for _, resource := range buildingType.Requirements {
		requirements = append(requirements, ResourceAmount{
			ResourceId: resource.ResourceId,
			Amount:     resource.Amount * float64(payload.Square),
		})
	}

	if !CheckEnoughResourcesAmount(m, userId, payload.X, payload.Y, requirements) {
		return errors.New("not enough resources for construction in cell")
	}

	return CreateBuilding(m, userId, payload, buildingType, requirements)
}

func CreateBuilding(m *mongo.Database, userId primitive.ObjectID, payload ConstructBuildingPayload, buildingType BuildingType, requirements []ResourceAmount) error {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(3*time.Second))
	defer cancel()

	if err := AddMoney(m, userId, (-1)*buildingType.Cost*float64(payload.Square)); err != nil {
		return err
	}

	if err := AddCivilSavings(m, payload.X, payload.Y, buildingType.Cost*float64(payload.Square)); err != nil {
		return err
	}

	for i := range requirements {
		requirements[i].Amount = (-1) * requirements[i].Amount
	}

	if err := AddResourceArray(m, userId, payload.X, payload.Y, requirements); err != nil {
		return err
	}

	now := time.Now()
	end := now.Add(time.Duration(float64(buildingType.BuildTime) * float64(payload.Square)))
	building := Building{
		TypeId:      payload.TypeId,
		UserId:      userId,
		X:           payload.X,
		Y:           payload.Y,
		Square:      payload.Square,
		Level:       1,
		Status:      ConstructionStatus,
		WorkStarted: now,
		WorkEnd:     end,
		HiringNeeds: 0,
		Salary:      0,
		Workers:     0,
		Production:  nil,
		Goods:       nil,
	}

	_, err := m.Collection("buildings").InsertOne(ctx, building)
	return err
}

type FindBuildingParams struct {
	Id             *primitive.ObjectID `json:"id"`
	UserId         *primitive.ObjectID `json:"userId"`
	NickName       *string             `json:"nickName"`
	X              *int                `json:"x"`
	Y              *int                `json:"y"`
	BuildingTypeId *uint               `json:"buildingTypeId"`
	Limit          *int                `json:"limit"`
	OrderField     *string             `json:"orderField"`
	Order          *string             `json:"order"`
	Page           *int                `json:"page"`
} // @name findBuildingParams

type BuildingWithData struct {
	Id              primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty" validate:"required"`
	TypeId          uint               `json:"typeId" validate:"required"`
	UserId          primitive.ObjectID `json:"userId" validate:"required"`
	X               int                `json:"x" validate:"required"`
	Y               int                `json:"y" validate:"required"`
	Square          int                `json:"square" validate:"required"`
	SquareInUse     float64            `json:"squareInUse" validate:"required"`
	Level           int                `json:"level" validate:"required"`
	Status          BuildingStatus     `json:"status" validate:"required"`
	WorkStarted     time.Time          `json:"workStarted"`
	WorkEnd         time.Time          `json:"workEnd"`
	HiringNeeds     int                `json:"hiringNeeds" validate:"required"`
	Salary          float64            `json:"salary" validate:"required"`
	Workers         int                `json:"workers" validate:"required"`
	OnStrike        bool               `json:"onStrike" validate:"required"`
	BuildingType    BuildingType       `json:"buildingType" validate:"required"`
	NickName        string             `json:"nickName" validate:"required"`
	Production      *Production        `json:"production"`
	Goods           *[]Goods           `json:"goods"`
	Logistics       *Logistics         `json:"logistics"`
	Bank            *Bank              `json:"bank"`
	Equipment       *[]Equipment       `json:"equipment"`
	EquipmentEffect *[]EquipmentEffect `json:"equipmentEffect"`
} // @name buildingWithData

func GetBuildings(m *mongo.Database, findBuildingParams FindBuildingParams) ([]BuildingWithData, error) {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(3*time.Second))
	defer cancel()

	filter := bson.D{}
	if findBuildingParams.Id != nil {
		filter = append(filter, bson.E{Key: "buildings._id", Value: *findBuildingParams.Id})
	}
	if findBuildingParams.UserId != nil {
		filter = append(filter, bson.E{Key: "userId", Value: *findBuildingParams.UserId})
	}
	if findBuildingParams.BuildingTypeId != nil {
		filter = append(filter, bson.E{Key: "typeId", Value: *findBuildingParams.BuildingTypeId})
	}
	if findBuildingParams.X != nil {
		filter = append(filter, bson.E{Key: "x", Value: *findBuildingParams.X})
	}
	if findBuildingParams.Y != nil {
		filter = append(filter, bson.E{Key: "y", Value: *findBuildingParams.Y})
	}

	matchStage := bson.D{{"$match", filter}}
	lookupBuildingType := bson.D{{"$lookup", bson.D{
		{"from", "buildingTypes"},
		{"localField", "typeId"},
		{"foreignField", "id"},
		{"as", "buildingType"},
	}}}
	lookupUser := bson.D{{"$lookup", bson.D{
		{"from", "users"},
		{"localField", "userId"},
		{"foreignField", "_id"},
		{"as", "user"},
	}}}

	unwindBuildingType := bson.D{{"$unwind", bson.D{
		{"path", "$buildingType"},
		{"preserveNullAndEmptyArrays", true},
	}}}
	unwindUser := bson.D{{"$unwind", bson.D{
		{"path", "$user"},
		{"preserveNullAndEmptyArrays", true},
	}}}

	projectStage := bson.D{{"$project", bson.D{
		{"typeId", 1},
		{"x", 1},
		{"y", 1},
		{"level", 1},
		{"status", 1},
		{"square", 1},
		{"buildingType.title", 1},
		{"nickName", "$user.nickName"},
	}}}

	sort := bson.D{}

	if findBuildingParams.OrderField != nil {
		if findBuildingParams.Order != nil {
			sort = append(filter, bson.E{Key: *findBuildingParams.OrderField, Value: *findBuildingParams.Order})
		} else {
			sort = append(filter, bson.E{Key: *findBuildingParams.OrderField, Value: 1})
		}
	}

	sortStage := bson.D{}

	if len(sort) != 0 {
		sortStage = bson.D{{"$sort", sort}}
	} else {
		sortStage = bson.D{{"$sort", bson.D{{"_id", -1}}}}
	}

	limit := 20
	if findBuildingParams.Limit != nil {
		limit = *findBuildingParams.Limit
	}
	limitStage := bson.D{{"$limit", limit}}

	skipStage := bson.D{{"$skip", 0}}
	if findBuildingParams.Page != nil {
		skipStage = bson.D{{"$skip", (*findBuildingParams.Page - 1) * limit}}
	}

	pipeline := mongo.Pipeline{matchStage, lookupBuildingType, lookupUser, unwindBuildingType, unwindUser,
		projectStage, sortStage, skipStage, limitStage}

	cursor, err := m.Collection("buildings").Aggregate(ctx, pipeline)
	if err != nil {
		log.Println("Can't get buildings: " + err.Error())
		return nil, err
	}

	var buildings []BuildingWithData
	if err = cursor.All(ctx, &buildings); err != nil {
		log.Println("Can't get buildings: " + err.Error())
	}
	return buildings, nil
}

func GetMyBuildings(m *mongo.Database, userId primitive.ObjectID, buildingId primitive.ObjectID) ([]BuildingWithData, error) {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(3*time.Second))
	defer cancel()

	filter := bson.D{}
	filter = append(filter, bson.E{Key: "userId", Value: userId})
	if buildingId != primitive.NilObjectID {
		filter = append(filter, bson.E{Key: "_id", Value: buildingId})
	}
	matchStage := bson.D{{"$match", filter}}

	lookupBuildingType := bson.D{{"$lookup", bson.D{
		{"from", "buildingTypes"},
		{"localField", "typeId"},
		{"foreignField", "id"},
		{"as", "buildingType"},
	}}}

	unwindBuildingType := bson.D{{"$unwind", bson.D{
		{"path", "$buildingType"},
		{"preserveNullAndEmptyArrays", true},
	}}}

	pipeline := mongo.Pipeline{matchStage, lookupBuildingType, unwindBuildingType}
	cursor, err := m.Collection("buildings").Aggregate(ctx, pipeline)
	if err != nil {
		log.Println("Can't get buildings: " + err.Error())
		return nil, err
	}

	var myBuildings []BuildingWithData
	if err = cursor.All(ctx, &myBuildings); err != nil {
		log.Println(err)
	}
	return myBuildings, nil
}

func GetBuildingById(m *mongo.Database, buildingId primitive.ObjectID) (Building, error) {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(3*time.Second))
	defer cancel()

	var building Building
	err := m.Collection("buildings").FindOne(ctx,
		bson.M{"_id": buildingId}).Decode(&building)
	if err != nil {
		log.Println("Can't get building by Id: " + err.Error())
	}
	return building, err
}

func DestroyBuilding(m *mongo.Database, userId primitive.ObjectID, buildingId primitive.ObjectID) error {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(3*time.Second))
	defer cancel()

	building, err := GetBuildingById(m, buildingId)
	if err != nil {
		log.Println("Can't get building: " + err.Error())
		return err
	}
	if userId != building.UserId && building.UserId != primitive.NilObjectID {
		return errors.New("for attempting to destroy someone else's building, inevitable punishment awaits you")
	}

	_, err = m.Collection("buildings").DeleteOne(ctx, bson.M{"_id": buildingId, "userId": userId})
	if err != nil {
		log.Println("Failed to delete building: " + err.Error())
		return err
	}

	return nil
}

func BuildingStatusUpdate(m *mongo.Database, buildingId primitive.ObjectID, status BuildingStatus) error {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(3*time.Second))
	defer cancel()

	_, err := m.Collection("buildings").UpdateOne(ctx,
		bson.M{"_id": buildingId},
		bson.M{"$set": bson.M{"status": status}})
	return err
}

func BuildingSetWorkStarted(m *mongo.Database, buildingId primitive.ObjectID, timeStart time.Time) error {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(3*time.Second))
	defer cancel()

	_, err := m.Collection("buildings").UpdateOne(ctx,
		bson.M{"_id": buildingId},
		bson.M{"$set": bson.M{"workStarted": timeStart}})
	return err
}

type HiringPayload struct {
	BuildingID  primitive.ObjectID `json:"buildingId" bson:"buildingId"`
	Salary      float64            `json:"salary" bson:"salary"`
	HiringNeeds int                `json:"hiringNeeds" bson:"hiringNeeds"`
} // @name hiringPayload

func SetHiring(m *mongo.Database, userId primitive.ObjectID, payload HiringPayload) error {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(3*time.Second))
	defer cancel()

	building, err := GetBuildingById(m, payload.BuildingID)
	if err != nil {
		return err
	}
	if userId != building.UserId && building.UserId != primitive.NilObjectID {
		return errors.New("this building doesn't belong you")
	}
	buildingType, err := GetBuildingTypeById(m, building.TypeId)
	if err != nil {
		return err
	}
	hiringMax := buildingType.Workers * building.Level * building.Square
	if payload.HiringNeeds > hiringMax {
		return errors.New("hiring needs more that maximum")
	}

	_, err = m.Collection("buildings").UpdateOne(ctx,
		bson.M{"_id": building.Id},
		bson.M{"$set": bson.M{"salary": payload.Salary, "hiringNeeds": payload.HiringNeeds}})
	if err != nil {
		log.Println("Can't update building: " + err.Error())
	}
	return err
}

func GetBuildingsForHiring(m *mongo.Database) ([]Building, error) {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(3*time.Second))
	defer cancel()

	filter := bson.M{"salary": bson.M{"$ne": 0}, "hiringNeeds": bson.M{"$ne": 0}}
	cursor, err := m.Collection("buildings").Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	var buildings []Building
	err = cursor.All(ctx, &buildings)
	return buildings, err
}

func GetAllReadyBuildingByGroup(m *mongo.Database, group string) ([]Building, error) {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(3*time.Second))
	defer cancel()

	var readyLogisticsHubs []Building

	logisticsTypes, err := GetBuildingTypesByBuildingGroup(m, group)
	if err != nil {
		return readyLogisticsHubs, err
	}

	var typeIds []uint
	for _, logisticsType := range logisticsTypes {
		typeIds = append(typeIds, logisticsType.Id)
	}

	filter := bson.M{"status": ReadyStatus, "onStrike": false, "typeId": bson.M{"$in": typeIds}}
	cursor, err := m.Collection("buildings").Find(ctx, filter)
	if err != nil {
		return readyLogisticsHubs, err
	}

	err = cursor.All(ctx, &readyLogisticsHubs)
	return readyLogisticsHubs, err
}

func CheckBuildingCell(m *mongo.Database, buildingId primitive.ObjectID, x int, y int) bool {
	building, err := GetBuildingById(m, buildingId)
	if err != nil {
		log.Println("can't check buildings", buildingId, "cell:", err)
		return false
	}
	return building.X == x && building.Y == y
}

type EmergencyHiringPayload struct {
	BuildingID primitive.ObjectID `json:"buildingId" validate:"required"`
} // @name emergencyHiringPayload

func EmergencyHiring(m *mongo.Database, userId primitive.ObjectID, payload EmergencyHiringPayload) error {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(3*time.Second))
	defer cancel()

	building, err := GetBuildingById(m, payload.BuildingID)
	if err != nil {
		log.Println("can't check buildings", payload.BuildingID, "cell:", err)
		return err
	}
	if userId != building.UserId && building.UserId != primitive.NilObjectID {
		return errors.New("this building doesn't belong you")
	}
	if building.Workers >= building.HiringNeeds {
		return errors.New("impossible")
	}
	cell, err := GetCell(m, building.X, building.Y)
	if err != nil {
		log.Println("can't get cell ", building.X, "x", building.Y, ":", err)
		return err
	}
	// FORMULA Emergence hiring
	price := float64(building.HiringNeeds-building.Workers) * cell.AverageSalary * 10
	if CheckEnoughMoney(m, building.UserId, price) {
		err = AddMoney(m, building.UserId, -1*price)
		if err != nil {
			return err
		}
		_, err := m.Collection("buildings").UpdateOne(ctx,
			bson.M{"_id": building.Id},
			bson.M{"$set": bson.M{"workers": building.HiringNeeds}})
		return err
	} else {
		return errors.New("not enough money")
	}
}
