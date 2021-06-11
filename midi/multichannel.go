package midi

import (
	"fmt"
)


type MultiChannelSelector struct {
	flags [16]bool
}

func NewMultiChannelSelector() *MultiChannelSelector {
	ary := [16]bool{false, false, false, false,
		false, false, false, false,
		false, false, false, false,
		false, false, false, false}
	s := new(MultiChannelSelector)
	s.flags = ary
	return s
}


func (mcs *MultiChannelSelector) String() string {
	s := "MIDI Channels: "
	for i, flag := range mcs.flags {
		if flag {
			s += fmt.Sprintf("%2d ", i+1)
		}
	}
	return s
}
				
	
func (mcs *MultiChannelSelector) ChannelMode() ChannelMode {
	return MultiChannel
}

func (mcs *MultiChannelSelector) EnableChannel(c MIDIChannel, flag bool) error {
	err := ValidateMIDIChannel(c)
	if err == nil {
		mcs.flags[c-1] = flag
	}
	return err
}

func (mcs *MultiChannelSelector) SelectChannel(c MIDIChannel) error {
	return mcs.EnableChannel(c, true)
}

func (mcs *MultiChannelSelector) SelectedChannelIndexes() []MIDIChannelIndex {
	acc := make([]MIDIChannelIndex, 0, 16)
	for i, flag := range mcs.flags {
		if flag {
			acc = append(acc, MIDIChannelIndex(i))
		}
	}
	return acc
}

func (mcs *MultiChannelSelector) ChannelIndexSelected(ci MIDIChannelIndex) bool {
	err := ValidateMIDIChannelIndex(ci)
	if err == nil {
		return mcs.flags[ci]
	} else {
		fmt.Println(err)
		return false
	}
}

func (mcs *MultiChannelSelector) DeselectAllChannels() {
	for i:=0; i<16; i++ {
		mcs.flags[i] = false
	}
}

func (mcs *MultiChannelSelector) SelectAllChannels() {
	for i:=0; i<16; i++ {
		mcs.flags[i] = true
	}
}
