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


var ioCache map[portmidi.DeviceID]*IOWrapper = make(map[portmidi.DeviceID]*IOWrapper)


func ioDeviceCached(id portmidi.DeviceID) bool {
	_, flag := ioCache[id]
	return flag
}

