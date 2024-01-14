package evolution

import (
	"backend/packages/models"
	"fmt"
	"gorm.io/gorm"
	"log"
)

func Payroll(db *gorm.DB) {
	buildings, err := models.GetBuildingsForHiring(db)
	if err != nil {
		log.Fatalln(err)
	}
	for bIndex, building := range buildings {
		if models.AddMoney(db, building.UserID, (-1)*float64(building.Workers)*building.Salary) != nil {
			buildings[bIndex].OnStrike = true
			continue
		}
		if err := models.AddCivilSavings(db, building.X, building.Y, float64(building.Workers)*building.Salary); err != nil {
			fmt.Println(err)
		}
		buildings[bIndex].OnStrike = false
	}
	db.Save(&buildings)
}
