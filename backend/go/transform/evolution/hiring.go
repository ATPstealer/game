package evolution

import (
	"backend/packages/models"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"math"
	"time"
)

// Hiring Formula is Worker N = HiringNeeds *(1 - alpha*( (salary N / salary Max) - 1)^2 )
// Count alpha is ( Σ HiringNeeds N - cell.population ) / Σ HiringNeeds N ((salary N / salary Max) - 1)
// Alpha is percentage of hired employees from required
// TODO: strike buildings not in count

func Hiring(m *mongo.Database) {
	buildings, err := models.GetBuildingsForHiring(m)
	if err != nil {
		log.Println("Can't get buildings for hiring: " + err.Error())
		return
	}

	cells, err := models.GetAllCells(m)
	if err != nil {
		log.Println("Can't get cells for hiring: " + err.Error())
		return
	}

	settings, err := models.GetSettings(m)
	if err != nil {
		log.Println("Can't get settings for hiring: " + err.Error())
		return
	}

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
			saveBuildings(m, &buildingsInCell)
			saveBuildings(m, &zeroWorkerBuildings)
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

func getCell(cells *[]models.Cell, x int, y int) models.Cell {
	var cellFound models.Cell
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

func saveBuildings(m *mongo.Database, buildings *[]models.Building) {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(60*time.Second))
	defer cancel()

	for _, building := range *buildings {
		filter := bson.M{"_id": building.ID}
		update := bson.M{
			"$set": bson.M{
				"workers":     building.Workers,
				"hiringNeeds": building.HiringNeeds,
				"onStrike":    building.OnStrike,
			},
		}

		_, err := m.Collection("buildings").UpdateOne(ctx, filter, update)
		if err != nil {
			log.Println(err)
		}
	}
}
