package op

import (
	"fmt"
	gomidi "gitlab.com/gomidi/midi/v2"
	goosc "github.com/hypebeast/go-osc/osc"
	"github.com/plewto/pigiron/midi"
)

// SingleChannelFilter is a restricted form of ChannelFilter.
// which only passes messages from a single MIDI channel.
// non-channel messages may optionaly be blocked.
//
type SingleChannelFilter struct {
	baseOperator
	enableSystemEvents bool 
}

func newSingleChannelFilter(name string) *SingleChannelFilter {
	op := new(SingleChannelFilter)
	initOperator(&op.baseOperator, "SingleChannelFilter", name, midi.SingleChannel)
	op.addCommandHandler("q-system-events-enabled", op.remoteQuerySystemEventsEnabled)
	op.addCommandHandler("enable-system-events", op.remoteEnableSystemEvents)
	op.Reset()
	return op
}

func (op *SingleChannelFilter) Reset() {
     op.enableSystemEvents = true
     op.EnableChannel(midi.MIDIChannel(1), true)
     base := &op.baseOperator
     base.Reset()
}

func (op *SingleChannelFilter) Send(msg gomidi.Message) {
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

func (op *SingleChannelFilter) Info() string {
	s := op.commonInfo()
	s += fmt.Sprintf("\tenable system events: %v\n", op.enableSystemEvents)
	return s
}


// osc /pig/op name q-system-events-enabled
// -> bool
//
func (op SingleChannelFilter) remoteQuerySystemEventsEnabled(_ *goosc.Message)([]string, error) {
	var err error
	s := fmt.Sprintf("%v", op.enableSystemEvents)
	return []string{s}, err
}


// osc /pig/op name enable-system-events flag
// -> Ack
//
func (op SingleChannelFilter) remoteEnableSystemEvents(msg *goosc.Message)([]string, error) {
	args, err := ExpectMsg("ssb", msg)
	op.enableSystemEvents = args[2].B
	return empty, err
}
	