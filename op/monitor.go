package op

import (
	"fmt"

	gomidi "gitlab.com/gomidi/midi"
	midi "github.com/plewto/pigiron/midi"
)

// Monitor is an Operator to display information about MIDI events.
//
type Monitor struct {
	Operator
	Enabled bool   // If false, do not print MIDI info.
}

func makeMonitor(name string) *Monitor {
	op := new(Monitor)
	initOperator(&op.Operator, "Monitor", name, midi.NoChannel)
	op.Enabled = true
	return op
}

func (op *Monitor) Info() string {
	s := op.commonInfo()
	s += fmt.Sprintf("\tmonitor Enabled : %v\n", op.Enabled)
	return s
}


func (op *Monitor) Send(msg gomidi.Message) {
	if op.MIDIEnabled() {
		op.distribute(msg)
	}
	if op.Enabled {
		fmt.Printf("%s\n", msg)
	}
}
