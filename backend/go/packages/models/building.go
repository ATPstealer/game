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

type BuildingStatus string

const (
	ConstructionStatus    BuildingStatus = "Construction"
	ReadyStatus           BuildingStatus = "Ready"
	ProductionStatus      BuildingStatus = "Production"
	ResourcesNeededStatus BuildingStatus = "ResourcesNeeded"
	StorageNeededStatus   BuildingStatus = "StorageNeeded"
)

type Building struct {
	Id          primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	TypeId      uint               `json:"typeId" bson:"typeId"`
	UserId      primitive.ObjectID `json:"userId" bson:"userId"`
	X           int                `json:"x" bson:"x"`
	Y           int                `json:"y" bson:"y"`
	Square      int                `json:"square" bson:"square"`
	Level       int                `json:"level" bson:"level"`
	Status      BuildingStatus     `json:"status" bson:"status"`
	WorkStarted time.Time          `json:"workStarted" bson:"workStarted"`
	WorkEnd     time.Time          `json:"workEnd" bson:"workEnd"`
	HiringNeeds int                `json:"hiringNeeds" bson:"hiringNeeds"`
	Salary      float64            `json:"salary" bson:"salary"`
	Workers     int                `json:"workers" bson:"workers"`
	OnStrike    bool               `json:"onStrike" bson:"onStrike"`
	Production  *Production        `json:"production" bson:"production"`
	Goods       *[]Goods           `json:"goods" bson:"goods"`
	Equipment   *[]Equipment       `json:"equipment" bson:"equipment"`
}

type Production struct {
	BlueprintId uint `json:"blueprintId" bson:"blueprintId"`
}

type Goods struct {
	ResourceTypeId uint             `json:"resourceTypeId" bson:"resourceTypeId"`
	Price          float64          `json:"price" bson:"price"`
	SellSum        int              `json:"sellSum" bson:"sellSum"`
	Revenue        float64          `json:"revenue" bson:"revenue"`
	SellStarted    time.Time        `json:"sellStarted" bson:"sellStarted"`
	Status         StoreGoodsStatus `json:"status" bson:"status"`
}

type Equipment struct {
	EquipmentTypeId uint `json:"equipmentTypeId" bson:"equipmentTypeId"`
	Amount          int  `json:"amount" bson:"amount"`
	Durability      int  `json:"durability" bson:"durability"`
}

type ConstructBuildingPayload struct {
	TypeId uint `json:"typeId"`
	X      int  `json:"x"`
	Y      int  `json:"y"`
	Square int  `json:"square"`
}

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
	return CreateBuilding(m, userId, payload, buildingType)
}

func CreateBuilding(m *mongo.Database, userId primitive.ObjectID, payload ConstructBuildingPayload, buildingType BuildingType) error {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(3*time.Second))
	defer cancel()

	if err := AddMoney(m, userId, (-1)*buildingType.Cost*float64(payload.Square)); err != nil {
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
}

type BuildingWithData struct {
	Id           primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	TypeId       uint               `json:"typeId"`
	UserId       primitive.ObjectID `json:"userId"`
	X            int                `json:"x"`
	Y            int                `json:"y"`
	Square       int                `json:"square"`
	Level        int                `json:"level"`
	Status       BuildingStatus     `json:"status"`
	WorkStarted  time.Time          `json:"workStarted"`
	WorkEnd      time.Time          `json:"workEnd"`
	HiringNeeds  int                `json:"hiringNeeds"`
	Salary       float64            `json:"salary"`
	Workers      int                `json:"workers"`
	OnStrike     bool               `json:"onStrike"`
	BuildingType BuildingType       `json:"buildingType"`
	NickName     string             `json:"nickName"`
	Production   *Production        `json:"production"`
	Goods        *[]Goods           `json:"goods"`
}

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

type HiringPayload struct {
	BuildingID  primitive.ObjectID `json:"buildingId" bson:"buildingId"`
	Salary      float64            `json:"salary" bson:"salary"`
	HiringNeeds int                `json:"hiringNeeds" bson:"hiringNeeds"`
}

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

func GetAllReadyStorages(m *mongo.Database) ([]Building, error) {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(3*time.Second))
	defer cancel()

	var readyStorages []Building

	filter := bson.M{"status": ReadyStatus, "onStrike": false, "typeId": 1}
	cursor, err := m.Collection("buildings").Find(ctx, filter)
	if err != nil {
		return readyStorages, err
	}

	err = cursor.All(ctx, &readyStorages)
	return readyStorages, err
}

func BuildingStatusUpdate(m *mongo.Database, buildingId primitive.ObjectID, status BuildingStatus) error {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(3*time.Second))
	defer cancel()

	_, err := m.Collection("buildings").UpdateOne(ctx,
		bson.M{"_id": buildingId},
		bson.M{"$set": bson.M{"status": status}})
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

func GetProduction(m *mongo.Database) ([]BuildingWithData, error) {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(3*time.Second))
	defer cancel()

	filter := bson.D{
		{"workEnd", bson.D{{"$gt", time.Now()}}},
		{"production", bson.D{{"$ne", nil}}},
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
		log.Println("Can't get productions: " + err.Error())
		return nil, err
	}

	var buildingWithData []BuildingWithData
	if err = cursor.All(ctx, &buildingWithData); err != nil {
		log.Println(err)
	}
	return buildingWithData, nil
}

func BuildingSetWorkStarted(m *mongo.Database, buildingId primitive.ObjectID, timeStart time.Time) error {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(3*time.Second))
	defer cancel()

	_, err := m.Collection("buildings").UpdateOne(ctx,
		bson.M{"_id": buildingId},
		bson.M{"$set": bson.M{"workStarted": timeStart}})
	return err
}

