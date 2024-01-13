package evolution

import (
	"backend/packages/models"
	"fmt"
	"gorm.io/gorm"
	"log"
	"math"
)

// Hiring formula is Worker N = HiringNeeds *(1 - alpha*( (salary N / salary Max) - 1)^2 )
// Count alpha is ( Σ HiringNeeds N - cell.population ) / Σ HiringNeeds N ((salary N / salary Max) - 1)
func Hiring(db *gorm.DB) {
	buildings, err := models.GetBuildingsForHiring(db)
	if err != nil {
		log.Fatalln(err)
	}
	cells := models.GetAllCells(db)
	settings := models.GetSettingsMap(db)
	for y := int(settings["mapMinY"]); y <= int(settings["mapMaxY"]); y++ {
		for x := int(settings["mapMinX"]); x <= int(settings["mapMaxX"]); x++ {
			buildingsInCell := findBuildingsInCell(&buildings, x, y)
			if len(buildingsInCell) == 0 {
				continue
			}

			var zeroWorkerBuildings []models.Building
			cell := getCell(&cells, x, y)
			salaryMax := findMaxSalary(&buildingsInCell)
			fmt.Println(salaryMax)
			again := true
			for again {
				again = false
				numerator := 0.0
				denominator := 0.0
				for _, building := range buildingsInCell {
					numerator += float64(building.HiringNeeds)
					denominator += float64(building.HiringNeeds) * math.Pow((building.Salary/salaryMax)-1, 2)
				}
				fmt.Println(numerator, denominator)
				alpha := (numerator - cell.Population) / denominator
				if alpha < 0 {
					alpha = 0
				}
				for bIndex, building := range buildingsInCell {
					newWorkers := int(float64(building.HiringNeeds) * (1 - alpha*math.Pow((building.Salary/salaryMax)-1, 2)))
					fmt.Println(newWorkers)
					if newWorkers < 0 {
						buildingsInCell = append(buildingsInCell[:bIndex], buildingsInCell[bIndex+1:]...)
						zeroWorkerBuildings = append(zeroWorkerBuildings, building)
						again = true
					}
				}
				if !again {
					for bIndex, building := range buildingsInCell {
						buildingsInCell[bIndex].Workers = int(float64(building.HiringNeeds) * (1 - alpha*math.Pow((building.Salary/salaryMax)-1, 2)))
					}
				}
			}
			db.Save(&buildingsInCell)
			for bIndex := range zeroWorkerBuildings {
				zeroWorkerBuildings[bIndex].Workers = 0
			}
			db.Save(&zeroWorkerBuildings)

			fmt.Println(x, y)
			fmt.Println(cell.Population)

			for _, building := range buildingsInCell {
				fmt.Println(building.ID, building.HiringNeeds, building.Salary)
			}
		}
	}
}

func findBuildingsInCell(buildings *[]models.Building, x int, y int) []models.Building {
	var buildingsInCell []models.Building
	for _, building := range *buildings {
		if building.X == x && building.Y == y {
			buildingsInCell = append(buildingsInCell, building)
		}
	}
	return buildingsInCell
}

func getCell(cells *[]models.CellResult, x int, y int) models.CellResult {
	var cellFound models.CellResult
	for _, cell := range *cells {
		if cell.X == x && cell.Y == y {
			cellFound = cell
			break
		}
	}
	return cellFound
}

func findMaxSalary(buildingsInCell *[]models.Building) float64 {
	salaryMax := 0.0
	for _, building := range *buildingsInCell {
		if building.Salary > salaryMax {
			salaryMax = building.Salary
		}
	}
	return salaryMax
}
