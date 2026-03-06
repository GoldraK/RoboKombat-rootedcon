//go:build windows

package main

import (
	"fmt"
	"syscall"
	"unsafe"
)

const (
	XINPUT_GAMEPAD_LEFT_THUMB_DEADZONE  = 7849
	XINPUT_GAMEPAD_RIGHT_THUMB_DEADZONE = 8689
)

type XInputState struct {
	PacketNumber uint32
	Gamepad      XInputGamepad
}

type XInputGamepad struct {
	Buttons      uint16
	LeftTrigger  uint8
	RightTrigger uint8
	ThumbLX      int16
	ThumbLY      int16
	ThumbRX      int16
	ThumbRY      int16
}

var (
	xinput         = syscall.NewLazyDLL("xinput1_4.dll")
	xinputGetState = xinput.NewProc("XInputGetState")
	xinputGetCaps  = xinput.NewProc("XInputGetCapabilities")
)

func XInputGetState(userIndex uint32) (*XInputState, error) {
	var state XInputState
	ret, _, _ := xinputGetState.Call(
		uintptr(userIndex),
		uintptr(unsafe.Pointer(&state)),
	)
	if ret != 0 {
		return nil, fmt.Errorf("XInputGetState failed with code %d", ret)
	}
	return &state, nil
}

func applyDeadzone(value int16, deadzone int16) int16 {
	if value > deadzone {
		return value
	} else if value < -deadzone {
		return value
	}
	return 0
}
