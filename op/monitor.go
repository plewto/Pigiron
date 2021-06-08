package op

import (
	"fmt"

	"github.com/rakyll/portmidi"
	"github.com/plewto/pigiron/midi"
)

type Monitor struct {
	baseOperator
}

func newMonitor(name string) *Monitor {
	op := new(Monitor)
	initOperator(&op.baseOperator, "Monitor", name, midi.NoChannel)
	return op
}

func (op *Monitor) Send(event portmidi.Event) {
	if op.MIDIEnabled() {
		op.distribute(event)
		fmt.Printf("%v\n", event)
	}
}
