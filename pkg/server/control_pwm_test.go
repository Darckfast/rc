package server

import (
	"fmt"
	"os"
	"rc/configs"
	"rc/shared"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func setup(t *testing.T) {
	servoPin := t.TempDir() + "/"
	os.Mkdir(servoPin+"pwm0/", os.ModePerm)
	configs.P.Servo.Pin = servoPin

	escPin := t.TempDir() + "/"
	os.Mkdir(escPin+"pwm0/", os.ModePerm)
	configs.P.Esc.Pin = escPin
}

func TestInitParams_ShouldCreateAndSetTheInitialValues(t *testing.T) {
	setup(t)

	InitPins()

	b, err := os.ReadFile(configs.P.Servo.Pin + "pwm0/period")

	assert.Nil(t, err)
	assert.Equal(t, []byte(fmt.Sprintf("%d", configs.P.Esc.Period)), b)

	b, err = os.ReadFile(configs.P.Servo.Pin + "pwm0/duty_cycle")

	assert.Nil(t, err)
	assert.Equal(t, []byte(fmt.Sprintf("%d", configs.P.Esc.Neutral)), b)

	b, err = os.ReadFile(configs.P.Servo.Pin + "pwm0/polarity")

	assert.Nil(t, err)
	assert.Equal(t, []byte(configs.P.Esc.Polarity), b)

	b, err = os.ReadFile(configs.P.Servo.Pin + "pwm0/enable")

	assert.Nil(t, err)
	assert.Equal(t, []byte("1"), b)

	b, err = os.ReadFile(configs.P.Esc.Pin + "pwm0/period")

	assert.Nil(t, err)
	assert.Equal(t, []byte(fmt.Sprintf("%d", configs.P.Esc.Period)), b)

	b, err = os.ReadFile(configs.P.Esc.Pin + "pwm0/duty_cycle")

	assert.Nil(t, err)
	assert.Equal(t, []byte(fmt.Sprintf("%d", configs.P.Esc.Neutral)), b)

	b, err = os.ReadFile(configs.P.Esc.Pin + "pwm0/polarity")

	assert.Nil(t, err)
	assert.Equal(t, []byte(configs.P.Esc.Polarity), b)

	b, err = os.ReadFile(configs.P.Esc.Pin + "pwm0/enable")

	assert.Nil(t, err)
	assert.Equal(t, []byte("1"), b)
}

func TestApplyControl_ShouldWriteValidValues(t *testing.T) {
	setup(t)

	InitPins()
	gamepad := shared.NormalizedGamepad{
		Tr: 100, // they are float between -1 and 1
		Lx: 100, // they are float between -1 and 1
	}

	ApplyControls(&gamepad)

	b, err := os.ReadFile(configs.P.Servo.Pin + "pwm0/duty_cycle")
	bn, _ := strconv.Atoi(string(b))
	assert.Nil(t, err)
	assert.GreaterOrEqual(t, configs.P.Servo.Max, uint32(bn))
	assert.LessOrEqual(t, configs.P.Servo.Min, uint32(bn))
	assert.NotEqual(t, configs.P.Servo.Neutral, uint32(bn))

	b, err = os.ReadFile(configs.P.Esc.Pin + "pwm0/duty_cycle")
	bn, _ = strconv.Atoi(string(b))
	assert.Nil(t, err)
	assert.GreaterOrEqual(t, configs.P.Esc.Max, uint32(bn))
	assert.LessOrEqual(t, configs.P.Esc.Min, uint32(bn))
	assert.NotEqual(t, configs.P.Esc.Neutral, uint32(bn))
}
