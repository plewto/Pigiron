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
	op.distribute(event)
	fmt.Printf(midi.EventToString(event))
	fmt.Println()
}
