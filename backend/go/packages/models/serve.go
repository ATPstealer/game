package models

import (
	"backend/packages/cfg"
	"context"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
	"log"
	"math/rand"
	"strconv"
	"time"
)

func DeleteObsoleteTokens(m *mongo.Database) {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(60*time.Second))
	defer cancel()

	now := time.Now()
	_, err := m.Collection("tokens").DeleteMany(ctx, bson.M{
		"$expr": bson.M{
			"$lt": []interface{}{
				bson.M{
					"$add": []interface{}{
						"$createdAt",
						bson.M{
							"$multiply": []interface{}{"$ttl", 1000}, // Multiply by 1000 to convert from milliseconds to seconds
						},
					},
				},
				now,
			},
		},
	})

	if err != nil {
		log.Println(err)
	}
}

func MongoIndex(m *mongo.Database) {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(60*time.Second))
	defer cancel()

	_, err := m.Collection("users").Indexes().CreateOne(ctx,
		mongo.IndexModel{Keys: bson.D{{"nickName", 1}}, Options: options.Index().SetUnique(true)},
		options.CreateIndexes().SetMaxTime(60*time.Second))
	if err != nil {
		log.Fatal(err)
	}
	_, err = m.Collection("users").Indexes().CreateOne(ctx,
		mongo.IndexModel{Keys: bson.D{{"email", 1}}, Options: options.Index().SetUnique(true)},
		options.CreateIndexes().SetMaxTime(60*time.Second))
	if err != nil {
		log.Fatal(err)
	}
}

func Init(m *mongo.Database, config cfg.Vars) {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(60*time.Second))
	defer cancel()

	if config.Init {
		count, err := m.Collection("settings").CountDocuments(ctx, bson.M{})
		if err != nil {
			log.Fatal(err)
		}
		if count == 0 {
			loadSettings(m)
		}
		count, err = m.Collection("cells").CountDocuments(ctx, bson.M{})
		if err != nil {
			log.Fatal(err)
		}
		if count == 0 {
			generateMap(m)
		}
		importFromDataSheet(m)
	}
}

func loadSettings(m *mongo.Database) {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(60*time.Second))
	defer cancel()

	settings := []Settings{
		{Key: "mapMinX", Value: -2},
		{Key: "mapMaxX", Value: 2},
		{Key: "mapMinY", Value: -2},
		{Key: "mapMaxY", Value: 2},
		{Key: "interestRate", Value: 0.5},
	}
	collection := m.Collection("settings")
	for _, setting := range settings {
		filter := bson.D{{"key", setting.Key}}
		update := bson.D{{"$set", setting}}
		_, err := collection.UpdateOne(ctx, filter, update, options.Update().SetUpsert(true))
		if err != nil {
			log.Fatal("Error while updating setting:", err)
		}
	}
}

func generateMap(m *mongo.Database) {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(60*time.Second))
	defer cancel()

	settings, _ := GetSettings(m)
	var newMapCell Cell
	for y := int(settings["mapMinY"]); y <= int(settings["mapMaxY"]); y++ {
		for x := int(settings["mapMinX"]); x <= int(settings["mapMaxX"]); x++ {
			r := rand.New(rand.NewSource(time.Now().UnixNano()))
			newMapCell = Cell{
				CellName:         strconv.Itoa(x) + "x" + strconv.Itoa(y),
				X:                x,
				Y:                y,
				SurfaceImagePath: "/map/grass/land" + strconv.Itoa(x) + "x" + strconv.Itoa(y) + ".png",
				Square:           int(rand.Float64() * 10000),
				Pollution:        r.Float64() * 1000,
				Population:       r.Float64() * 10000,
				CivilSavings:     r.Float64() * 1000000,
				SpendRate:        r.Float64() * 0.1,
				Education:        r.Float64() * 10,
				Crime:            r.Float64() * 10,
				Medicine:         r.Float64() * 10,
			}
			collection := m.Collection("cells")
			filter := bson.D{{"X", newMapCell.X}, {"Y", newMapCell.Y}}
			update := bson.D{{"$set", newMapCell}}
			_, err := collection.UpdateOne(ctx, filter, update, options.Update().SetUpsert(true))
			if err != nil {
				log.Fatal("Error while updating cells:", err)
			}
		}
	}
}

