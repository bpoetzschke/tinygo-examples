package example2

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
	input.Configure(machine.PinConfig{Mode: machine.PinInput})

	for {
		led.Set(!input.Get())
		time.Sleep(time.Second)
	}
}
