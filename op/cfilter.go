package op

import (
	"fmt"
	gomidi "gitlab.com/gomidi/midi/v2"
	goosc "github.com/hypebeast/go-osc/osc"
	"github.com/plewto/pigiron/midi"
)

// ChannelFilter is an Operator which selectively blocks MIDI channels.
// Only MIDI messages with channels enabled are allowed through.
// Separately non-channel messages may also be filtered.
//
type ChannelFilter struct {
	baseOperator
	enableSystemEvents bool 
}

func newChannelFilter(name string) *ChannelFilter {
	op := new(ChannelFilter)
	initOperator(&op.baseOperator, "ChannelFilter", name, midi.MultiChannel)
	op.addCommandHandler("q-system-events-enabled", op.remoteQuerySystemEventsEnabled)
	op.addCommandHandler("enable-system-events", op.remoteEnableSystemEvents)
	op.Reset()
	return op
}

func (op *ChannelFilter) Reset() {
	op.enableSystemEvents = true
	for c := 1; c < 17; c++ {
		op.EnableChannel(midi.MIDIChannel(c), true)
	}
	base := &op.baseOperator
	base.Reset()
}

func (op *ChannelFilter) Send(msg gomidi.Message) {
	st := midi.StatusByte(msg.Data[0])
	if midi.IsChannelStatus(st) {
		ci := midi.MIDIChannelNibble(st & 0x0F)
		if op.ChannelIndexSelected(ci) {
			op.distribute(msg)
		}
	} else {
		if op.enableSystemEvents {
			op.distribute(msg)
		}
	}
}

func (op *ChannelFilter) Info() string {
	s := op.commonInfo()
	s += fmt.Sprintf("\tenable system events: %v\n", op.enableSystemEvents)
	return s
}


// osc /pig/op name q-system-events-enabled
// -> bool
//
func (op *ChannelFilter) remoteQuerySystemEventsEnabled(_ *goosc.Message)([]string, error) {
	var err error
	s := fmt.Sprintf("%v", op.enableSystemEvents)
	return []string{s}, err
}


// osc /pig/op name enable-system-events flag
// -> Ack
//
func (op *ChannelFilter) remoteEnableSystemEvents(msg *goosc.Message)([]string, error) {
	args, err := ExpectMsg("ssb", msg)
	op.enableSystemEvents = args[2].B
	return empty, err
}
	
