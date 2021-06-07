package op

import "github.com/rakyll/portmidi"

type DeviceWrapper interface {
	Operator
	DeviceID() portmidi.DeviceID
	Stream() *portmidi.Stream
	DeviceName() string
	IsOpen() bool
	Close() error
}

func Foo() {
}
