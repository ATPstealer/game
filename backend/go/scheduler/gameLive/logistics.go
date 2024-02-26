package gameLive

import (
	"backend/packages/models"
	"gorm.io/gorm"
	"log"
)

func LogisticsDone(db *gorm.DB) {
	var logistics []models.Logistic
	res := db.Model(&models.Logistic{}).Where("work_end < NOW()").Scan(&logistics)
	if res.Error != nil {
		log.Println(res.Error)
	}
	log.Println(logistics)
	for _, logistic := range logistics {
		err := models.AddResource(db, logistic.ResourceTypeID, logistic.UserID, logistic.ToX, logistic.ToY, logistic.Amount)
		if err != nil {
			log.Println(err)
		}
		db.Delete(&logistic)
	}
}
