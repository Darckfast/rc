package configs

type Pwm struct {
	Pin      string
	Max      uint32
	Neutral  uint32
	Min      uint32
	Polarity string
	Period   uint32
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
		Max:      1700000,
		Min:      1300000,
		Neutral:  1500000,
		Polarity: "normal",
		Period:   20000000,
	},
	Esc: esc{
		Pwm: Pwm{
			Pin:      "/sys/devices/platform/fe6f0000.pwm/pwm/pwmchip0/", // PIN_16
			Max:      1600000,
			Min:      1300000,
			Neutral:  1500000,
			Polarity: "normal",
			Period:   20000000,
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
