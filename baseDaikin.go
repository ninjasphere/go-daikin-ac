package daikin

import (
	"fmt"
	"log"
	"net/url"
	"strconv"
	"time"
)

type baseDaikin struct {
	info            *BasicInfo
	host            string
	refreshInterval time.Duration

	controlState ControlState
	sensorState  SensorState
	listeners    []chan ACState
}

type AC interface {
	AutoRefresh(interval time.Duration)
	BasicInfo() *BasicInfo
	RefreshBasicInfo() (*BasicInfo, error)
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

func (d *baseDaikin) BasicInfo() *BasicInfo {
	return d.info
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

	if has(values, "stemp") {

		s.TargetTemperature, err = strconv.ParseFloat(values["stemp"][0], 64)
		if err != nil {
			log.Printf("Warning: Couldn't parse target temperature: %s - %s", values["stemp"][0], err)
		}
	}

	if has(values, "shum") {

		targetHumidity, err := strconv.ParseInt(values["shum"][0], 10, 64)
		if err != nil {
			log.Printf("Warning: Couldn't parse target humidity: %s", err)
		}
		s.TargetHumidity = int(targetHumidity)

	}

	if has(values, "f_rate") {
		s.Fan = parseFan(true, values["f_rate"][0])
	}

	if has(values, "f_dir") {
		s.FanDirection = parseFanDirection(true, values["f_dir"][0])
	}
}
func has(values url.Values, name string) bool {
	_, ok := values[name]
	return ok
}
func getVal(values url.Values, name string) (val string) {

	if v, ok := values[name]; ok {
		val = v[0]
	}

	return
}

type SensorState struct {
	// sensor info
	InsideTemperature  float64
	InsideHumidity     int
	OutsideTemperature float64
}

func (s *SensorState) GetWirelessValues() url.Values {
	return url.Values{
		"htemp": []string{fmt.Sprintf("%.2f", s.InsideTemperature)},
		"hhum":  []string{fmt.Sprintf("%d", s.InsideHumidity)},
		"otemp": []string{fmt.Sprintf("%.2f", s.OutsideTemperature)},
	}
}

func (s *SensorState) ParseWirelessValues(values url.Values) {
	var err error

	s.InsideTemperature, err = strconv.ParseFloat(values["htemp"][0], 64)
	if err != nil {
		log.Printf("Warning: Couldn't parse inside temperature: %s - %s", values["htemp"][0], err)
	}

	insideHumidity, err := strconv.ParseInt(values["hhum"][0], 10, 64)
	if err != nil {
		log.Printf("Warning: Couldn't parse inside temperature: %s", err)
	}
	s.InsideHumidity = int(insideHumidity)

	s.OutsideTemperature, err = strconv.ParseFloat(values["otemp"][0], 64)
	if err != nil {
		log.Printf("Warning: Couldn't parse inside temperature: %s", err)
	}

}
