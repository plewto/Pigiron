package op

import (

	gomidi "gitlab.com/gomidi/midi"
	midi "github.com/plewto/pigiron/midi"
)

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
	return true
}


