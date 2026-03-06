//go:build !windows

package main

import (
	"fmt"
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

func XInputGetState(userIndex uint32) (*XInputState, error) {
	return nil, fmt.Errorf("XInput is only supported on Windows")
}

func applyDeadzone(value int16, deadzone int16) int16 {
	if value > deadzone {
		return value
	} else if value < -deadzone {
		return value
	}
	return 0
}
