package example7

import (
	"fmt"
	"image/color"
	"machine"
	"strings"
	"time"

	"tinygo.org/x/drivers/dht"
	"tinygo.org/x/drivers/net"
	"tinygo.org/x/drivers/net/http"
	"tinygo.org/x/drivers/waveshare-epd/epd2in13x"
	"tinygo.org/x/drivers/wifinina"
	"tinygo.org/x/tinydraw"
	"tinygo.org/x/tinyfont"
	"tinygo.org/x/tinyfont/freesans"
)

var (
	buffer [1024]byte

	// these are the default pins for the Arduino Nano33 IoT.
	spi = machine.NINA_SPI

	// this is the ESP chip that has the WIFININA firmware flashed on it
	adaptor *wifinina.Device

	input           = machine.D2
	dhtInput        = machine.D3
	led             = machine.D4
	lastButtonState = false
	display         epd2in13x.Device
	colorBlack      = color.RGBA{1, 1, 1, 255}
	colorWhite      = color.RGBA{0, 0, 0, 255}
	colored         = color.RGBA{255, 0, 0, 255}
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

func setupDisplay() {
	machine.SPI0.Configure(machine.SPIConfig{
		Frequency: 8000000,
		Mode:      0,
	})

	display = epd2in13x.New(machine.SPI0, machine.D10, machine.D9, machine.D8, machine.D7)
	display.Configure(epd2in13x.Config{})
}

func Run() {
	waitSerial()

	setup()
	connectToAP()

	setupDisplay()

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
					displayDHTReading(temp, hum)
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

func displayDHTReading(temp int16, humidity uint16) {

	display.ClearBuffer()
	display.ClearDisplay()

	tinydraw.FilledRectangle(&display, 10, 0, 30, 212, colored)

	tinyfont.WriteLineRotated(&display, &freesans.Regular9pt7b, 50, 10, fmt.Sprintf("Temperature: %02d.%d°C", temp/10, temp%10), colorBlack, tinyfont.ROTATION_90)
	tinyfont.WriteLineRotated(&display, &freesans.Regular9pt7b, 20, 10, fmt.Sprintf("Humidity: %02d.%d%%", humidity/10, humidity%10), colorWhite, tinyfont.ROTATION_90)

	display.Display()
	display.WaitUntilIdle()

}
