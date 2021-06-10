package op


import (
	"github.com/rakyll/portmidi"
	"github.com/plewto/pigiron/midi"
)

type Distributor struct {
	baseOperator
}

func newDistributor(name string) *Distributor {
	op := new(Distributor)
	initOperator(&op.baseOperator, "ChannelFilter", name, midi.MultiChannel)
	op.Reset()
	return op
}

func (op *Distributor) Reset() {
	for c := 1; c < 17; c++ {
		op.EnableChannel(c, true)
	}
}

func (op *Distributor) Send(event portmidi.Event) {
	if op.MIDIEnabled() {
		s := event.Status
		if midi.IsChannelStatus(s) {
			cmd := midi.StatusChannelCommand(s)
			for _, c := range op.SelectedChannels() {
				ci := int64(c - 1)
				s2 := cmd | ci
				event.Status = s2
				op.distribute(event)
			}
		} else {
			op.distribute(event)
		}
	}
}
	
