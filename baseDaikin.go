package daikin

import "time"

type baseDaikin struct {
	host            string
	refreshInterval time.Duration

	controlState ControlState
	sensorState  SensorState
}

type DaikinAC interface {
	UpdateState() error
	GetControlState() *ControlState
	GetSensorState() *SensorState
}

func (d *baseDaikin) GetControlState() *ControlState {
	return &d.controlState
}

func (d *baseDaikin) GetSensorState() *SensorState {
	return &d.sensorState
}

func defaultControlState() ControlState {
	return ControlState{
		Power:        PowerOff,
		Mode:         ModeNone,
		Fan:          FanNone,
		FanDirection: FanDirectionNone,
		Timer:        TimerNone,
	}
}

type ControlState struct {
	Power             Power
	TargetTemperature float64
	TargetHumidity    int
	Mode              Mode
	Fan               Fan
	FanDirection      FanDirection
	Timer             Timer
}

type SensorState struct {
	// sensor info
	InsideTemperature  float64
	InsideHumidity     float64
	OutsideTemperature float64
}
