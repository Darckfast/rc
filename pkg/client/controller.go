//go:build windows

package client

import (
	"log"
	"math"
	"syscall"
	"unsafe"
)

// source: https://learn.microsoft.com/en-us/windows/win32/xinput/getting-started-with-xinput
const xinput_dll = "XINPUT1_4.DLL"
const xinput_get_state = "XInputGetState"

var getState *syscall.Proc

// Device Button                 Bitmask
// ----------------------------  --------
// XINPUT_GAMEPAD_DPAD_UP        0x0001
// XINPUT_GAMEPAD_DPAD_DOWN      0x0002
// XINPUT_GAMEPAD_DPAD_LEFT      0x0004
// XINPUT_GAMEPAD_DPAD_RIGHT     0x0008
// XINPUT_GAMEPAD_START          0x0010
// XINPUT_GAMEPAD_BACK           0x0020
// XINPUT_GAMEPAD_LEFT_THUMB     0x0040
// XINPUT_GAMEPAD_RIGHT_THUMB    0x0080
// XINPUT_GAMEPAD_LEFT_SHOULDER  0x0100
// XINPUT_GAMEPAD_RIGHT_SHOULDER 0x0200
// XINPUT_GAMEPAD_A              0x1000
// XINPUT_GAMEPAD_B              0x2000
// XINPUT_GAMEPAD_X              0x4000
// XINPUT_GAMEPAD_Y              0x8000
type Gamepad struct {
	_             uint32
	button_mask   uint16
	left_trigger  uint8
	right_trigger uint8
	left_x        int16
	left_y        int16
	right_x       int16
	right_y       int16
}

func init() {
	dll, err := syscall.LoadDLL(xinput_dll)

	if err != nil {
		log.Println("error loading DLL"+xinput_dll, err)
		panic(err)
	}

	getState, err = dll.FindProc(xinput_get_state)
	if err != nil {
		log.Println("error finding proc get state", err)
		panic(err)
	}
}

const xinput_gamepad_left_thumb_deadzone = 7849
const xinput_gamepad_right_thumb_deadzone = 8689
const xinput_gamepad_trigger_threshold = 30

func calculateDeadZone(gamepad *Gamepad) (float64, float64, float64, float64, float64, float64) {
	lx, ly, _ := normalizeThumb(gamepad.left_x, gamepad.left_y, xinput_gamepad_left_thumb_deadzone)
	rx, ry, _ := normalizeThumb(gamepad.right_x, gamepad.right_y, xinput_gamepad_right_thumb_deadzone)
	_, tl := normalizeTrigger(gamepad.left_trigger, xinput_gamepad_trigger_threshold)
	_, tr := normalizeTrigger(gamepad.right_trigger, xinput_gamepad_trigger_threshold)

	return lx, ly, rx, ry, tl, tr
}

func normalizeTrigger(trigger uint8, deadzone float64) (float64, float64) {
	t64 := float64(trigger)

	if t64 <= deadzone {
		return 0, 0
	}

	t64 -= deadzone
	normalizedTrigger := t64 / (255 - deadzone)
	normalizedTriggerCubic := math.Pow(t64/(255-deadzone), 3)

	return normalizedTrigger, normalizedTriggerCubic
}

func normalizeThumb(x, y int16, deadzone float64) (float64, float64, float64) {
	x64 := float64(x)
	y64 := float64(y)
	magnitude := math.Sqrt((x64 * x64) + (y64 * y64))

	normalizedX := float64(0)
	normalizedY := float64(0)

	if magnitude > 0 {
		normalizedX = x64 / magnitude
		normalizedY = y64 / magnitude
	}

	normalizedMagnitude := float64(0)

	if magnitude <= deadzone {
		return 0, 0, 0
	}

	if magnitude > 32767 {
		magnitude = 32767
	}

	magnitude -= deadzone
	normalizedMagnitude = math.Pow(magnitude/(32767-deadzone), 3)

	return normalizedX, normalizedY, normalizedMagnitude
}

func GetControllerState() {
	var gamepad Gamepad
	for {
		result, _, _ := getState.Call(uintptr(0), uintptr(unsafe.Pointer(&gamepad)))

		if result != 0 {
			log.Println("controller is not connected")
			break
		}

		calculateDeadZone(&gamepad)
	}
}
