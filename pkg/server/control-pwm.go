package server

import (
	"fmt"
	"log"
	"math"
	"os"
	"rc/shared"
	"time"
)

const servo_pwm_pin_18 = "/sys/devices/platform/fe6f0010.pwm/pwm/pwmchip1/" // PIN_18
const esc_pwm_pin_16 = "/sys/devices/platform/fe6f0000.pwm/pwm/pwmchip0/"   // PIN_16

const neutral_duty_cycle = 1500000
const polarity = "normal"
const period = 20000000
const init_forward_cycle = 35_000
const init_reverse_cycle = 45_000
const min_esc_cycle = 1300000
const max_esc_cycle = 1550000
const reverse_trigger_cycle = 35_000
const forward_trigger_cycle = 25_000

var lastPacket time.Time = time.Now()

func InitPins() {
	setInitParams(servo_pwm_pin_18, period, neutral_duty_cycle, polarity)
	log.Println("servo pwm enabled")

	setInitParams(esc_pwm_pin_16, period, neutral_duty_cycle, polarity)
	log.Println("esc pwm enabled")

	ticker := time.NewTicker(100 * time.Millisecond)

	go func() {
		for range ticker.C {
			log.Println("latency:", time.Since(lastPacket))
			if time.Since(lastPacket) > 100*time.Millisecond {
				log.Println("no packets received, resetting controls to neutral")
				os.WriteFile(servo_pwm_pin_18+"pwm0/duty_cycle", []byte(fmt.Sprintf("%d", neutral_duty_cycle)), 0644)
				os.WriteFile(esc_pwm_pin_16+"pwm0/duty_cycle", []byte(fmt.Sprintf("%d", neutral_duty_cycle)), 0644)
			}
		}
	}()

	log.Println("fail-safe in place, checking every 100ms")
}

func setInitParams(path string, period uint32, neutral_duty_cycle uint32, polarity string) {
	_, err := os.Stat(path)

	if err != nil {
		panic(err)
	}

	_, err = os.Stat(path + "pwm0")

	if err != nil && os.IsNotExist(err) {
		err = os.WriteFile(path+"export", []byte("0"), 0200)

		if err != nil {
			panic(err)
		}
	} else if err != nil {
		panic(err)
	}

	err = os.WriteFile(path+"pwm0/period", []byte(fmt.Sprintf("%d", period)), 0644)

	if err != nil {
		panic(err)
	}

	err = os.WriteFile(path+"pwm0/duty_cycle", []byte(fmt.Sprintf("%d", neutral_duty_cycle)), 0644)

	if err != nil {
		panic(err)
	}

	err = os.WriteFile(path+"pwm0/polarity", []byte(polarity), 0644)

	if err != nil {
		panic(err)
	}

	err = os.WriteFile(path+"pwm0/enable", []byte("1"), 0644)

	if err != nil {
		panic(err)
	}
}

func Move(gamepad *shared.NormalizedGamepad) {
	lastPacket = time.Now()
	if gamepad.Lx > 1 {
		gamepad.Lx = 1
	} else if gamepad.Lx < -1 {
		gamepad.Lx = -1
	}

	// limits steering to +/- 20k in duty cycle
	steeringCycle := int(math.Round(gamepad.Lx*200_000)) + neutral_duty_cycle

	if steeringCycle < 1300000 {
		steeringCycle = 1300000
	} else if steeringCycle > 1700000 {
		steeringCycle = 1700000
	}

	// log.Printf("%d\n", steeringCycle)
	err := os.WriteFile(servo_pwm_pin_18+"pwm0/duty_cycle", []byte(fmt.Sprintf("%d", steeringCycle)), 0644)

	if err != nil {
		// disable pwm and exit the process
		os.WriteFile(servo_pwm_pin_18+"pwm0/enable", []byte("0"), 0644)
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
		// limits steering to +/- 2k in duty cycle
		escCycle := int(math.Round(gamepad.Tl*-reverse_trigger_cycle)) + neutral_duty_cycle

		if escCycle != neutral_duty_cycle {
			escCycle -= init_reverse_cycle
		}

		if escCycle > max_esc_cycle {
			escCycle = max_esc_cycle
		} else if escCycle < min_esc_cycle {
			escCycle = min_esc_cycle
		}

		err = os.WriteFile(esc_pwm_pin_16+"pwm0/duty_cycle", []byte(fmt.Sprintf("%d", escCycle)), 0644)

		if err != nil {
			// disable pwm and exit the process
			os.WriteFile(esc_pwm_pin_16+"pwm0/enable", []byte("0"), 0644)
			panic(err)
		}
	} else { // Forward
		// limits steering to +/- 2k in duty cycle
		escCycle := int(math.Round(gamepad.Tr*forward_trigger_cycle)) + neutral_duty_cycle

		if escCycle != neutral_duty_cycle {
			escCycle += init_forward_cycle
		}

		if escCycle > max_esc_cycle {
			escCycle = max_esc_cycle
		} else if escCycle < min_esc_cycle {
			escCycle = min_esc_cycle
		}

		err = os.WriteFile(esc_pwm_pin_16+"pwm0/duty_cycle", []byte(fmt.Sprintf("%d", escCycle)), 0644)

		if err != nil {
			// disable pwm and exit the process
			os.WriteFile(esc_pwm_pin_16+"pwm0/enable", []byte("0"), 0644)
			panic(err)
		}
	}
}
