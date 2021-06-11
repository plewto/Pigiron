package op

import (
	"fmt"

	"github.com/rakyll/portmidi"
	"github.com/plewto/pigiron/midi"
)

type ChannelFilter struct {
	baseOperator
	enableSystemEvents bool 
}

func newChannelFilter(name string) *ChannelFilter {
	op := new(ChannelFilter)
	initOperator(&op.baseOperator, "ChannelFilter", name, midi.MultiChannel)
	op.Reset()
	return op
}

func (op *ChannelFilter) Reset() {
	op.enableSystemEvents = true
	for c := 1; c < 17; c++ {
		op.EnableChannel(midi.MIDIChannel(c), true)
	}
}

func (op *ChannelFilter) Send(event portmidi.Event) {
	s := event.Status
	if midi.IsChannelStatus(s) {
		ci := midi.StatusChannelIndex(s)
		if op.ChannelIndexSelected(ci) {
			op.distribute(event)
		}
	} else {
		if op.enableSystemEvents {
			op.distribute(event)
		}
	}
}

func (op *ChannelFilter) Info() string {
	s := op.commonInfo()
	s += fmt.Sprintf("\tenable system events: %v\n", op.enableSystemEvents)
	return s
}
