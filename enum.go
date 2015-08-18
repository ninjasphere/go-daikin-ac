package daikin

import (
	"fmt"
	"strings"
)

const unsupported = "unsupported"

type command struct {
	wired    string
	wireless string
}

type Power command

var (
	PowerOn  = Power{"On", "1"}
	PowerOff = Power{"Off", "0"}
	powerAll = []Power{PowerOn, PowerOff}
)

func parsePower(wireless bool, v string) Power {
	for _, c := range powerAll {
		if (wireless && cmp(c.wireless, v)) || (!wireless && cmp(c.wired, v)) {
			return c
		}
	}
	fmt.Printf("Warning: Unknown power value: %s", v)
	return PowerOff
}

type Fan command

var (
	FanAuto = Fan{"FAuto", "A"}
	FanF1   = Fan{"Fun1", "3"}
	FanF2   = Fan{"Fun2", "4"}
	FanF3   = Fan{"Fun3", "5"}
	FanF4   = Fan{"Fun4", "6"}
	FanF5   = Fan{"Fun5", "7"}
	FanNone = Fan{"FunNone", "B"}
	fanAll  = []Fan{FanAuto, FanF1, FanF2, FanF3, FanF4, FanF5, FanNone}
)

func parseFan(wireless bool, v string) Fan {
	for _, c := range fanAll {
		if (wireless && cmp(c.wireless, v)) || (!wireless && cmp(c.wired, v)) {
			return c
		}
	}
	fmt.Printf("Warning: Unknown fan value: %s", v)
	return FanNone
}

type FanDirection command

var (
	FanDirectionOff                   = FanDirection{"Off", ""}
	FanDirectionNone                  = FanDirection{"None", "0"}
	FanDirectionVertical              = FanDirection{"Ud", "1"}
	FanDirectionHorizontal            = FanDirection{unsupported, "2"}
	FanDirectionVerticalAndHorizontal = FanDirection{unsupported, "3"}
	fanDirectionAll                   = []FanDirection{FanDirectionOff, FanDirectionNone, FanDirectionVertical, FanDirectionHorizontal, FanDirectionVerticalAndHorizontal}
)

func parseFanDirection(wireless bool, v string) FanDirection {
	for _, c := range fanDirectionAll {
		if (wireless && cmp(c.wireless, v)) || (!wireless && cmp(c.wired, v)) {
			return c
		}
	}
	fmt.Printf("Warning: Unknown fan direction value: %s", v)
	return FanDirectionNone
}

type Mode command

var (
	// only these options are available for wireless daikins
	ModeAuto = Mode{"Auto", "0"}
	ModeDry  = Mode{"Dry", "2"}
	ModeCool = Mode{"Cool", "3"}
	ModeHeat = Mode{"Heat", "4"}
	ModeFan  = Mode{"Fan", "6"}

	// the non-wireless daikins also support the following:
	ModeOnlyFun = Mode{"OnlyFun", unsupported}
	ModeNight   = Mode{"Night", unsupported}
	ModeNone    = Mode{"None", ""}

	modeAll = []Mode{ModeAuto, ModeDry, ModeCool, ModeHeat, ModeFan, ModeOnlyFun, ModeNight, ModeNone}
)

func parseMode(wireless bool, v string) Mode {
	for _, c := range modeAll {
		if (wireless && cmp(c.wireless, v)) || (!wireless && cmp(c.wired, v)) {
			return c
		}
	}
	fmt.Printf("Warning: Unknown mode value: %s", v)
	return ModeNone
}

func cmp(a, b string) bool {
	return strings.ToUpper(a) == strings.ToUpper(b)
}

type Timer command

var (
	TimerOffOff = Timer{"OFF/OFF", unsupported}
	TimerOnOff  = Timer{"ON/OFF", unsupported}
	TimerOffOn  = Timer{"OFF/ON", unsupported}
	TimerOnOn   = Timer{"ON/ON", unsupported}
	TimerNone   = Timer{"NONE", unsupported}
)
