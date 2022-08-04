package example1

import (
	"machine"
	"time"
)

var (
	led = machine.D4
)

func Run() {
	println("Hello, TinyGo")
	led.Configure(machine.PinConfig{Mode: machine.PinOutput})
	for {
		led.Low()
		time.Sleep(time.Second)

		led.High()
		time.Sleep(time.Second)
	}
}
