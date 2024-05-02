//go:build arduino

package example8

import (
	"machine"
	"time"
)

var (
	pwm = machine.Timer1
	pin = machine.D9
)

var period uint64 = 1e9 / 500

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
		for i := uint32(1); i <= pwm.Top(); i += 100 {
			// This performs a stylish fade-out blink
			pwm.Set(ch, i)
			time.Sleep(4 * time.Millisecond)
		}
		for i := pwm.Top(); i > 0; i -= 100 {
			// This performs a stylish fade-out blink
			pwm.Set(ch, i)
			time.Sleep(4 * time.Millisecond)
		}
	}
}
