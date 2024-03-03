package evolution

import (
	"backend/packages/models"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
)

type AverageSalary struct {
	salary float64
	worker int
	x      int
	y      int
}

func PayrollMongo(m *mongo.Database) {
	buildings, err := models.GetBuildingsForHiringMongo(m)
	if err != nil {
		log.Fatalln(err)
	}
	var averageSalary []AverageSalary

	for bIndex, building := range buildings {
		if models.AddMoneyMongo(m, building.UserID, (-1)*float64(building.Workers)*building.Salary) != nil {
			buildings[bIndex].OnStrike = true
			buildings[bIndex].HiringNeeds = 0
			continue
		}
		if err := models.AddCivilSavingsMongo(m, building.X, building.Y, float64(building.Workers)*building.Salary); err != nil {
			fmt.Println(err)
		}
		buildings[bIndex].OnStrike = false
		addPayrollMongo(&building, &averageSalary)
	}
	saveBuildings(m, &buildings)
	setCellAverageSalaryMongo(m, &averageSalary)
}

func addPayrollMongo(building *models.BuildingMongo, averageSalary *[]AverageSalary) {
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

func setCellAverageSalaryMongo(m *mongo.Database, averageSalary *[]AverageSalary) {
	for _, avg := range *averageSalary {
		filter := bson.M{"x": avg.x, "y": avg.y}
		update := bson.M{
			"$set": bson.M{
				"averageSalary": avg.salary / float64(avg.worker),
			},
		}

		_, err := m.Collection("cells").UpdateOne(context.TODO(), filter, update)
		if err != nil {
			log.Println(err)
		}
	}
}
