package op


import (
	"github.com/rakyll/portmidi"
	"github.com/plewto/pigiron/midi"
)

// Distributor is an Operator for changing MIDI channel numbers.
// Incoming channel-message are rebroadcast on each of the distributors
// selected channels.   The original message channel is ignored.
// All non-channel messages are passed unchanged.
//
type Distributor struct {
	baseOperator
}

func newDistributor(name string) *Distributor {
	op := new(Distributor)
	initOperator(&op.baseOperator, "Distributor", name, midi.MultiChannel)
	op.Reset()
	return op
}

func (op *Distributor) Reset() {
	op.DeselectAllChannels()
	op.EnableChannel(midi.MIDIChannel(1), true)
	base := &op.baseOperator
	base.Reset()
}

func (op *Distributor) Send(event portmidi.Event) {
	s := event.Status
	if midi.IsChannelStatus(s) {
		cmd := midi.StatusChannelCommand(s)
		for _, ci := range op.SelectedChannelIndexes() {
			s2 := cmd | int64(ci)
			event.Status = s2
			op.distribute(event)
		}
	} else {
		op.distribute(event)
	}
}
				
