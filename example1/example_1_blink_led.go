package example1

import (
	"machine"
	"time"
)

var (
	led = machine.D13 // onboard LED
)

func Run() {
	println("Hello, TinyGo")
	led.Configure(machine.PinConfig{Mode: machine.PinOutput})
	for {
		println("Led off")
		led.Low()
		time.Sleep(time.Second)

		println("Led on")
		led.High()
		time.Sleep(time.Second)
	}
}
