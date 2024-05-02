package example1

import (
	"machine"
	"time"
)

var (
	//led = machine.D12 // onboard LED
	led = machine.D3
)

func Run() {
	println("Hello, TinyGo")
	led.Configure(machine.PinConfig{Mode: machine.PinOutput})
	for {
		println("Led off")
		led.Low()
		time.Sleep(500 * time.Millisecond)

		println("Led on")
		led.High()
		time.Sleep(1000 * time.Millisecond)
	}
}
