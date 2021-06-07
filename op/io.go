package op

import "github.com/rakyll/portmidi"

type IOWrapper interface {
	Operator
	DeviceID() portmidi.DeviceID
	Stream() *portmidi.Stream
	DeviceName() string
	IsOutput() bool
}

