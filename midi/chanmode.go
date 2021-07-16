package midi

import (
	"errors"
	"fmt"
)
	

type MIDIChannel int       // 1..16
type MIDIChannelIndex int  // 0..15

type ChannelMode int

const (
	NoChannel ChannelMode = iota
	SingleChannel
	MultiChannel
)

func (m ChannelMode) String() string {
	return [...]string{"NoChannel", "SingleChannel", "MultiChannel"}[m]
}

func ValidateMIDIChannel(c MIDIChannel) error {
	var err error
	if c < 1 || c > 16 {
		msg := fmt.Sprintf("Illegal MIDIChannel: %d", c)
		err = errors.New(msg)
	}
	return err
}

func ValidateMIDIChannelIndex(ci MIDIChannelIndex) error {
	var err error
	if ci < 0 || ci > 15 {
		msg := fmt.Sprintf("Illegal MIDIChannelIndex: %d", ci)
		err = errors.New(msg)
	}
	return err
}


type ChannelSelector interface {
	ChannelMode() ChannelMode
	EnableChannel(c MIDIChannel, flag bool) error
	SelectChannel(c MIDIChannel) error
	SelectedChannelIndexes() []MIDIChannelIndex
	ChannelIndexSelected(ci MIDIChannelIndex) bool
	DeselectAllChannels()
	SelectAllChannels()
}


/*
 * nullChannelSeletor singleton, default instance for ChannelSelector
 *
 */

var nullSelectorInstance *NullChannelSelector = new (NullChannelSelector)

type NullChannelSelector struct {}

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

func (ncs *NullChannelSelector) SelectedChannelIndexes() []MIDIChannelIndex {
	ary := make([]MIDIChannelIndex, 0, 0)
	return ary
}

func (ncs *NullChannelSelector) ChannelIndexSelected(ci MIDIChannelIndex) bool {
	return false
}

func (ncs *NullChannelSelector) DeselectAllChannels() {}

func (ncs *NullChannelSelector) SelectAllChannels() {}



/*
 * SingleChannelSelector struct selectes one, and only one, MIDI channel at a time.
 *
 */

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

/*
 * MultiChannelSelector selects an arbitary set of MIDI channels.
 *
*/

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
