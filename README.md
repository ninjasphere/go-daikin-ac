# Daikin AC - Golang API

Ported from https://bitbucket.org/JonathanGiles/jdaikin by NinjaBlocks.

WIP. This hasn't actually been tested with any actual hardware.

TODO:
- There are other endpoints and data being ignored.
- Is there any way to do discovery?

### Getting started

```
ac := daikin.NewWirelessAC("192.168.0.100", time.Second*3)
```

If you don't actually have a Daikin air conditioner handy (what, you don't carry one with you?), you can fire up the *very* simple emulator.

```
emulator.StartWirelessAC(9999)
ac := daikin.NewWirelessAC("localhost:9999", time.Second*3)
```

### Listening for state updates

```
go func() {
  for state := range ac.OnStateUpdate() {
    log.Printf("Inside temperature: %f. Target temperature %f.", state.InsideTemperature, state.TargetTemperature)
  }
}()
```

### Controlling the AC

```
ac.ControlState().Power = daikin.PowerOn
ac.ControlState().Fan = daikin.FanAuto
ac.ControlState().FanDirection = daikin.FanDirectionHorizontal
ac.ControlState().Mode = daikin.ModeCool
ac.ControlState().TargetTemperature = 21

err := ac.SendState()
if err != nil {
  panic(err)
}
```
