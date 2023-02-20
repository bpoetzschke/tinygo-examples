package example4

import (
	"machine"
	"time"
)

var (
	led1 = machine.D13 // onboard LED
	led2 = machine.D4
)

func Run() {
	led1.Configure(machine.PinConfig{Mode: machine.PinOutput})
	led2.Configure(machine.PinConfig{Mode: machine.PinOutput})

	go func() {
		for {
			led1.Low()
			time.Sleep(time.Second)

			led1.High()
			time.Sleep(time.Second)
		}
	}()

	go func() {
		for {
			led2.Low()
			time.Sleep(500 * time.Millisecond)

			led2.High()
			time.Sleep(500 * time.Millisecond)
		}
	}()

	for {
		println(".")
		time.Sleep(10 * time.Second)
	}
}
