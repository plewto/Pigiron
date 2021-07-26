package midi

import "fmt"

// ChannelMode enum indicates the manner in which MIDI channels are selected.
// There are three possible values:
//    NoChannel
//    SingleChannel
//    MultiChannel
//
type ChannelMode int

const (
	NoChannel ChannelMode = iota
	SingleChannel
	MultiChannel
)

func (m ChannelMode) String() string {
	return [...]string{"NoChannel", "SingleChannel", "MultiChannel"}[m]
}

// ChannelSelector interface defines MIDI channel selection.
//
// cs.ChannelMode() ChannelMode
//    ChannelMode function returns the ChannelMode for cs.
//
// cs.EnableChannel(c MIDIChannel, flag bool) error
//    EnableChannel function enables/disables specific MIDI channel.
//
//    For SingleChannel mode the flag argument is ignored and c is
//    set as the current channel.
//
//    For MultiChannel mode, c is enabled/disabled without effecting
//    the state of the other MIDI channels.
//    Returns non-nill error if c is outside the interval [1,16].
//
//    Has no effect for NoChannel mode.
//
// cs.SelectChannel(c MIDIChannel) error
//    SelectChannel is a convenience function identical to
//    cs.EnableChannel(c, true)
//
// cs.SelectedChannelIndexes() []MIDIChannelNibble
//    SelectedChannelIndexes returns a list of all enabled MIDI channels.
//    The channels are specified as transmitted in the interval [0,15].
//    For NoChannel mode the result is an empty list
//    For SingleChannel mode the result will always be a single value.
//    For MultiChannel mode the result is a list between 0 and 16 items.
//
// cs.ChannelIndexSelected(ci MIDIChannelNibble) bool
//    Returns true if the MIDI channel-index (0,15) is enabled.
//    Always returns false for NoChannel Mode.
//
// cs.DeselectAllChannels()
//    Sets all MIDI channels to disabled.
//    Ignored by NoChannel and SingleChannel modes.
//
// cs.SelectAllChannels()
//    Sets all MIDI channels as enabled.
//    Ignored by NoChannel and SingleChannel modes.
//
type ChannelSelector interface {
	ChannelMode() ChannelMode
	EnableChannel(c MIDIChannel, flag bool) error
	SelectChannel(c MIDIChannel) error
	SelectedChannelIndexes() []MIDIChannelNibble
	ChannelIndexSelected(ci MIDIChannelNibble) bool
	DeselectAllChannels()
	SelectAllChannels()
}


/*
** NullChannel singleton
**
*/

var nullSelectorInstance *NullChannelSelector = new (NullChannelSelector)

type NullChannelSelector struct {}

// NewNullChannelSelector() returns a ChannelSelector using NoChannel mode.
// The result is always the same singleton object.
//
func NewNullChannelSelector() *NullChannelSelector {
	return nullSelectorInstance
}

func (ncs *NullChannelSelector) String() string {
	return "<no channel>"
}

func (ncs *NullChannelSelector) ChannelMode() ChannelMode {
	return NoChannel
}

func (ncs *NullChannelSelector) EnableChannel(c MIDIChannel, flag bool) error {
	var err error
	return err
}

func (ncs *NullChannelSelector) SelectChannel(c MIDIChannel) error {
	var err error
	return err
}

func (ncs *NullChannelSelector) SelectedChannelIndexes() []MIDIChannelNibble {
	ary := make([]MIDIChannelNibble, 0, 0)
	return ary
}

func (ncs *NullChannelSelector) ChannelIndexSelected(ci MIDIChannelNibble) bool {
	return false
}

func (ncs *NullChannelSelector) DeselectAllChannels() {}

func (ncs *NullChannelSelector) SelectAllChannels() {}



/*
** SingleChannelSelector implements ChannelSelector for SingleChannel mode.
**
*/

type SingleChannelSelector struct {
	channelIndex MIDIChannelNibble
}

// NewSingleChannelSelector() returns a new ChannelSelector with SingleChannel mode.
//
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
		scs.channelIndex = MIDIChannelNibble(c - 1)
	}
	return err
}

func (scs *SingleChannelSelector) SelectChannel(c MIDIChannel) error {
	return scs.EnableChannel(c, true)
}

func (scs *SingleChannelSelector) SelectedChannelIndexes() []MIDIChannelNibble {
	acc := make([]MIDIChannelNibble, 1, 1)
	acc[0] = scs.channelIndex
	return acc
}

func (scs *SingleChannelSelector) ChannelIndexSelected(ci MIDIChannelNibble) bool {
	return ci == scs.channelIndex
}

func (scs *SingleChannelSelector) DeselectAllChannels() {}

func (scs *SingleChannelSelector) SelectAllChannels() {}


/*
** MultiChannelSelector implements ChannelSelector for MultiChannel mode.
**
*/

type MultiChannelSelector struct {
	flags [16]bool
}


// NewMultChannelSelector() returns new ChannelSelector with MultiChannel mode.
//
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

func (mcs *MultiChannelSelector) SelectedChannelIndexes() []MIDIChannelNibble {
	acc := make([]MIDIChannelNibble, 0, 16)
	for i, flag := range mcs.flags {
		if flag {
			acc = append(acc, MIDIChannelNibble(i))
		}
	}
	return acc
}

func (mcs *MultiChannelSelector) ChannelIndexSelected(ci MIDIChannelNibble) bool {
	err := ValidateMIDIChannelNibble(ci)
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
