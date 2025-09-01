package dto

type DeviceListItemDTO struct {
	ID         uint   `json:"id"`
	DeviceName string `json:"device_name"`
	DeviceCode string `json:"device_code"`
	Status     string `json:"status"`
}

type SensorDTO struct {
	ID         uint   `json:"id"`
	SensorType string `json:"sensor_type"`
}

type ActuatorDTO struct {
	ID           uint   `json:"id"`
	ActuatorName string `json:"actuator_name"`
	Type         string `json:"type"`
	PinNumber    string `json:"pin_number"`
	Status       string `json:"status"`
}

type DeviceDetailDTO struct {
	ID         uint          `json:"id"`
	DeviceName string        `json:"device_name"`
	DeviceCode string        `json:"device_code"`
	Status     string        `json:"status"`
	Sensors    []SensorDTO   `json:"sensors"`
	Actuators  []ActuatorDTO `json:"actuators"`
}
