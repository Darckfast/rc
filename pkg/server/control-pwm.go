package server

import (
	"os"
)

const servo_pwm = "/sys/devices/platform/fe6f0010.pwm/pwm/pwmchip1/" // PIN_18
// const esc_pwm = "/sys/devices/platform/fe6f0000.pwm/pwm/pwmchip0/"   // PIN_16

const servo_neutral_duty_cycle = "1500000"
const servo_polarity = "normal"
const servo_period = "20000000"

func InitPins() {
	_, err := os.Stat(servo_pwm)

	if err != nil {
		panic(err)
	}

	_, err = os.Stat(servo_pwm + "pwm0")

	if err != nil && os.IsNotExist(err) {
		err = os.WriteFile(servo_pwm+"export", []byte("0"), 0200)

		if err != nil {
			panic(err)
		}
	} else if err != nil {
		panic(err)
	}

	err = os.WriteFile(servo_pwm+"pwm0/polarity", []byte(servo_polarity), 0200)

	if err != nil {
		panic(err)
	}

	err = os.WriteFile(servo_pwm+"pwm0/period", []byte(servo_period), 0200)

	if err != nil {
		panic(err)
	}

	err = os.WriteFile(servo_pwm+"pwm0/duty_cycle", []byte(servo_neutral_duty_cycle), 0200)

	if err != nil {
		panic(err)
	}

	err = os.WriteFile(servo_pwm+"pwm0/enable", []byte("1"), 0200)

	if err != nil {
		panic(err)
	}
}
