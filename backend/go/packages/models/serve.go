package models

import (
	"backend/packages/cfg"
	"context"
	"fmt"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
	"gorm.io/gorm"
	"log"
	"reflect"
	"strconv"
	"time"
)

func AutoMigrateModel(db *gorm.DB) {
	models := []interface{}{&User{}, &Token{}, &Cell{}, &LandLord{}, &Building{},
		&BuildingType{}, &Resource{}, &ResourceType{}, &Logistic{}, &Blueprint{},
		&Storage{}, &Order{}, &Settings{}, &StoreGoods{}, &EvolutionPrice{}}

	for _, model := range models {
		if err := db.AutoMigrate(model); err != nil {
			log.Println(fmt.Sprintf("Can't migrate model: %v Error: %s", reflect.TypeOf(model), err.Error()))
		}
	}
}

func DeleteObsoleteTokens(db *gorm.DB) {
	err := db.Unscoped().Where("deleted_at IS NOT NULL OR " +
		"NOW() > (created_at + INTERVAL TTL SECOND)").Delete(&Token{}).Error
	if err != nil {
		fmt.Println("Error occurred:", err)
	} else {
		fmt.Println("Deletion was successful.")
	}
}

func Init(db *gorm.DB, config cfg.Vars) {
	if config.Init {
		var mapCells []Cell
		db.Find(&mapCells)
		if len(mapCells) == 0 {
			generateMap(db)
		}
		importFromDataSheet(db)
		var settings []Settings
		db.Find(&settings)
		if len(settings) == 0 {
			loadSettings(db)
		}
	}
}

func loadSettings(db *gorm.DB) {
	settings := [4]Settings{
		{
			Key:   "mapMinX",
			Value: -2,
		},
		{
			Key:   "mapMaxX",
			Value: 2,
		},
		{
			Key:   "mapMinY",
			Value: -2,
		},
		{
			Key:   "mapMaxY",
			Value: 2,
		},
	}
	db.Save(&settings)
}

// GenerateMap Use it if necessary
func generateMap(db *gorm.DB) {
	settings := GetSettingsMap(db)
	var newMapCell Cell
	for y := int(settings["mapMinY"]); y <= int(settings["mapMaxY"]); y++ {
		for x := int(settings["mapMinX"]); x <= int(settings["mapMaxX"]); x++ {
			newMapCell = Cell{
				CellName:         strconv.Itoa(x) + "x" + strconv.Itoa(y),
				X:                x,
				Y:                y,
				SurfaceImagePath: "/map/grass/land" + strconv.Itoa(x) + "x" + strconv.Itoa(y) + ".png",
				Square:           10000,
				Pollution:        1000,
				Population:       10000,
				CivilSavings:     1000000,
				SpendRate:        0.1,
				Education:        10,
				Crime:            10,
				Medicine:         10,
				ElementarySchool: 10,
				HigherSchool:     10,
			}
			db.Create(&newMapCell)
		}
	}
}

func importFromDataSheet(db *gorm.DB) { // ImportDataFromGoogleSheets() {
	ctx := context.Background()
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
			if err := importDataInTable(db, sheet.Properties.Title, rows.Values); err != nil {
				fmt.Println("Func importDataInTable failed")
			}
		}
	}
}

func importDataInTable(db *gorm.DB, tableName string, rows [][]interface{}) error {
	if tableName == "building_types" {
		if err := buildingTypesImport(db, rows); err != nil {
			fmt.Println("Import table building_types failed")
		}
	}
	if tableName == "resource_types" {
		if err := resourceTypesImport(db, rows); err != nil {
			fmt.Println("Import table resource_types failed")
		}
	}
	if tableName == "blueprints" {
		if err := productionBlueprintsImport(db, rows); err != nil {
			fmt.Println("Import table production_blueprints failed")
		}
	}
	return nil
}