func importFromDataSheet(m *mongo.Database) {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(60*time.Second))
	defer cancel()

	conf, err := google.JWTConfigFromJSON([]byte(cfg.Config.GoogleAPI), "https://www.googleapis.com/auth/spreadsheets.readonly")
	if err != nil {
		log.Println("Can't get data from Google Sheet: ", err)
	}
	srv, err := sheets.NewService(ctx, option.WithHTTPClient(conf.Client(ctx)))
	if err != nil {
		log.Println("Can't get data from Google Sheet: ", err)
	}
	sheetsList, err := srv.Spreadsheets.Get(cfg.Config.GoogleSheetID).Do()
	if err != nil {
		log.Println("Unable to retrieve data from sheet: ", err)
	}
	for _, sheet := range sheetsList.Sheets {
		rows, err := srv.Spreadsheets.Values.Get(cfg.Config.GoogleSheetID, sheet.Properties.Title).Do()
		if err != nil {
			log.Fatalf("Unable to retrieve data from sheet: %v", err)
		}
		if len(rows.Values) == 0 {
			fmt.Println("No data found in sheet: ", sheet.Properties.Title)
		} else {
			if err := importDataInTable(m, sheet.Properties.Title, rows.Values); err != nil {
				fmt.Println("Func importDataInTable failed")
			}
		}
	}
}

func importDataInTable(m *mongo.Database, tableName string, rows [][]interface{}) error {
	if tableName == "buildingTypes" {
		if err := buildingTypesImport(m, rows); err != nil {
			fmt.Println("Import table building_types failed")
		}
	}
	if tableName == "resourceTypes" {
		if err := resourceTypesImport(m, rows); err != nil {
			fmt.Println("Import table resource_types failed")
		}
	}
	if tableName == "blueprints" {
		if err := productionBlueprintsImport(m, rows); err != nil {
			fmt.Println("Import table production_blueprints failed")
		}
	}
	if tableName == "equipmentTypes" {
		if err := equipmentTypesImport(m, rows); err != nil {
			fmt.Println("Import table production_blueprints failed")
		}
	}
	return nil
}

func buildingTypesImport(m *mongo.Database, rows [][]interface{}) error {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(60*time.Second))
	defer cancel()

	for _, row := range rows[1:] {
		id, err := strconv.ParseInt(row[1].(string), 10, 64)
		if err != nil {
			log.Println("Can't get time.Duration from Google sheet id field: ", err)
		}
		cost, err := strconv.ParseFloat(row[4].(string), 32)
		if err != nil {
			log.Println("Can't get float from Google sheet Cost field: ", err)
		}
		buildTime, err := strconv.ParseInt(row[6].(string), 10, 64)
		if err != nil {
			log.Println("Can't get time.Duration from Google sheet BuildTime field: ", err)
		}
		capacity, err := strconv.ParseInt(row[9].(string), 10, 32)
		if err != nil {
			log.Println("Can't get Int from Google sheet Capacity field: ", err)
		}
		workers, err := strconv.ParseInt(row[10].(string), 10, 32)
		if err != nil {
			log.Println("Can't get Int from Google sheet Workers field: ", err)
		}

		buildingTypeMongo := BuildingType{
			Id:               uint(id),
			Title:            row[2].(string),
			Description:      row[3].(string),
			Cost:             float64(cost),
			Requirements:     row[5].(string),
			BuildTime:        time.Second * time.Duration(buildTime),
			BuildingGroup:    row[7].(string),
			BuildingSubGroup: row[8].(string),
			Capacity:         float64(capacity),
			Workers:          int(workers),
		}

		collection := m.Collection("buildingTypes")
		filter := bson.D{{"id", buildingTypeMongo.Id}}
		update := bson.D{{"$set", buildingTypeMongo}}
		_, err = collection.UpdateOne(ctx, filter, update, options.Update().SetUpsert(true))
		if err != nil {
			log.Fatal("Error while updating buildingTypes:", err)
		}
	}
	return nil
}

func resourceTypesImport(m *mongo.Database, rows [][]interface{}) error {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(60*time.Second))
	defer cancel()

	for _, row := range rows[1:] {
		id, err := strconv.ParseInt(row[1].(string), 10, 64)
		if err != nil {
			log.Println("Can't get time.Duration from Google sheet Id field: ", err)
		}
		volume, err := strconv.ParseFloat(row[3].(string), 32)
		if err != nil {
			log.Println("Can't get float from Google sheet Volume field: ", err)
		}
		weight, err := strconv.ParseFloat(row[4].(string), 32)
		if err != nil {
			log.Println("Can't get float from Google sheet Weight field: ", err)
		}
		demand, err := strconv.ParseFloat(row[5].(string), 32)
		if err != nil {
			log.Println("Can't get float from Google sheet Demand field: ", err)
		}

		resourceTypeMongo := ResourceType{
			Id:         uint(id),
			Name:       row[2].(string),
			Volume:     float64(volume),
			Weight:     float64(weight),
			Demand:     float64(demand),
			StoreGroup: row[6].(string),
		}

		collection := m.Collection("resourceTypes")
		filter := bson.D{{"id", resourceTypeMongo.Id}}
		update := bson.D{{"$set", resourceTypeMongo}}
		_, err = collection.UpdateOne(ctx, filter, update, options.Update().SetUpsert(true))
		if err != nil {
			log.Fatal("Error while updating resourceTypes:", err)
		}
	}
	return nil
}

