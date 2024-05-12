package models

type EquipmentType struct {
	Id             uint    `json:"id" bson:"id"`
	Name           string  `json:"name" bson:"name"`
	ResourceTypeId uint    `json:"resourceTypeId" bson:"resourceTypeId"`
	Durability     int     `json:"durability" bson:"durability"`
	BlueprintIds   []uint  `json:"blueprintIds" bson:"blueprintIds"`
	EffectId       uint    `json:"effectId" bson:"effectId"`
	Value          float64 `json:"value" bson:"value"`
	Square         float64 `json:"square" bson:"square"`
}
