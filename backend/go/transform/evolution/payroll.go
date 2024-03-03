package evolution

import (
	"backend/packages/models"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"time"
)

type AverageSalary struct {
	salary float64
	worker int
	x      int
	y      int
}

func Payroll(m *mongo.Database) {
	buildings, err := models.GetBuildingsForHiring(m)
	if err != nil {
		log.Fatalln(err)
	}
	var averageSalary []AverageSalary

	for bIndex, building := range buildings {
		if models.AddMoney(m, building.UserID, (-1)*float64(building.Workers)*building.Salary) != nil {
			buildings[bIndex].OnStrike = true
			buildings[bIndex].HiringNeeds = 0
			continue
		}
		if err := models.AddCivilSavings(m, building.X, building.Y, float64(building.Workers)*building.Salary); err != nil {
			fmt.Println(err)
		}
		buildings[bIndex].OnStrike = false
		addPayroll(&building, &averageSalary)
	}
	saveBuildings(m, &buildings)
	setCellAverageSalary(m, &averageSalary)
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
			salary: building.Salary * float64(building.Workers),
			worker: building.Workers,
			x:      building.X,
			y:      building.Y,
		})
	}
}

func setCellAverageSalary(m *mongo.Database, averageSalary *[]AverageSalary) {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(60*time.Second))
	defer cancel()

	for _, avg := range *averageSalary {
		filter := bson.M{"x": avg.x, "y": avg.y}
		update := bson.M{
			"$set": bson.M{
				"averageSalary": avg.salary / float64(avg.worker),
			},
		}

		_, err := m.Collection("cells").UpdateOne(ctx, filter, update)
		if err != nil {
			log.Println(err)
		}
	}
}
