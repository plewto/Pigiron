package op


import (
	"github.com/rakyll/portmidi"
	"github.com/plewto/pigiron/midi"
)

// Distributor is an Operator for changing MIDI channels.
// Incoming channel-message are re-broadcast on each selected-channel.
// The original message channel is ignored and all non-channel messages
// are passed unchanged.
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
	s := byte(event.Status)
	if midi.IsChannelStatus(s) {
		cmd := int64(s & 0xF0)
		for _, ci := range op.SelectedChannelIndexes() {
			s2 := cmd | int64(ci)
			event.Status = s2
			op.distribute(event)
		}
	} else {
		op.distribute(event)
	}
}
				
