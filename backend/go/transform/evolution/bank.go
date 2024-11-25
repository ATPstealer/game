package evolution

import (
	"backend/packages/models"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
)

func BankCount(m *mongo.Database) {
	LimitsCount(m)
}

func LimitsCount(m *mongo.Database) {
	buildings, err := models.GetAllReadyBuildingByGroup(m, "Bank")
	if err != nil {
		log.Println("bank count: can't get banks: ", err)
		return
	}
	buildingTypes, err := models.GetAllBuildingTypes(m)
	if err != nil {
		log.Println("bank count: can't get building types: ", err)
		return
	}
	for _, building := range buildings {
		buildingLimitCont(m, building, &buildingTypes)
	}
}

func buildingLimitCont(m *mongo.Database, building models.Building, buildingTypes *[]models.BuildingType) {
	buildingType := getType(building.TypeId, buildingTypes)
	log.Println(building.Workers, buildingType.Workers)
	efficiency := float64(building.Workers) / float64(buildingType.Workers)
	additionalLoanLimit := findEffect(building.EquipmentEffect, 4) // 4 - Increase loans limit
	loanLimit := 1000000 * efficiency * (1 + additionalLoanLimit/10)
	additionalBorrowingsLimit := findEffect(building.EquipmentEffect, 5) // 5 - Increases limit of borrowings
	borrowedLimit := 500000 * efficiency * (1 + additionalBorrowingsLimit/10)

	bank := models.Bank{
		LoansLimit:    loanLimit,
		BorrowedLimit: borrowedLimit,
	}
	if building.Bank != nil {
		bank.LoansAmount = building.Bank.LoansAmount
		bank.BorrowedFromState = building.Bank.BorrowedFromState
	} else {
		bank.LoansAmount = 0
		bank.BorrowedFromState = 0
	}
	models.UpdateBankLimits(m, building.Id, bank)
}

func getType(id uint, buildingTypes *[]models.BuildingType) models.BuildingType {
	for _, buildingType := range *buildingTypes {
		if buildingType.Id == id {
			return buildingType
		}
	}
	return models.BuildingType{}
}

func findEffect(equipmentEffects *[]models.EquipmentEffect, effectId uint) float64 {
	if equipmentEffects != nil {
		for _, equipmentEffect := range *equipmentEffects {
			if equipmentEffect.EffectId == effectId {
				return equipmentEffect.Value
			}
		}
	}
	return 0
}
