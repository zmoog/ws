package ws

type Location struct {
	Ulc          string `json:"ulc"`
	Registration string `json:"registrationKey"`
	SerialNumber int    `json:"serialNumber"`
	Attributes   struct {
		Mode          string  `json:"mode"`
		VacationOn    bool    `json:"vacationOn"`
		VacationUntil *string `json:"vacationUntil"`
		Outdoor       struct {
			Temperature float64 `json:"temperature"`
		} `json:"outdoor"`
		Dst bool `json:"dst"`
	} `json:"attributes"`
}

type Room struct {
	Name             string  `json:"name"`
	Code             string  `json:"code"`
	Season           string  `json:"season"`
	Status           string  `json:"status"`
	Thermo           string  `json:"thermoStatus"`
	Dryer            string  `json:"dryerStatus"`
	TempDesired      float64 `json:"tempDesired"`
	TempAlarmLow     float64 `json:"tempAlarmLow"`
	TempCurrent      float64 `json:"tempCurrent"`
	TempAirCurrent   float64 `json:"tempAirCurrent"`
	TempFloorCurrent float64 `json:"tempFloorCurrent"`
	TempManual       float64 `json:"tempManual"`
	TempEco          float64 `json:"tempEco"`
	TempComfort      float64 `json:"tempComfort"`
	TempExtra        float64 `json:"tempExtra"`
	TempLimit        struct {
		Minimum float64 `json:"minimum"`
		Maximum float64 `json:"maximum"`
	} `json:"tempLimit"`
	TempLimitEco struct {
		Minimum float64 `json:"minimum"`
		Maximum float64 `json:"maximum"`
	} `json:"tempLimitEco"`
	TempLimitComfort struct {
		Minimum float64 `json:"minimum"`
		Maximum float64 `json:"maximum"`
	} `json:"tempLimitComfort"`
	TempLimitExtra struct {
		Minimum float64 `json:"minimum"`
		Maximum float64 `json:"maximum"`
	} `json:"tempLimitExtra"`
	HumidityDesired float64 `json:"humidityDesired"`
	HumidityCurrent float64 `json:"humidityCurrent"`
}
