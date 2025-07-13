package ws

import "time"

type Devices struct {
	Devices []Device `json:"devices"`
}

type Device struct {
	Name              string     `json:"name"`
	CreateTime        time.Time  `json:"createTime"`
	UpdateTime        time.Time  `json:"updateTime"`
	SerialNumber      string     `json:"serialNumber"`
	RegistrationKey   string     `json:"registrationKey"`
	FirmwareAvailable string     `json:"firmwareAvailable"`
	FirmwareInstalled string     `json:"firmwareInstalled"`
	Type              string     `json:"type"`
	LastHeartbeat     time.Time  `json:"lastHeartbeat"`
	LastConfig        LastConfig `json:"lastConfig"`
	HcMode            string     `json:"hcMode"`
}

type LastConfig struct {
	Name      string    `json:"name"`
	Timestamp time.Time `json:"timestamp"`
	Sentio    Sentio    `json:"sentio"`
}

type Sentio struct {
	Title                     string                     `json:"title"`
	TitlePersonalized         string                     `json:"titlePersonalized"`
	Rooms                     []Room                     `json:"rooms"`
	OutdoorTemperatureSensors []OutdoorTemperatureSensor `json:"outdoorTemperatureSensors"`
	HcMode                    string                     `json:"hcMode"`
	AutomaticHcMode           string                     `json:"automaticHcMode"`
	AvailableHcModes          []string                   `json:"availableHcModes"`
	StandbyMode               string                     `json:"standbyMode"`
	VacationSettings          VacationSettings           `json:"vacationSettings"`
	QuietSettings             QuietSettings              `json:"quietSettings"`
}

type OutdoorTemperatureSensor struct {
	ID                 string  `json:"id"`
	OutdoorTemperature float64 `json:"outdoorTemperature"`
}

type Room struct {
	ID                      string                   `json:"id"`
	Title                   string                   `json:"title"`
	TitlePersonalized       string                   `json:"titlePersonalized"`
	AirTemperature          float64                  `json:"airTemperature"`
	Humidity                float64                  `json:"humidity"`
	SetpointTemperature     float64                  `json:"setpointTemperature"`
	MinSetpointTemperature  float64                  `json:"minSetpointTemperature"`
	MaxSetpointTemperature  float64                  `json:"maxSetpointTemperature"`
	VacationMode            string                   `json:"vacationMode"`
	LockMode                string                   `json:"lockMode"`
	TemperatureState        string                   `json:"temperatureState"`
	TemperaturePresets      []TemperaturePreset      `json:"temperaturePresets"`
	SystemModes             []SystemMode             `json:"systemModes"`
	DehumidificationPresets []DehumidificationPreset `json:"dehumidificationPresets"`
	DehumidifierState       string                   `json:"dehumidifierState"`
	WeeklySchedule          WeeklySchedule           `json:"weeklySchedule"`
}

type TemperaturePreset struct {
	Type                   string  `json:"type"`
	HcMode                 string  `json:"hcMode"`
	SetpointTemperature    float64 `json:"setpointTemperature"`
	MinSetpointTemperature float64 `json:"minSetpointTemperature"`
	MaxSetpointTemperature float64 `json:"maxSetpointTemperature"`
}

type SystemMode struct {
	Type                   string  `json:"type"`
	HcMode                 string  `json:"hcMode"`
	SetpointTemperature    float64 `json:"setpointTemperature"`
	MinSetpointTemperature float64 `json:"minSetpointTemperature"`
	MaxSetpointTemperature float64 `json:"maxSetpointTemperature"`
}

type DehumidificationPreset struct {
	HcMode              string  `json:"hcMode"`
	Setpoint            float64 `json:"setpoint"`
	MinHumiditySetpoint float64 `json:"minHumiditySetpoint"`
	MaxHumiditySetpoint float64 `json:"maxHumiditySetpoint"`
}

type WeeklySchedule struct {
	DefaultPresetType string     `json:"defaultPresetType"`
	ScheduleMode      string     `json:"scheduleMode"`
	Intervals         []Interval `json:"intervals"`
}

type Interval struct {
	Day              string   `json:"day"`
	PresetTimeframes []string `json:"presetTimeframes"`
}

type VacationSettings struct {
	VacationMode      string `json:"vacationMode"`
	VacationModeUntil string `json:"vacationModeUntil"`
}

type QuietSettings struct {
	Mode string `json:"mode"`
}