type StartWorkPayload struct {
	BuildingId  primitive.ObjectID
	BlueprintId uint
	Duration    time.Duration
}

func StartWork(m *mongo.Database, userId primitive.ObjectID, payload StartWorkPayload) error {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(3*time.Second))
	defer cancel()

	building, err := GetBuildingById(m, payload.BuildingId)
	if err != nil {
		log.Println("Can't find buildings: " + err.Error())
		return err
	}
	if building.Status != ReadyStatus {
		return errors.New("building busy")
	}
	if building.UserId != userId {
		err := errors.New("this building don't belong you")
		log.Println(err)
		return err
	}
	blueprintResult, err := GetBlueprintById(m, payload.BlueprintId)
	if err != nil {
		log.Println("invalid blueprint" + err.Error())
		return err
	}
	if blueprintResult.ProducedInId != building.TypeId {
		err := errors.New("can't product it here")
		return err
	}

	now := time.Now()
	end := now.Add(payload.Duration)

	production := Production{
		BlueprintId: payload.BlueprintId,
	}

	_, err = m.Collection("buildings").UpdateOne(ctx, bson.M{"_id": building.Id}, bson.M{
		"$set": bson.M{
			"status":      ProductionStatus,
			"workStarted": now,
			"workEnd":     end,
			"production":  &production,
		},
	})
	if err != nil {
		log.Println("Failed to update building: " + err.Error())
		return err
	}
	return nil
}

func StopWork(m *mongo.Database, userId primitive.ObjectID, payload StartWorkPayload) error {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(3*time.Second))
	defer cancel()

	building, err := GetBuildingById(m, payload.BuildingId)
	if err != nil {
		log.Println("Can't find buildings: " + err.Error())
		return err
	}

	if building.UserId != userId {
		err := errors.New("this building don't belong you")
		log.Println(err)
		return err
	}

	_, err = m.Collection("buildings").UpdateOne(ctx, bson.M{"_id": building.Id}, bson.M{
		"$set": bson.M{
			"status":      ReadyStatus,
			"production":  nil,
			"workStarted": nil,
			"workEnd":     nil,
		},
	})

	return err
}

func GetBuildingsStores(m *mongo.Database) ([]BuildingWithData, error) {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(3*time.Second))
	defer cancel()

	filter := bson.D{
		{"workEnd", bson.D{{"$gt", time.Now()}}},
		{"goods", bson.D{{"$ne", nil}}},
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
		log.Println("Can't get productions: " + err.Error())
		return nil, err
	}

	var buildingWithData []BuildingWithData
	if err = cursor.All(ctx, &buildingWithData); err != nil {
		log.Println(err)
	}
	return buildingWithData, nil
}

func BuildingGoodsStatusUpdate(m *mongo.Database, buildingId primitive.ObjectID, resourceTypeId uint, status StoreGoodsStatus) error {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(3*time.Second))
	defer cancel()

	filter := bson.D{{"_id", buildingId}}
	update := bson.D{{"$set", bson.D{{"goods.$[elem].status", status}}}}
	updateOpts := options.Update().SetArrayFilters(options.ArrayFilters{
		Filters: []interface{}{
			bson.D{{"elem.resourceTypeId", resourceTypeId}},
		},
	})

	_, err := m.Collection("buildings").UpdateOne(ctx, filter, update, updateOpts)
	return err
}

func BuildingSetSellStarted(m *mongo.Database, buildingId primitive.ObjectID, resourceTypeId uint, timeStart time.Time) error {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(3*time.Second))
	defer cancel()

	filter := bson.D{{"_id", buildingId}}
	update := bson.D{{"$set", bson.D{{"goods.$[elem].sellStarted", timeStart}}}}
	updateOpts := options.Update().SetArrayFilters(options.ArrayFilters{
		Filters: []interface{}{
			bson.D{{"elem.resourceTypeId", resourceTypeId}},
		},
	})

	_, err := m.Collection("buildings").UpdateOne(ctx, filter, update, updateOpts)
	return err
}

func BuildingGoodsStatsUpdate(m *mongo.Database, buildingId primitive.ObjectID, resourceTypeId uint, sellSum int, revenue float64) error {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(3*time.Second))
	defer cancel()

	filter := bson.D{{"_id", buildingId}}
	update := bson.D{{"$set",
		bson.D{
			{"goods.$[elem].sellSum", sellSum},
			{"goods.$[elem].revenue", revenue},
		},
	}}
	updateOpts := options.Update().SetArrayFilters(options.ArrayFilters{
		Filters: []interface{}{
			bson.D{{"elem.resourceTypeId", resourceTypeId}},
		},
	})

	_, err := m.Collection("buildings").UpdateOne(ctx, filter, update, updateOpts)
	return err
}

type InstallEquipmentPayload struct {
	BuildingId      primitive.ObjectID `json:"buildingId"`
	EquipmentTypeId uint               `json:"equipmentTypeId"`
	Amount          int                `json:"amount"`
}

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
		return updateEquipmentAmount(m, ctx, building, index, equipmentType.Id, installEquipment.Amount)
	} else {
		newEquipment := Equipment{EquipmentTypeId: installEquipment.EquipmentTypeId, Amount: installEquipment.Amount, Durability: equipmentType.Durability}
		_, err = m.Collection("buildings").UpdateOne(ctx, bson.M{"_id": installEquipment.BuildingId},
			bson.M{
				"$push": bson.M{
					"equipment": newEquipment,
				},
			},
		)
	}
	return err

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
