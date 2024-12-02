package evolution

import (
	"backend/packages/models"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
)

func BankCount(m *mongo.Database) {
	LimitsCount(m)
	CountPayments(m)
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
		bank.LoansAmountNewUsers = building.Bank.LoansAmountNewUsers
	} else {
		bank.LoansAmount = 0
		bank.BorrowedFromState = 0
		bank.LoansAmountNewUsers = 0
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

func CountPayments(m *mongo.Database) {
	Loans, err := models.GetAllLoans(m)
	if err != nil {
		log.Println("bank count: can't get loans: ", err)
		return
	}
	for _, loan := range Loans {
		payment := loan.Amount * loan.Rate / 100
		if models.CheckEnoughMoney(m, loan.BorrowerUserId, payment) {
			if err = models.AddMoney(m, loan.BorrowerUserId, (-1)*payment); err != nil {
				log.Println("bank count: can't add money: ", err)
			}
			if err = models.IncreaseCreditRating(m, loan.BorrowerUserId, payment); err != nil {
				log.Println("bank count: can't increase credit rating: ", err)
			}
			if !loan.StateLoan {
				if err = models.AddMoney(m, loan.LenderUserId, payment); err != nil {
					log.Println("bank count: can't add money: ", err)
				}
			}
			if loan.Status != models.Paying {
				if err = models.UpdateLoanStatus(m, loan.Id, models.Paying); err != nil {
					log.Println("bank count: can't update loan status: ", err)
				}
			}

		} else {
			if err = models.IncreaseCreditRating(m, loan.BorrowerUserId, (-10)*payment); err != nil {
				log.Println("bank count: can't increase credit rating: ", err)
			}
			if loan.Status != models.LoanDefault {
				if err = models.UpdateLoanStatus(m, loan.Id, models.LoanDefault); err != nil {
					log.Println("bank count: can't update loan status: ", err)
				}
			}

		}
	}
}
