package op

import (
	//"fmt"

	gomidi "gitlab.com/gomidi/midi"
	midi "github.com/plewto/pigiron/midi"
)


// ChannelFilter is an Operator which filters messages by MIDI channel.
//
// It may also filter non-chanel messages.
//
type ChannelFilter struct {
	Operator
	BlockNonChannel bool
}

func makeChannelFilter(name string) *ChannelFilter {
	op := new(ChannelFilter)
	initOperator(&op.Operator, "ChannelFilter", name, midi.MultiChannel)
	op.BlockNonChannel = false
	return op
}

func (op *ChannelFilter) Accept(msg gomidi.Message) bool {
	raw := msg.Raw()	
	status := raw[0]
	result := true
	if isChannelMessage(status) {
		c := int((status & statusMask) + 1)
		result = op.ChannelSelected(c)
	} else {
		result = op.BlockNonChannel
	}
	return result
}


func (op *ChannelFilter) Send(msg gomidi.Message) {
	if op.Accept(msg) && op.MIDIEnabled() {
		op.distribute(msg)
	}
}



