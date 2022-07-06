package main

import (
	"machine"
	"time"
)

const led = machine.LED

func main() {
	println("Hello, TinyGo")
	led.Configure(machine.PinConfig{Mode: machine.PinOutput})
	for {
		led.Low()
		time.Sleep(time.Second)

		led.High()
		time.Sleep(time.Second)
	}
}
