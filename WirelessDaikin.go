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
	ac := &wirelessAC{baseDaikin: baseDaikin{
		host:            host,
		refreshInterval: refreshInterval,
		controlState:    defaultControlState(),
	}}

	ac.timer = time.AfterFunc(refreshInterval, func() {
		_, _, err := ac.RefreshState()
		if err != nil {
			fmt.Sprintf("Failed to refresh AC state: %s", err)
		}
		ac.timer.Reset(refreshInterval)
	})

	return ac
}

type wirelessAC struct {
	baseDaikin
	timer *time.Timer
}

func (d *wirelessAC) SendState() error {
	return post(d.host, setControlInfo, d.ControlState().GetWirelessValues())
}

func (d *wirelessAC) RefreshState() (*ControlState, *SensorState, error) {
	controlVals, err := get(d.host, getControlInfo)

	if err != nil {
		return nil, nil, err
	}

	d.ControlState().ParseWirelessValues(controlVals)

	d.emitStateUpdate()

	return d.ControlState(), d.SensorState(), nil
}
