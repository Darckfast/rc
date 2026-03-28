package server

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
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

	ticker := time.NewTicker(50 * time.Millisecond)

	go func() {
		for range ticker.C {
			if time.Since(lastPacket) > 100*time.Millisecond {
				log.Println("no packets received, resetting controls to neutral")

				os.WriteFile(filepath.Join(configs.P.Servo.Pin, "pwm0/duty_cycle"), []byte(fmt.Sprintf("%d", configs.P.Servo.Neutral)), 0644)
				os.WriteFile(filepath.Join(configs.P.Esc.Pin, "pwm0/duty_cycle"), []byte(fmt.Sprintf("%d", configs.P.Esc.Neutral)), 0644)
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
		err = os.WriteFile(filepath.Join(p.Pin, "export"), []byte("0"), 0200)

		if err != nil {
			panic(err)
		}
	} else if err != nil {
		panic(err)
	}

	err = os.WriteFile(filepath.Join(p.Pin, "pwm0/period"), []byte(fmt.Sprintf("%d", p.Period)), 0644)
	if err != nil {
		panic(err)
	}

	err = os.WriteFile(filepath.Join(p.Pin, "pwm0/duty_cycle"), []byte(fmt.Sprintf("%d", p.Neutral)), 0644)
	if err != nil {
		panic(err)
	}

	err = os.WriteFile(filepath.Join(p.Pin, "pwm0/polarity"), []byte(p.Polarity), 0644)
	if err != nil {
		panic(err)
	}

	err = os.WriteFile(filepath.Join(p.Pin, "pwm0/enable"), []byte("1"), 0644)
	if err != nil {
		panic(err)
	}
}

func clampUint32(v, max, min uint32) uint32 {
	if v > max {
		return max
	} else if v < min {
		return min
	}

	return v
}

func clampFloat64(v, max, min float64) float64 {
	if v > max {
		return max
	} else if v < min {
		return min
	}

	return v
}

func Steering(gamepad *shared.NormalizedGamepad) uint32 {
	steeringCycle := uint32(gamepad.Lx*float64(configs.P.Servo.Scale)) + configs.P.Servo.Neutral
	steeringCycle = clampUint32(steeringCycle, configs.P.Servo.Max, configs.P.Servo.Min)

	err := os.WriteFile(filepath.Join(configs.P.Servo.Pin, "pwm0/duty_cycle"), []byte(fmt.Sprintf("%d", steeringCycle)), 0644)

	if err != nil {
		// disable pwm and exit the process
		os.WriteFile(filepath.Join(configs.P.Servo.Pin, "pwm0/enable"), []byte("0"), 0644)
		panic(err)
	}

	return steeringCycle
}

func ForwardOrReverse(gamepad *shared.NormalizedGamepad) uint32 {
	escCycle := uint32(configs.P.Esc.Neutral)
	if gamepad.Tr != 0 { // Forward
		escCycle = uint32(gamepad.Tr*float64(configs.P.Esc.Forward.Scale)) + configs.P.Esc.Neutral

		if escCycle != configs.P.Esc.Neutral {
			escCycle += configs.P.Esc.Forward.Init
		}
	}

	// overwrites forward movement
	if gamepad.Tl != 0 { // Reverse
		escCycle = uint32(gamepad.Tl*-float64(configs.P.Esc.Reverse.Scale)) + configs.P.Esc.Neutral

		if escCycle != configs.P.Esc.Neutral {
			escCycle -= configs.P.Esc.Reverse.Init
		}
	}

	escCycle = clampUint32(escCycle, configs.P.Esc.Max, configs.P.Esc.Min)
	err := os.WriteFile(filepath.Join(configs.P.Esc.Pin, "pwm0/duty_cycle"), []byte(fmt.Sprintf("%d", escCycle)), 0644)

	if err != nil {
		// disable pwm and exit the process
		os.WriteFile(configs.P.Esc.Pin+"pwm0/enable", []byte("0"), 0644)
		panic(err)
	}

	return escCycle
}

func ApplyControls(gamepad *shared.NormalizedGamepad, size int) {
	latency := time.Since(lastPacket)
	lastPacket = time.Now()

	gamepad.Lx = clampFloat64(gamepad.Lx, 1, -1)
	gamepad.Tl = clampFloat64(gamepad.Tl, 1, -1)
	gamepad.Tr = clampFloat64(gamepad.Tr, 1, -1)

	st := Steering(gamepad)
	fr := ForwardOrReverse(gamepad)

	log.Printf(
		"Steering: %-4d | Direction: %-4d | Latency: %-6d ms | Size: %-4d B\n",
		st, fr, latency, size,
	)
}
