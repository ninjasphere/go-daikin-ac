package daikin

import (
	"fmt"
	"time"
)

const (
	getControlInfo = "/aircon/get_control_info"
	setControlInfo = "/aircon/set_control_info"
	getSensorInfo  = "/aircon/get_sensor_info"
)

func NewWirelessAC(host string, refreshInterval time.Duration) DaikinAC {
	return &wirelessDaikin{baseDaikin{
		host:            host,
		refreshInterval: refreshInterval,
		controlState:    defaultControlState(),
	}}
}

type wirelessDaikin struct {
	baseDaikin
}

func (d *wirelessDaikin) UpdateState() error {
	query := map[string]string{
		"pow":    d.controlState.Power.wireless,
		"mode":   d.controlState.Mode.wireless,
		"stemp":  fmt.Sprintf("%.2f", d.controlState.TargetTemperature),
		"shum":   fmt.Sprintf("%d", d.controlState.TargetHumidity),
		"f_rate": d.controlState.Fan.wireless,
		"f_dir":  d.controlState.FanDirection.wireless,
	}

	return post(d.host, setControlInfo, query)
}
