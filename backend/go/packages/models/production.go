package models

import (
	"errors"
	"gorm.io/gorm"
	"log"
	"time"
)

type Production struct {
	gorm.Model
	BuildingID  uint       `json:"buildingId"`
	BlueprintID uint       `json:"blueprintId"`
	WorkStarted *time.Time `json:"workStarted"`
	WorkEnd     *time.Time `json:"workEnd"`
}

type StartWorkPayload struct {
	BuildingID  uint
	BlueprintID uint
	Duration    time.Duration
}

func StartWork(db *gorm.DB, userID uint, payload StartWorkPayload) error {
	building, err := GetBuildingByID(db, payload.BuildingID)
	if err != nil {
		log.Println("Can't find buildings: " + err.Error())
		return err
	}
	if building.Status != ReadyStatus {
		return errors.New("Building not ready. Status is " + string(building.Status))
	}
	if building.UserID != userID {
		err := errors.New("this building don't belong you")
		log.Println(err)
		return err
	}
	blueprintResult, err := GetBlueprintByID(db, payload.BlueprintID)
	if blueprintResult.ID == 0 {
		err := errors.New("invalid blueprint")
		log.Println(err)
		return err
	}
	if blueprintResult.ProducedInID != building.TypeID {
		err := errors.New("can't product it here")
		log.Println(err)
		return err
	}

	log.Println(building.WorkStarted)
	now := time.Now()
	end := now.Add(payload.Duration)
	building.WorkStarted = &now
	building.WorkEnd = &end
	building.Status = ProductionStatus
	db.Save(&building)

	db.Save(&Production{
		BuildingID:  payload.BuildingID,
		BlueprintID: payload.BlueprintID,
		WorkStarted: &now,
		WorkEnd:     &end,
	})

	return nil
}
