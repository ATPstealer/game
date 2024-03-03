package gameLive

import (
	"backend/packages/models"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"time"
)

func LogisticsDone(m *mongo.Database) {
	filter := bson.D{{"workEnd", bson.D{{"$lt", time.Now()}}}}
	cursor, err := m.Collection("logistics").Find(context.TODO(), filter)
	if err != nil {
		log.Println("Logistics: " + err.Error())
		return
	}

	var logistics []models.Logistic
	if err = cursor.All(context.TODO(), &logistics); err != nil {
		log.Println("Logistics: " + err.Error())
		return
	}

	for _, logistic := range logistics {
		err := models.AddResource(m, logistic.ResourceTypeID, logistic.UserID, logistic.ToX, logistic.ToY, logistic.Amount)
		if err != nil {
			log.Println("Logistics: " + err.Error())
		}
		_, err = m.Collection("logistics").DeleteOne(context.TODO(), bson.M{"_id": logistic.ID})
		if err != nil {
			log.Println("Logistics: " + err.Error())
		}
	}
}
