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
		op.EnableChannel(c, true)
	}
}


func (op *ChannelFilter) Send(event portmidi.Event) {
	if op.MIDIEnabled() {
		s := event.Status
		if midi.IsChannelStatus(s) {
			c := midi.StatusChannelIndex(s) + 1
			if op.ChannelSelected(int(c)) {
				op.distribute(event)
			}
		} else {
			if op.enableSystemEvents {
				op.distribute(event)
			}
		}
	}
}

func (op *ChannelFilter) Info() string {
	s := op.commonInfo()
	s += fmt.Sprintf("\tenable system events: %v\n", op.enableSystemEvents)
	return s
}
