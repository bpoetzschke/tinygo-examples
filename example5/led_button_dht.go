package example5

import (
	"fmt"
	"machine"
	"time"
	"tinygo.org/x/drivers/dht"
)

var (
	input           = machine.D2
	dhtInput        = machine.D3
	led             = machine.D4
	lastButtonState = false
)

func Run() {
	led.Configure(machine.PinConfig{Mode: machine.PinOutput})
	input.Configure(machine.PinConfig{Mode: machine.PinInputPullup})
	dhtSensor := dht.New(dhtInput, dht.DHT22)

	for {
		buttonState := !input.Get()
		if buttonState != lastButtonState {
			lastButtonState = buttonState
			led.Set(buttonState)
			if buttonState {
				temp, hum, err := dhtSensor.Measurements()
				if err == nil {
					println(fmt.Sprintf("Temperature: %02d.%d°C, Humidity: %02d.%d%%", temp/10, temp%10, hum/10, hum%10))
				} else {
					println(fmt.Sprintf("Could not take measurements from the sensor: %s", err.Error()))
				}
			}
		}
		time.Sleep(10 * time.Millisecond)
	}
}
