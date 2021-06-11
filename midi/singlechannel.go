package midi

import (
	"fmt"
)

type SingleChannelSelector struct {
	channelIndex MIDIChannelIndex
}

func NewSingleChannelSelector() *SingleChannelSelector {
	scs := new(SingleChannelSelector)
	scs.channelIndex = 0
	return scs
}


func (scs *SingleChannelSelector) String() string {
	return fmt.Sprintf("MIDI Channel: %2d", scs.channelIndex + 1)
}

func (scs *SingleChannelSelector) ChannelMode() ChannelMode {
	return SingleChannel
}

func (scs *SingleChannelSelector) EnableChannel(c MIDIChannel, _ bool) error {
	err := ValidateMIDIChannel(c)
	if err == nil {
		scs.channelIndex = MIDIChannelIndex(c - 1)
	}
	return err
}

func (scs *SingleChannelSelector) SelectChannel(c MIDIChannel) error {
	return scs.EnableChannel(c, true)
}

func (scs *SingleChannelSelector) SelectedChannelIndexes() []MIDIChannelIndex {
	acc := make([]MIDIChannelIndex, 1, 1)
	acc[0] = scs.channelIndex
	return acc
}

func (scs *SingleChannelSelector) ChannelIndexSelected(ci MIDIChannelIndex) bool {
	return ci == scs.channelIndex
}

func (scs *SingleChannelSelector) DeselectAllChannels() {}

func (scs *SingleChannelSelector) SelectAllChannels() {}

