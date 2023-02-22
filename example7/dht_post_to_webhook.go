package example7

import (
	"fmt"
	"machine"
	"strings"
	"time"

	"tinygo.org/x/drivers/dht"
	"tinygo.org/x/drivers/net"
	"tinygo.org/x/drivers/net/http"
	"tinygo.org/x/drivers/wifinina"
)

var (
	buffer [1120]byte

	// these are the default pins for the Arduino Nano33 IoT.
	spi = machine.NINA_SPI

	// this is the ESP chip that has the WIFININA firmware flashed on it
	adaptor *wifinina.Device

	input           = machine.D2
	dhtInput        = machine.D3
	led             = machine.D4
	lastButtonState = false
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

	time.Sleep(2 * time.Second)
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
	connected := false
	for i := 0; i < 3; i++ {
		println("Connecting to " + ssid)
		err := adaptor.ConnectToAccessPoint(ssid, pass, 10*time.Second)
		if err != nil {
			if err.Error() == net.ErrWiFiConnectTimeout.Error() {
				println("Connection error while connecting to wifi. Will retry.")
				time.Sleep(time.Second)
				continue
			}
			PrintErrIfNotNil(err)
		}
		connected = true
	}

	if !connected {
		PrintErrIfNotNil(net.ErrWiFiConnectTimeout)
	}

	println("Connected.")

	time.Sleep(500 * time.Millisecond)
	ip, _, _, err := adaptor.GetIP()
	PrintErrIfNotNil(err)
	println(ip.String())
}

func Run() {
	waitSerial()

	setup()
	connectToAP()

	http.SetBuf(buffer[:])

	// if onboard led is configured as output then the display does not work; D13 is SCK and LED!!!
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
					postDHTReadingToWebhook(temp, hum)
				} else {
					println(fmt.Sprintf("Could not take measurements from the sensor: %s", err.Error()))
				}
			}
		}
		time.Sleep(10 * time.Millisecond)
	}
}

func postDHTReadingToWebhook(temp int16, humidity uint16) {
	rdr := strings.NewReader(fmt.Sprintf(`{"temp":"%02d.%d°C","humidity":"%02d.%d%%"}`, temp/10, temp%10, humidity/10, humidity%10))

	res, err := http.Post(
		webhookSiteUrl,
		"application/json",
		rdr,
	)
	if err != nil {
		println(fmt.Sprintf("Error posting dht reading to webhook: %s\n", err))
	}

	println(fmt.Sprintf("Webhook response status code: %d\n", res.StatusCode))
}
