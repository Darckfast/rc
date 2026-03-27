package server

import (
	"fmt"
	"log"
	"math"
	"os"
	"rc/configs"
	"rc/shared"
	"time"
)

var lastPacket time.Time = time.Now()

func InitPins() {
	setInitParams(&configs.P.Servo)
	log.Println("servo pwm enabled")

	setInitParams(&configs.P.Esc.Pwm)
	log.Println("esc pwm enabled")

	ticker := time.NewTicker(100 * time.Millisecond)

	go func() {
		for range ticker.C {
			log.Println("latency:", time.Since(lastPacket))
			if time.Since(lastPacket) > 100*time.Millisecond {
				log.Println("no packets received, resetting controls to neutral")
				os.WriteFile(configs.P.Servo.Pin+"pwm0/duty_cycle", []byte(fmt.Sprintf("%d", configs.P.Servo.Neutral)), 0644)
				os.WriteFile(configs.P.Esc.Pin+"pwm0/duty_cycle", []byte(fmt.Sprintf("%d", configs.P.Esc.Neutral)), 0644)
			}
		}
	}()

	log.Println("fail-safe in place, checking every 100ms")
}

func setInitParams(p *configs.Pwm) {
	_, err := os.Stat(p.Pin)

	if err != nil {
		panic(err)
	}

	_, err = os.Stat(p.Pin + "pwm0")

	if err != nil && os.IsNotExist(err) {
		err = os.WriteFile(p.Pin+"export", []byte("0"), 0200)

		if err != nil {
			panic(err)
		}
	} else if err != nil {
		panic(err)
	}

	err = os.WriteFile(p.Pin+"pwm0/period", []byte(fmt.Sprintf("%d", p.Period)), 0644)

	if err != nil {
		panic(err)
	}

	err = os.WriteFile(p.Pin+"pwm0/duty_cycle", []byte(fmt.Sprintf("%d", p.Neutral)), 0644)

	if err != nil {
		panic(err)
	}

	err = os.WriteFile(p.Pin+"pwm0/polarity", []byte(p.Polarity), 0644)

	if err != nil {
		panic(err)
	}

	err = os.WriteFile(p.Pin+"pwm0/enable", []byte("1"), 0644)

	if err != nil {
		panic(err)
	}
}

func clamp(v, max, min uint32) uint32 {
	if v > max {
		return max
	} else if v < min {
		return min
	}

	return v
}

func ApplyControls(gamepad *shared.NormalizedGamepad) {
	lastPacket = time.Now()
	if gamepad.Lx > 1 {
		gamepad.Lx = 1
	} else if gamepad.Lx < -1 {
		gamepad.Lx = -1
	}

	// limits steering to +/- 20k in duty cycle
	steeringCycle := uint32(math.Round(gamepad.Lx*200_000)) + configs.P.Servo.Neutral
	steeringCycle = clamp(steeringCycle, configs.P.Servo.Max, configs.P.Servo.Min)

	err := os.WriteFile(configs.P.Servo.Pin+"pwm0/duty_cycle", []byte(fmt.Sprintf("%d", steeringCycle)), 0644)

	if err != nil {
		// disable pwm and exit the process
		os.WriteFile(configs.P.Servo.Pin+"pwm0/enable", []byte("0"), 0644)
		panic(err)
	}

	if gamepad.Tl > 1 {
		gamepad.Tl = 1
	} else if gamepad.Tl < -1 {
		gamepad.Tl = -1
	}

	if gamepad.Tr > 1 {
		gamepad.Tr = 1
	} else if gamepad.Tr < -1 {
		gamepad.Tr = -1
	}

	if gamepad.Tl != 0 { // Reverse
		escCycle := uint32(math.Round(gamepad.Tl*-float64(configs.P.Esc.Reverse.Scale))) + configs.P.Esc.Neutral

		if escCycle != configs.P.Esc.Neutral {
			escCycle -= configs.P.Esc.Reverse.Init
		}

		escCycle = clamp(escCycle, configs.P.Esc.Max, configs.P.Servo.Min)
		err = os.WriteFile(configs.P.Esc.Pin+"pwm0/duty_cycle", []byte(fmt.Sprintf("%d", escCycle)), 0644)

		if err != nil {
			// disable pwm and exit the process
			os.WriteFile(configs.P.Esc.Pin+"pwm0/enable", []byte("0"), 0644)
			panic(err)
		}
	} else { // Forward
		escCycle := uint32(math.Round(gamepad.Tr*float64(configs.P.Esc.Forward.Scale))) + configs.P.Servo.Neutral

		if escCycle != configs.P.Esc.Neutral {
			escCycle += configs.P.Esc.Forward.Init
		}

		escCycle = clamp(escCycle, configs.P.Esc.Max, configs.P.Servo.Min)
		err = os.WriteFile(configs.P.Esc.Pin+"pwm0/duty_cycle", []byte(fmt.Sprintf("%d", escCycle)), 0644)

		if err != nil {
			// disable pwm and exit the process
			os.WriteFile(configs.P.Esc.Pin+"pwm0/enable", []byte("0"), 0644)
			panic(err)
		}
	}
}
