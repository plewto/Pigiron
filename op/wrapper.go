package op

import "github.com/rakyll/portmidi"

type DeviceWrapper interface {
	PigOp
	DeviceID() portmidi.DeviceID
	Stream() portmidi.Stream
	DeviceName() string
	IsInput() bool
	IsOutput() bool
	IsOpen() bool
	Close() error
}

	
func Foo() {
	//fmt.Println("op.Foo()   GET RID OF THIS")
}
