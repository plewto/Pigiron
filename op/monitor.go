package op

import (
	"fmt"

	gomidi "gitlab.com/gomidi/midi"
)

type Monitor struct {
	Operator
	Enabled bool 
}

func makeMonitor(name string) *Monitor {
	op := new(Monitor)
	initOperator(&op.Operator, "Monitor", name, NoChannel)
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
