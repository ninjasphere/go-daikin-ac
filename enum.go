package daikin

const unsupported = "unsupported"

type command struct {
	wired    string
	wireless string
}

type Power command

var (
	PowerOn  = Power{"On", "1"}
	PowerOff = Power{"Off", "0"}
)

type Fan command

var (
	FanAuto = Fan{"FAuto", "A"}
	FanF1   = Fan{"Fun1", "3"}
	FanF2   = Fan{"Fun2", "4"}
	FanF3   = Fan{"Fun3", "5"}
	FanF4   = Fan{"Fun4", "6"}
	FanF5   = Fan{"Fun5", "7"}
	FanNone = Fan{"FunNone", "B"}
)

type FanDirection command

var (
	FanDirectionOff                   = FanDirection{"Off", ""}
	FanDirectionNone                  = FanDirection{"None", "0"}
	FanDirectionVertical              = FanDirection{"Ud", "1"}
	FanDirectionHorizontal            = FanDirection{unsupported, "2"}
	FanDirectionVerticalAndHorizontal = FanDirection{unsupported, "3"}
)

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
)

type Timer int

const (
	TimerOffOff = iota
	TimerOnOff
	TimerOffOn
	TimerOnOn
	TimerNone
)
