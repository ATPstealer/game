package evolution

import (
	"backend/packages/models"
	"fmt"
	"gorm.io/gorm"
	"log"
)

type AverageSalary struct {
	salary float64
	worker int
	x      int
	y      int
}

func Payroll(db *gorm.DB) {
	buildings, err := models.GetBuildingsForHiring(db)
	if err != nil {
		log.Fatalln(err)
	}
	var averageSalary []AverageSalary

	for bIndex, building := range buildings {
		if models.AddMoney(db, building.UserID, (-1)*float64(building.Workers)*building.Salary) != nil {
			buildings[bIndex].OnStrike = true
			buildings[bIndex].HiringNeeds = 0
			continue
		}
		if err := models.AddCivilSavings(db, building.X, building.Y, float64(building.Workers)*building.Salary); err != nil {
			fmt.Println(err)
		}
		buildings[bIndex].OnStrike = false
		addPayroll(&building, &averageSalary)
	}
	db.Save(&buildings)
	log.Println(averageSalary)
	cellAverageSalary(db, &averageSalary)
}

func addPayroll(building *models.Building, averageSalary *[]AverageSalary) {
	founded := false
	for i, avg := range *averageSalary {
		if avg.x == building.X && avg.y == building.Y {
			(*averageSalary)[i].salary += building.Salary * float64(building.Workers)
			(*averageSalary)[i].worker += building.Workers
			founded = true
			break
		}
	}
	if !founded {
		*averageSalary = append(*averageSalary, AverageSalary{
			salary: building.Salary,
			worker: building.Workers,
			x:      building.X,
			y:      building.Y,
		})
	}
}

func cellAverageSalary(db *gorm.DB, averageSalary *[]AverageSalary) {
	cells := models.GetCells(db)
	for iCell, cell := range cells {
		for _, avg := range *averageSalary {
			if avg.x == cell.X && avg.y == cell.Y {
				cells[iCell].AverageSalary = avg.salary / float64(avg.worker)
			}
		}
	}
	db.Save(&cells)
}
