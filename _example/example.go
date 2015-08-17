package main

import (
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/ninjasphere/go-daikin-ac"
)

func main() {
	ac := daikin.NewWirelessAC("192.168.0.100", time.Second*10)
	ac.GetControlState().Power = daikin.PowerOn
	ac.GetControlState().Fan = daikin.FanAuto
	ac.GetControlState().FanDirection = daikin.FanDirectionHorizontal
	ac.GetControlState().Mode = daikin.ModeCool
	ac.GetControlState().TargetTemperature = 21

	spew.Dump(ac.UpdateState())
}
