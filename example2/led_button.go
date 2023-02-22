package example2

import (
	"machine"
	"time"
)

var (
	input = machine.D2
	led   = machine.D13 // onboard LED
)

func Run() {
	led.Configure(machine.PinConfig{Mode: machine.PinOutput})
	input.Configure(machine.PinConfig{Mode: machine.PinInput})

	for {
		led.Set(!input.Get())
		time.Sleep(250 * time.Millisecond)
	}
}
