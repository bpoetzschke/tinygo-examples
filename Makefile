flash:
	tinygo flash -target=arduino-nano33

example4:
	tinygo flash -target=arduino-nano33 example_4.go

serial:
	screen /dev/cu.usbmodem14313201 9600