func buildingTypesImport(db *gorm.DB, rows [][]interface{}) error {
	for i, row := range rows[1:] {
		var buildingType BuildingType
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
			log.Println("Can't get UInt from Google sheet Capacity field: ", err)
		}

		db.Unscoped().Model(&BuildingType{}).Where("id = ?", i+1).First(&buildingType)
		if buildingType == (BuildingType{}) {
			buildingType = BuildingType{
				Title:            row[2].(string),
				Description:      row[3].(string),
				Cost:             float32(cost),
				Requirements:     row[5].(string),
				BuildTime:        time.Second * time.Duration(buildTime),
				BuildingGroup:    row[7].(string),
				BuildingSubGroup: row[8].(string),
				Capacity:         float32(capacity),
			}
			db.Create(&buildingType)
		} else {
			if row[0] == "TRUE" {
				db.Delete(&buildingType)
			} else {
				db.Unscoped().Model(&BuildingType{}).Where("id = ?", i+1).Update("DeletedAt", nil)
				buildingType = BuildingType{}
				db.Model(&BuildingType{}).Where("id = ?", i+1).First(&buildingType)
				buildingType.Title = row[2].(string)
				buildingType.Description = row[3].(string)
				buildingType.Cost = float32(cost)
				buildingType.Requirements = row[5].(string)
				buildingType.BuildTime = time.Second * time.Duration(buildTime)
				buildingType.BuildingGroup = row[7].(string)
				buildingType.BuildingSubGroup = row[8].(string)
				buildingType.Capacity = float32(capacity)
				db.Save(&buildingType)
			}
		}
	}
	return nil
}

func resourceTypesImport(db *gorm.DB, rows [][]interface{}) error {
	for i, row := range rows[1:] {
		var resourceType ResourceType
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

		db.Unscoped().Model(&ResourceType{}).Where("id = ?", i+1).First(&resourceType)
		if resourceType == (ResourceType{}) {
			resourceType = ResourceType{
				Name:       row[2].(string),
				Volume:     float32(volume),
				Weight:     float32(weight),
				Demand:     float32(demand),
				StoreGroup: row[6].(string),
			}
			db.Create(&resourceType)
		} else {
			if row[0] == "TRUE" {
				db.Delete(&resourceType)
			} else {
				db.Unscoped().Model(&ResourceType{}).Where("id = ?", i+1).Update("DeletedAt", nil)
				resourceType = ResourceType{}
				db.Model(&ResourceType{}).Where("id = ?", i+1).First(&resourceType)
				resourceType.Name = row[2].(string)
				resourceType.Volume = float32(volume)
				resourceType.Weight = float32(weight)
				resourceType.Demand = float32(demand)
				resourceType.StoreGroup = row[6].(string)
				db.Save(&resourceType)
			}
		}
	}
	return nil
}

func productionBlueprintsImport(db *gorm.DB, rows [][]interface{}) error {
	for i, row := range rows[1:] {
		var blueprint Blueprint
		producedInID, err := strconv.ParseUint(row[4].(string), 10, 32)
		if err != nil {
			log.Println("Can't get UInt from Google sheet ProducedInID field: ", err)
		}
		productionTime, err := strconv.ParseInt(row[5].(string), 10, 32)
		if err != nil {
			log.Println("Can't get UInt from Google sheet ProductionTime field: ", err)
		}

		db.Unscoped().Model(&Blueprint{}).Where("id = ?", i+1).First(&blueprint)
		if blueprint == (Blueprint{}) {
			blueprint = Blueprint{
				Name:              row[6].(string),
				ProducedResources: row[2].(string),
				UsedResources:     row[3].(string),
				ProducedInID:      uint(producedInID),
				ProductionTime:    time.Second * time.Duration(productionTime),
			}
			db.Create(&blueprint)
		} else {
			if row[0] == "TRUE" {
				db.Delete(&blueprint)
			} else {
				db.Unscoped().Model(&Blueprint{}).Where("id = ?", i+1).Update("DeletedAt", nil)
				blueprint = Blueprint{}
				db.Model(&Blueprint{}).Where("id = ?", i+1).First(&blueprint)
				blueprint.Name = row[6].(string)
				blueprint.ProducedResources = row[2].(string)
				blueprint.UsedResources = row[3].(string)
				blueprint.ProducedInID = uint(producedInID)
				blueprint.ProductionTime = time.Second * time.Duration(productionTime)
				db.Save(&blueprint)
			}
		}
	}
	return nil
}