func productionBlueprintsImport(m *mongo.Database, rows [][]interface{}) error {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(60*time.Second))
	defer cancel()

	for _, row := range rows[1:] {
		id, err := strconv.ParseInt(row[1].(string), 10, 64)
		if err != nil {
			log.Println("Can't get time.Duration from Google sheet Id field: ", err)
		}
		var producedResources, usedResources []ResourceAmount
		if err := json.Unmarshal([]byte(row[2].(string)), &producedResources); err != nil {
			log.Println("Error while unmarshalling ProducedResources:", err)
			return err
		}

		if err := json.Unmarshal([]byte(row[3].(string)), &usedResources); err != nil {
			log.Println("Error while unmarshalling UsedResources:", err)
			return err
		}
		producedInID, err := strconv.ParseUint(row[4].(string), 10, 32)
		if err != nil {
			log.Println("Can't get UInt from Google sheet ProducedInId field: ", err)
		}
		productionTime, err := strconv.ParseInt(row[5].(string), 10, 32)
		if err != nil {
			log.Println("Can't get UInt from Google sheet ProductionTime field: ", err)
		}

		blueprintMongo := Blueprint{
			Id:                uint(id),
			Name:              row[6].(string),
			ProducedResources: producedResources,
			UsedResources:     usedResources,
			ProducedInId:      uint(producedInID),
			ProductionTime:    time.Second * time.Duration(productionTime),
		}

		collection := m.Collection("blueprints")
		filter := bson.D{{"id", blueprintMongo.Id}}
		update := bson.D{{"$set", blueprintMongo}}
		_, err = collection.UpdateOne(ctx, filter, update, options.Update().SetUpsert(true))
		if err != nil {
			log.Fatal("Error while updating resourceTypes:", err)
		}
	}
	return nil
}

func equipmentTypesImport(m *mongo.Database, rows [][]interface{}) error {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(60*time.Second))
	defer cancel()

	for _, row := range rows[1:] {
		id, err := strconv.ParseUint(row[1].(string), 10, 64)
		if err != nil {
			log.Println("Can't get Uint from Google sheet Id field: ", err)
		}
		resourceTypeId, err := strconv.ParseUint(row[3].(string), 10, 64)
		if err != nil {
			log.Println("Can't get Uint from Google sheet ResourceTypeId field: ", err)
		}
		durability, err := strconv.ParseInt(row[4].(string), 10, 32)
		if err != nil {
			log.Println("Can't get Int from Google sheet Durability field: ", err)
		}
		blueprintIds := make([]uint, 0)
		if err := json.Unmarshal([]byte(row[5].(string)), &blueprintIds); err != nil {
			log.Println("Error while unmarshalling BlueprintIds:", err)
			return err
		}
		effectId, err := strconv.ParseUint(row[6].(string), 10, 64)
		if err != nil {
			log.Println("Can't get Uint from Google sheet EffectId field: ", err)
		}
		value, err := strconv.ParseFloat(row[7].(string), 64)
		if err != nil {
			log.Println("Can't get Float64 from Google sheet Value field: ", err)
		}
		square, err := strconv.ParseFloat(row[8].(string), 64)
		if err != nil {
			log.Println("Can't get Float64 from Google sheet Square field: ", err)
		}

		equipmentTypeMongo := EquipmentType{
			Id:             uint(id),
			Name:           row[2].(string),
			ResourceTypeId: uint(resourceTypeId),
			Durability:     int(durability),
			BlueprintIds:   blueprintIds,
			EffectId:       uint(effectId),
			Value:          value,
			Square:         square,
		}

		collection := m.Collection("equipmentTypes")
		filter := bson.D{{"id", equipmentTypeMongo.Id}}
		update := bson.D{{"$set", equipmentTypeMongo}}
		_, err = collection.UpdateOne(ctx, filter, update, options.Update().SetUpsert(true))
		if err != nil {
			log.Fatal("Error while updating equipmentType:", err)
		}
	}
	return nil
}
