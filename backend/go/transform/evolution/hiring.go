package evolution

import (
	"backend/packages/models"
	"gorm.io/gorm"
	"log"
	"math"
)

// Hiring Formula is Worker N = HiringNeeds *(1 - alpha*( (salary N / salary Max) - 1)^2 )
// Count alpha is ( Σ HiringNeeds N - cell.population ) / Σ HiringNeeds N ((salary N / salary Max) - 1)
// Alpha is percentage of hired employees from required
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

			again := true
			for again {
				again = false

				alpha := getAlpha(&buildingsInCell, cell.Population, salaryMax)
				for bIndex, building := range buildingsInCell {
					newWorkers := int(float64(building.HiringNeeds) * (1 - alpha*math.Pow((building.Salary/salaryMax)-1, 2)))
					if newWorkers < 0 {
						buildingsInCell = append(buildingsInCell[:bIndex], buildingsInCell[bIndex+1:]...)
						building.Workers = 0
						zeroWorkerBuildings = append(zeroWorkerBuildings, building)
						again = true
					}
				}
				if !again {
					for bIndex, building := range buildingsInCell {
						buildingsInCell[bIndex].Workers = int(math.Round(float64(building.HiringNeeds) * (1 - alpha*math.Pow((building.Salary/salaryMax)-1, 2))))
					}
				}
			}
			db.Save(&buildingsInCell)
			db.Save(&zeroWorkerBuildings)
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

func getAlpha(buildingsInCell *[]models.Building, population float64, salaryMax float64) float64 {
	numerator := 0.0
	denominator := 0.0
	for _, building := range *buildingsInCell {
		numerator += float64(building.HiringNeeds)
		denominator += float64(building.HiringNeeds) * math.Pow((building.Salary/salaryMax)-1, 2)
	}
	alpha := (numerator - population) / denominator
	if alpha < 0 {
		return 0
	}
	return alpha
}
