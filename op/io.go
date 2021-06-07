package op

import "github.com/rakyll/portmidi"

type IOWrapper interface {
	Operator
	DeviceID() portmidi.DeviceID
	Stream() *portmidi.Stream
	DeviceName() string
	IsOpen() bool
	IsInput() bool
	IsOutput() bool
}

