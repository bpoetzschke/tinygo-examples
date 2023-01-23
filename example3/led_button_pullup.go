package example3

import (
	"machine"
	"time"
)

var (
	input = machine.D2
	led   = machine.D4
)

func Run() {
	led.Configure(machine.PinConfig{Mode: machine.PinOutput})
	input.Configure(machine.PinConfig{Mode: machine.PinInputPullup})

	for {
		led.Set(!input.Get())
		time.Sleep(time.Second)
	}
}