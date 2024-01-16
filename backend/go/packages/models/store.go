package models

import (
	"errors"
	"gorm.io/gorm"
	"time"
)

type StoreGoodsStatus string

const (
	Selling              StoreGoodsStatus = "Selling"
	DemandSatisfied      StoreGoodsStatus = "DemandSatisfied"
	HighPrice            StoreGoodsStatus = "HighPrice"
	NotEnoughMinerals    StoreGoodsStatus = "NotEnoughMinerals"
	SpendingLimitReached StoreGoodsStatus = "SpendingLimitReached"
	CapacityReached      StoreGoodsStatus = "CapacityReached"
)

type StoreGoods struct {
	gorm.Model
	BuildingID     uint
	ResourceTypeID uint
	Price          float64
	SellSum        int
	Revenue        float64
	SellStarted    *time.Time
	Status         StoreGoodsStatus
}

func SetStoreGoods(db *gorm.DB, userID uint, buildingID uint, resourceTypeID uint, price float64) error {
	building, err := GetBuildingByID(db, buildingID)
	if err != nil {
		return err
	}
	if building.UserID != userID {
		return errors.New("this building don't belong you")
	}

	buildingType, err := GetBuildingTypeByID(db, building.TypeID)
	if err != nil {
		return err
	}
	if buildingType.BuildingGroup != "Store" {
		return errors.New("this is not a store")
	}

	resourceType, err := GetResourceTypesByID(db, resourceTypeID)
	if err != nil {
		return err
	}
	if resourceType.StoreGroup != buildingType.BuildingSubGroup {
		return errors.New("can't sell here")
	}

	now := time.Now()
	var storeGoods StoreGoods
	db.Where(StoreGoods{BuildingID: buildingID, ResourceTypeID: resourceTypeID}).
		Assign(StoreGoods{SellStarted: &now, Status: Selling}).
		FirstOrCreate(&storeGoods)
	storeGoods.Price = price
	result := db.Save(&storeGoods)
	return result.Error
}

type StoreGoodsShortResult struct {
	ID             uint             `json:"id" gorm:"primary_key"`
	BuildingID     uint             `json:"buildingId"`
	ResourceTypeID uint             `json:"resourceTypeId"`
	Price          float64          `json:"price"`
	SellSum        int              `json:"sellSum"`
	Revenue        float64          `json:"revenue"`
	SellStarted    *time.Time       `json:"sellStarted"`
	Status         StoreGoodsStatus `json:"status"`
}

func GetStoreGoods(db *gorm.DB, buildingID uint) ([]StoreGoodsShortResult, error) {
	var storeGoodsShortResult []StoreGoodsShortResult
	res := db.Model(&StoreGoods{}).Where("building_id = ?", buildingID).Find(&storeGoodsShortResult)
	if res.Error != nil {
		return storeGoodsShortResult, res.Error
	}
	return storeGoodsShortResult, nil
}

type StoreGoodsResult struct {
	ID             uint             `json:"id" gorm:"primary_key"`
	BuildingID     uint             `json:"buildingId"`
	ResourceTypeID uint             `json:"resourceTypeId"`
	Price          float64          `json:"price"`
	SellSum        int              `json:"sellSum"`
	Revenue        float64          `json:"revenue"`
	SellStarted    *time.Time       `json:"sellStarted"`
	Status         StoreGoodsStatus `json:"status"`
	X              int              `json:"x"`
	Y              int              `json:"y"`
	Square         int              `json:"square"`
	Level          int              `json:"level"`
	Capacity       int              `json:"capacity"`
	UserID         uint             `json:"userID"`
	OnStrike       bool             `json:"onStrike"`
	Workers        int              `json:"workers"`
	MaxWorkers     int              `json:"maxWorkers"`
}

func GetAllStoreGoods(db *gorm.DB) ([]StoreGoodsResult, error) {
	var storeGoodsResult []StoreGoodsResult
	res := db.Model(&StoreGoods{}).
		Select("store_goods.id", "building_id", "resource_type_id", "price", "sell_sum", "revenue",
			"sell_started", "store_goods.status", "x", "y", "square", "level", "building_types.capacity AS capacity",
			"users.id AS user_id", "buildings.on_strike AS on_strike", "buildings.workers",
			"building_types.workers AS max_workers").
		Joins("left join buildings on store_goods.building_id = buildings.id").
		Joins("left join building_types on buildings.type_id = building_types.id").
		Joins("left join users on users.id = buildings.user_id").
		Where("on_strike = ?", false).
		Find(&storeGoodsResult)
	if res.Error != nil {
		return nil, res.Error
	}
	return storeGoodsResult, nil
}
