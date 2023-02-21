package example5

import (
	"fmt"
	"machine"
	"time"
)

var adcInputPin = machine.A0

// Wait for user to open serial console
func waitSerial() {
	for !machine.Serial.DTR() {
		time.Sleep(100 * time.Millisecond)
	}
}

func Run() {
	waitSerial()

	machine.InitADC()
	adcInput := machine.ADC{Pin: adcInputPin}
	adcInput.Configure(machine.ADCConfig{})

	for {
		adcReading := adcInput.Get()
		voltage := float64(adcReading) / float64(20_000)Z

		fmt.Printf("Voltage reading: %.2f\r\n", voltage)

		time.Sleep(time.Second)
	}
}
