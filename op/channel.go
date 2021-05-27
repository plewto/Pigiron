package op

/*
 * Provides MIDI channel selection.
*/


import (
	"errors"
	"fmt"
)

// ChannelMode is an enum indicating how an Operator selects MIDI channels.
// NoChannel     --> There is no channel selection.
// SingleChannel --> One, and only one, channel is selected at any time.
// MultiChannel  --> There may be any number of enabled MIDI channels.
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


// ValidateMIDIChannel(channel int) checks for valid channel range.
// Returns nil if channel in interval [1, 16].
// Returns error if channel is invalid.
//
func ValidateMIDIChannel(channel int) error {
	var err error
	if channel < 1 || 16 < channel {
		msg := fmt.Sprintf("Illegal MIDI channel: %d", channel)
		err = errors.New(msg)
	}
	return err
}


// ChannelSelector interface defines Operator MIDI channel selection.
//
// Mode() Returns ChannelMode
//
// EnableChannel(channel int, flag bool) Enables/disables indicated MIDI channel.
//   NoChannel mode, ignored
//   SingleChannel mode, selects channel, flag argument is ignored.
//   MultiChannel mode, Enable/disable channel.
//   Returns error for invalid MIDI channel.
//
// SelectChannel(channel int) syntax sugar for EnableChannel(channel, true).
//   
// SelectedChannels() returns list of currently selected channels.
//
// ChannelSelected(channel int) returns true if channel is currently selected.
// Returns false for out-of bounds channels.
//
// DeselectAllChannels()
//   SingleChannel mode, ignored.
//   MultiChannel mode, sets all channels as deselected.
//
type ChannelSelector interface {
	Mode() ChannelMode
	EnableChannel(channel int, flag bool) error
	SelectChannel(channel int) error
	SelectedChannels() []int
	ChannelSelected(channel int) bool
	DeselectAllChannels()
}



// NullChannelSelector implements place-holder ChannelSelector for NoChannel ChannelMode.
// Do not construct directly, use MakeNullChannelSelector.
//
type NullChannelSelector struct {}

var nsc *NullChannelSelector = new(NullChannelSelector)

// MakeNullChannelSelector returns singleton NullChannelSelector.
//
func MakeNullChannelSelector() *NullChannelSelector {
	return nsc
}

func (n *NullChannelSelector) Mode() ChannelMode {
	return NoChannel
}

func (n *NullChannelSelector) EnableChannel(_ int, _ bool) error {
	return ValidateMIDIChannel(1)
}


func (n *NullChannelSelector) SelectChannel(_ int) error {
	return ValidateMIDIChannel(1)
}


func (n *NullChannelSelector) SelectedChannels() []int {
	return make([]int, 0, 0)
}

func (n *NullChannelSelector) ChannelSelected(_ int) bool {
	return false
}

func (n *NullChannelSelector) DeselectAllChannels(){}

func (n *NullChannelSelector) String() string {
	return "NoChannel"
}



// SingleChannelSelector implements ChannelSelector for Single ChannelMode.
//
type SingleChannelSelector struct {
	channel int
}

// MakeSingleChannelSelector returns new instance of SingleChannelSelector.
//
func MakeSingleChannelSelector() *SingleChannelSelector {
	s := new(SingleChannelSelector)
	s.SelectChannel(1)
	return s
}


func (s *SingleChannelSelector) String() string {
	return fmt.Sprintf("SingleChannel: %2d", s.channel)
}


func (s *SingleChannelSelector) Mode() ChannelMode {
	return SingleChannel
}


func (s *SingleChannelSelector) EnableChannel(channel int, _ bool) error {
	var err error = ValidateMIDIChannel(channel)
	if err == nil {
		s.channel = channel
	}
	return err
}


func (s *SingleChannelSelector) SelectChannel(channel int) error {
	return s.EnableChannel(channel, true)
}

func (s *SingleChannelSelector) SelectedChannels() []int {
	result := make([]int, 0, 1)
	result = append(result, s.channel)
	return result
}

func (s *SingleChannelSelector) ChannelSelected(channel int) bool {
	if channel <= 0 || 16 < channel {
		return false
	} else {
		return channel == s.channel
	}
}

func (s *SingleChannelSelector) DeselectAllChannels() {}





// MultiChannelSelector implements ChannelSelector for MultiChannel ChannelMode.
//
type MultiChannelSelector struct {
	channels [16]bool
}



// MakeMultiChannelSelector returns new instance of MultiChannelSelector.
//
func MakeMultiChannelSelector() *MultiChannelSelector {
	ary := [16]bool{false, false, false, false,
		false, false, false, false,
		false, false, false, false,
		false, false, false, false}
	s := new(MultiChannelSelector)
	s.channels = ary
	return s
}

func (s *MultiChannelSelector) String() string {
	lst := s.SelectedChannels()
	return fmt.Sprintf("MultiChannel: %v", lst)
}


func (s *MultiChannelSelector) Mode() ChannelMode {
	return MultiChannel
}


func (s *MultiChannelSelector) EnableChannel(channel int, flag bool) error {
	var err error = ValidateMIDIChannel(channel)
	if err == nil {
		s.channels[channel-1] = flag
	}
	return err
}


func (s *MultiChannelSelector) SelectChannel(channel int) error {
	return s.EnableChannel(channel, true)
}

func (s *MultiChannelSelector) SelectedChannels() []int {
	result := make([]int, 0, 16)
	for i, flag := range s.channels {
		if flag {
			result = append(result, i+1)
		}
	}	
	return result
}

func (s *MultiChannelSelector) ChannelSelected(channel int) bool {
	if channel <= 0 || 16 < channel {
		return false
	} else {
		return s.channels[channel - 1]
	}
}


func (s *MultiChannelSelector) DeselectAllChannels() {
	for i, _ := range s.channels {
		s.channels[i] = false
	}
}



