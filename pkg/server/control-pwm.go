package server

import "os"

const SERVO_PWM = "/sys/devices/platform/fe6f0010.pwm/pwm/pwmchip1/" // PIN_18
const ESC_PWM = "/sys/devices/platform/fe6f0000.pwm/pwm/pwmchip0/"   // PIN_16

func InitPins() {
	// check for pin 18 and 16
	_, err := os.Stat(SERVO_PWM)

	if err != nil {
		panic(err)
	}

	_, err = os.Stat(ESC_PWM)

	if err != nil {
		panic(err)
	}

	os.WriteFile(SERVO_PWM+"export", []byte("0"), 0200)
	_, err = os.Stat(SERVO_PWM + "pwm0")

	if err != nil {
		panic(err)
	}
}
