//go:build sam && atsamd21

package example9

import (
	"machine"
	"time"
	"tinygo.org/x/drivers/servo"
)

var (
	pwm = machine.TCC0
	pin = machine.D12
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

	serv, err := servo.New(pwm, pin)
	if err != nil {
		for {
			println(err.Error())
			time.Sleep(time.Second)
		}
	}

	start := int16(485)
	stop := int16(2480)
	for {
		for i := start; i <= stop; i++ {
			serv.SetMicroseconds(i)
			time.Sleep(time.Millisecond * 5)
		}
		for i := stop; i >= start; i-- {
			serv.SetMicroseconds(i)
			time.Sleep(time.Millisecond * 5)
		}
		//serv.SetMicroseconds(485)
		//time.Sleep(time.Second * 2)
		////serv.SetMicroseconds(1500)
		////time.Sleep(time.Second * 1)
		////serv.SetMicroseconds(2480)
		////time.Sleep(time.Second * 2)
	}
}
