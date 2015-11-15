package daikin

import (
	"fmt"
	"time"
)

const (
	basicInfo      = "/common/basic_info"
	getControlInfo = "/aircon/get_control_info"
	setControlInfo = "/aircon/set_control_info"
	getSensorInfo  = "/aircon/get_sensor_info"
)

func NewWirelessAC(host string) *wirelessAC {
	ac := &wirelessAC{baseDaikin: baseDaikin{
		host:         host,
		controlState: defaultControlState(),
	}}

	return ac
}

type wirelessAC struct {
	baseDaikin
	timer *time.Timer
}

func (d *wirelessAC) AutoRefresh(refreshInterval time.Duration) {
	d.timer = time.AfterFunc(refreshInterval, func() {
		_, _, err := d.RefreshState()
		if err != nil {
			fmt.Printf("Failed to refresh AC state: %s\n", err)
		}
		d.timer.Reset(refreshInterval)
	})
}

func (d *wirelessAC) SendState() error {
	return post(d.host, setControlInfo, d.ControlState().GetWirelessValues())
}

func (d *wirelessAC) RefreshBasicInfo() (*BasicInfo, error) {
	vals, err := get(d.host, basicInfo)

	if err != nil {
		return nil, err
	}

	info := &BasicInfo{}

	err = mapValues(info, vals)
	if err != nil {
		return nil, err
	}

	d.info = info

	return info, nil
}

func (d *wirelessAC) RefreshState() (*ControlState, *SensorState, error) {
	controlVals, err := get(d.host, getControlInfo)

	if err != nil {
		return nil, nil, err
	}

	if err := d.ControlState().ParseWirelessValues(controlVals); err != nil {
		return nil, nil, fmt.Errorf("failed: host=%s: %s", d.host, err)
	}

	sensorVals, err := get(d.host, getSensorInfo)

	if err != nil {
		return nil, nil, err
	}

	d.SensorState().ParseWirelessValues(sensorVals)

	d.emitStateUpdate()

	return d.ControlState(), d.SensorState(), nil
}
