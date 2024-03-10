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
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(60*time.Second))
	defer cancel()

	filter := bson.D{{"workEnd", bson.D{{"$lt", time.Now()}}}}
	cursor, err := m.Collection("logistics").Find(ctx, filter)
	if err != nil {
		log.Println("Logistics: " + err.Error())
		return
	}

	var logistics []models.Logistic
	if err = cursor.All(ctx, &logistics); err != nil {
		log.Println("Logistics: " + err.Error())
		return
	}

	for _, logistic := range logistics {
		err := models.AddResource(m, logistic.ResourceTypeID, logistic.UserID, logistic.ToX, logistic.ToY, logistic.Amount)
		if err != nil {
			log.Println("Logistics: " + err.Error())
		}
		_, err = m.Collection("logistics").DeleteOne(ctx, bson.M{"_id": logistic.ID})
		if err != nil {
			log.Println("Logistics: " + err.Error())
		}
	}
}
