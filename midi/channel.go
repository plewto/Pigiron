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




var nsc *NullChannelSelector = new (NullChannelSelector)

type NullChannelSelector struct {}

func NewNullChannelSelector() *NullChannelSelector {
	return nsc
}

func (ns *NullChannelSelector) String() string {
	return "<no channel>"
}

func (ns *NullChannelSelector) ChannelMode() ChannelMode {
	return NoChannel
}

func (ns *NullChannelSelector) EnableChannel(c MIDIChannel, flag bool) error {
	var err error
	return err
}

func (ns *NullChannelSelector) SelectChannel(c MIDIChannel) error {
	var err error
	return err
}

func (ns *NullChannelSelector) SelectedChannelIndexes() []MIDIChannelIndex {
	ary := make([]MIDIChannelIndex, 0, 0)
	return ary
}

func (ns *NullChannelSelector) ChannelIndexSelected(ci MIDIChannelIndex) bool {
	return false
}

func (ns *NullChannelSelector) DeselectAllChannels() {}

func (ns *NullChannelSelector) SelectAllChannels() {}
