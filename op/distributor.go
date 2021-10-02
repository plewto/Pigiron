package op


import (
	gomidi "gitlab.com/gomidi/midi/v2"
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

func (op *Distributor) Send(msg gomidi.Message) {
	st := midi.StatusByte(msg.Data[0])
	if midi.IsChannelStatus(st) {
		cmd := byte(st & 0xF0)
		for _, ci := range op.SelectedChannelIndexes() {
			msg.Data[0] = cmd | byte(ci)
			op.distribute(msg)
		}
	} else {
		op.distribute(msg)
	}
}
		
		


