package server

import (
	"fmt"
	"log"
	"math"
	"os"
	"rc/shared"
)

const servo_pwm_pin_18 = "/sys/devices/platform/fe6f0010.pwm/pwm/pwmchip1/" // PIN_18
const esc_pwm_pin_16 = "/sys/devices/platform/fe6f0000.pwm/pwm/pwmchip0/"   // PIN_16

const neutral_duty_cycle = 1500000
const polarity = "normal"
const period = 20000000

func InitPins() {
	setInitParams(servo_pwm_pin_18, period, neutral_duty_cycle, polarity)
	log.Println("servo pwm enabled")

	setInitParams(esc_pwm_pin_16, period, neutral_duty_cycle, polarity)
	log.Println("esc pwm enabled")
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
	log.Println(gamepad.Lx, gamepad.Ly)
	if gamepad.Lx > 1 {
		gamepad.Lx = 1
	} else if gamepad.Lx < -1 {
		gamepad.Lx = 1
	}

	// limits steering to +/- 20k in duty cycle
	steeringCycle := int(math.Round(gamepad.Lx * 20_000))

	log.Printf("%d\n", neutral_duty_cycle+steeringCycle)

}
