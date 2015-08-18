package emulator

import (
	"fmt"
	"net/http"
	"time"

	"github.com/fatih/structs"
	"github.com/ninjasphere/go-daikin-ac"
)

type emulatedWirelessAC struct {
	controlState daikin.ControlState
	sensorState  daikin.SensorState
}

func StartWirelessAC(port int) {
	ac := &emulatedWirelessAC{}
	ac.start(port)
}

func (d *emulatedWirelessAC) start(port int) {

	http.HandleFunc("/aircon/get_control_info", d.getControlInfo)
	http.HandleFunc("/aircon/set_control_info", d.setControlInfo)
	http.HandleFunc("/aircon/get_sensor_info", d.getSensorInfo)

	go func() {
		for {
			time.Sleep(time.Second * 2)
			d.controlState.TargetTemperature += 0.5
			d.controlState.TargetHumidity++
		}
	}()

	fmt.Printf("Starting emulated Daikin Wireless AC on port %d", port)
	go http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}

func (d *emulatedWirelessAC) getControlInfo(w http.ResponseWriter, r *http.Request) {
	var out string

	for k, v := range d.controlState.GetWirelessValues() {
		out += fmt.Sprintf(",%s=%s", k, v[0])
	}

	//spew.Dump(out)

	fmt.Fprint(w, out[1:])
}

func (d *emulatedWirelessAC) getSensorInfo(w http.ResponseWriter, r *http.Request) {
	var out string

	for k, v := range structs.Map(d.controlState) {
		out += fmt.Sprintf(",%s=%s", k, v)
	}

	fmt.Fprint(w, out[1:])
}

func (d *emulatedWirelessAC) setControlInfo(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	//spew.Dump(r.Form)

	d.controlState.ParseWirelessValues(r.Form)

	fmt.Fprint(w, "Success? Who knows.")
}
