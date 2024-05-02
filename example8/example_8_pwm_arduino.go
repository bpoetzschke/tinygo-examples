//go:build arduino

package example8

import (
	"machine"
	"time"
)

var (
	pwm = machine.Timer2
	pin = machine.D3
)

var period uint64 = (uint64(1e9) * 256 * 8) / uint64(machine.CPUFrequency())

func Run() {
	err := pwm.Configure(machine.PWMConfig{
		Period: period,
	})

	if err != nil {
		for {
			println(err.Error())
			time.Sleep(time.Second)
		}
	}

	ch, err := pwm.Channel(pin)
	if err != nil {
		println(err.Error())
		return
	}

	for {
		for i := uint32(1); i <= pwm.Top(); i++ {
			// This performs a stylish fade-out blink
			pwm.Set(ch, i)
			time.Sleep(4 * time.Millisecond)
		}
		for i := pwm.Top(); i > 0; i-- {
			// This performs a stylish fade-out blink
			pwm.Set(ch, i)
			time.Sleep(4 * time.Millisecond)
		}
	}
}
