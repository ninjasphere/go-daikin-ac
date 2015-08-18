package daikin

import (
	"fmt"
	"net/url"
	"strconv"
	"time"
)

type baseDaikin struct {
	host            string
	refreshInterval time.Duration

	controlState ControlState
	sensorState  SensorState
	listeners    []chan ACState
}

type DaikinAC interface {
	SendState() error
	RefreshState() (*ControlState, *SensorState, error)
	ControlState() *ControlState
	SensorState() *SensorState
	OnStateUpdate() chan ACState
}

func (d *baseDaikin) ControlState() *ControlState {
	return &d.controlState
}

func (d *baseDaikin) SensorState() *SensorState {
	return &d.sensorState
}

func (d *baseDaikin) OnStateUpdate() chan ACState {
	c := make(chan ACState)

	d.listeners = append(d.listeners, c)

	return c
}

func (d *baseDaikin) emitStateUpdate() {
	s := ACState{d.controlState, d.sensorState}
	for _, c := range d.listeners {
		go func() {
			select {
			case c <- s:
			default:
			}
		}()
	}
}

type ACState struct {
	ControlState
	SensorState
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

func (s *ControlState) GetWirelessValues() url.Values {
	return url.Values{
		"pow":    []string{s.Power.wireless},
		"mode":   []string{s.Mode.wireless},
		"stemp":  []string{fmt.Sprintf("%.2f", s.TargetTemperature)},
		"shum":   []string{fmt.Sprintf("%d", s.TargetHumidity)},
		"f_rate": []string{s.Fan.wireless},
		"f_dir":  []string{s.FanDirection.wireless},
	}
}

func (s *ControlState) ParseWirelessValues(values url.Values) {
	s.Power = parsePower(true, values["pow"][0])
	s.Mode = parseMode(true, values["mode"][0])

	var err error

	s.TargetTemperature, err = strconv.ParseFloat(values["stemp"][0], 64)
	if err != nil {
		fmt.Printf("Warning: Couldn't parse target temperature: %s", err)
	}

	targetHumidity, err := strconv.ParseInt(values["shum"][0], 10, 64)
	if err != nil {
		fmt.Printf("Warning: Couldn't parse target humidity: %s", err)
	}
	s.TargetHumidity = int(targetHumidity)

	s.Fan = parseFan(true, values["f_rate"][0])
	s.FanDirection = parseFanDirection(true, values["f_dir"][0])
}

type SensorState struct {
	// sensor info
	InsideTemperature  float64
	InsideHumidity     float64
	OutsideTemperature float64
}
