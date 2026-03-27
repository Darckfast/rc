package configs

type Pwm struct {
	Pin      string
	Max      uint32
	Neutral  uint32
	Min      uint32
	Polarity string
	Period   uint32
	Scale    uint32
}

type move struct {
	Init  uint32
	Scale uint32
}

type esc struct {
	Pwm
	Reverse move
	Forward move
}

type conf struct {
	Servo Pwm
	Esc   esc
}

var P = conf{
	Servo: Pwm{
		Pin:      "/sys/devices/platform/fe6f0010.pwm/pwm/pwmchip1/", // PIN_18
		Max:      1_700_000,
		Min:      1_300_000,
		Neutral:  1_500_000,
		Polarity: "normal",
		Period:   20_000_000,
		Scale:    200_000,
	},
	Esc: esc{
		Pwm: Pwm{
			Pin:      "/sys/devices/platform/fe6f0000.pwm/pwm/pwmchip0/", // PIN_16
			Max:      1_600_000,
			Min:      1_300_000,
			Neutral:  1_500_000,
			Polarity: "normal",
			Period:   20_000_000,
		},
		Reverse: move{
			Init:  45_000,
			Scale: 35_000,
		},
		Forward: move{
			Init:  35_000,
			Scale: 25_000,
		},
	},
}
