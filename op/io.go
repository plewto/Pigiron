package op

import "github.com/rakyll/portmidi"


// IOWrapper interface defines common behavior for MIDIInput and MIDIOutput Operators.
//
type IOWrapper interface {
	Operator
	DeviceID() portmidi.DeviceID
	Stream() *portmidi.Stream
	DeviceName() string
	IsOutput() bool
}

