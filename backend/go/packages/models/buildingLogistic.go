package models

type Logistics struct {
	CapacityMax float64 `json:"capacityMax" bson:"capacityMax"`
	Capacity    float64 `json:"capacity" bson:"capacity"`
	Speed       float64 `json:"speed" bson:"speed"`
	Price       float64 `json:"price" bson:"price"`
}

//func LogisticsRecount(m *mongo.Database)
