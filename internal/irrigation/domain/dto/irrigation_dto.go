package dto

type PlantDTO struct {
	ID               uint               `json:"id"`
	DeviceID         uint64             `json:"device_id"`
	IrrigationRuleID uint64             `json:"irrigation_rule_id"`
	PlantName        string             `json:"plant_name"`
	ImageURL         string             `json:"image_url"`
	Rule             *IrrigationRuleDTO `json:"rule,omitempty"`
}

type IrrigationRuleDTO struct {
	ID                uint    `json:"id"`
	PlantName         string  `json:"plant_name"`
	MinMoisture       float64 `json:"min_moisture"`
	MaxMoisture       float64 `json:"max_moisture"`
	PreferredTemp     float64 `json:"preferred_temp"`
	PreferredHumidity float64 `json:"preferred_humidity"`
}
