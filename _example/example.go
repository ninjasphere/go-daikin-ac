package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/ninjasphere/go-daikin-ac"
	"github.com/ninjasphere/go-daikin-ac/emulator"
)

func main() {

	emulator.StartWirelessAC(9999)

	ac := daikin.NewWirelessAC("localhost:9999", time.Second*3)

	go func() {
		for state := range ac.OnStateUpdate() {
			log.Printf("Inside temperature: %f. Target temperature %f.", state.InsideTemperature, state.TargetTemperature)
		}
	}()

	ac.ControlState().Power = daikin.PowerOn
	ac.ControlState().Fan = daikin.FanAuto
	ac.ControlState().FanDirection = daikin.FanDirectionHorizontal
	ac.ControlState().Mode = daikin.ModeCool
	ac.ControlState().TargetTemperature = 21

	err := ac.SendState()
	if err != nil {
		panic(err)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)

	// Block until a signal is received.
	s := <-c
	fmt.Println("Got signal:", s)

}
