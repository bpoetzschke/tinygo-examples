package main

import (
	"fmt"
	"machine"
	"time"
	"tinygo.org/x/drivers/wifinina"
)

var (
	counter = 0
	buffer  [1024]byte

	// these are the default pins for the Arduino Nano33 IoT.
	spi = machine.NINA_SPI

	// this is the ESP chip that has the WIFININA firmware flashed on it
	adaptor *wifinina.Device

	led   = machine.LED
	input = machine.D2
)

func setup() {
	// Configure SPI for 8Mhz, Mode 0, MSB First
	spi.Configure(machine.SPIConfig{
		Frequency: 8 * 1e6,
		SDO:       machine.NINA_SDO,
		SDI:       machine.NINA_SDI,
		SCK:       machine.NINA_SCK,
	})

	adaptor = wifinina.New(spi,
		machine.NINA_CS,
		machine.NINA_ACK,
		machine.NINA_GPIO0,
		machine.NINA_RESETN)
	adaptor.Configure()
}

// Wait for user to open serial console
func waitSerial() {
	for !machine.Serial.DTR() {
		time.Sleep(100 * time.Millisecond)
	}
}

func PrintErrIfNotNil(err error) {
	if err == nil {
		return
	}

	for {
		println(err)
		time.Sleep(time.Second)
	}
}

// connect to access point
func connectToAP() {
	time.Sleep(2 * time.Second)
	println("Connecting to " + ssid)
	err := adaptor.ConnectToAccessPoint(ssid, pass, 10*time.Second)
	if err != nil { // error connecting to AP
		for {
			println(err)
			time.Sleep(1 * time.Second)
		}
	}

	println("Connected.")

	ip, _, _, err := adaptor.GetIP()
	for ; err != nil; ip, _, _, err = adaptor.GetIP() {
		PrintErrIfNotNil(err)
		time.Sleep(1 * time.Second)
	}
	println(ip.String())
}

func main() {
	waitSerial()

	led.Configure(machine.PinConfig{Mode: machine.PinOutput})
	input.Configure(machine.PinConfig{Mode: machine.PinInputPullup})

	//err := input.SetInterrupt(machine.PinToggle, func(pin machine.Pin) {
	//	println(fmt.Sprintf("Pin toggled to %t", pin.Get()))
	//	led.Set(pin.Get())
	//})
	//PrintErrIfNotNil(err)

	for {
		println(fmt.Sprintf("Pin toggled to %t", input.Get()))
		led.Set(input.Get())
		time.Sleep(time.Second)
	}
}
