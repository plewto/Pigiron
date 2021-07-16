package midi

import (
	"fmt"
	"github.com/rakyll/portmidi"
)


// UniversalEvent may represent a portmidi.Event or a meta event.
// If metaType is META_NONE, then portmidi message is used.
// If metaType != META_NONE, the SysEx byte slice of the portmidi message
// is co-opted for the meta message data
//
type UniversalEvent struct {
	time float64
	metaType MetaType
	message portmidi.Event
}

func validateMetaType(mtype MetaType) error {
	var err error
	_, flag := metaMnemonics[mtype]
	if !flag {
		errmsg := "Expected valid MetaType, got 0x%2X"
		err = fmt.Errorf(errmsg, byte(mtype)) 
	}
	return err
}

func MakeMetaEvent(mtype MetaType, data []byte) (*UniversalEvent, error) {
	var err error
	var ue = &UniversalEvent{}
	err = validateMetaType(mtype)
	if err != nil {
		return ue, err
	}
	ue.time = 0.0
	ue.metaType = mtype
	ue.message = portmidi.Event{0, int64(META), 0, 0, data}
	return ue, err
}

func validateMetaTextType(mtype MetaType) error {
	var err error
	_, flag := metaTextTypes[mtype]
	if !flag {
		errmsg := "Expected valid text MetaType, got 0x%2X, using default 0x01"
		err = fmt.Errorf(errmsg, byte(mtype))
	}
	return err
}
	
func MakeMetaTextEvent(mtype MetaType, text string) (*UniversalEvent, error) {
	err := validateMetaTextType(mtype)
	if err != nil {
		mtype = META_TEXT
	}
	return MakeMetaEvent(mtype, []byte(text))
}
	
func MakeSysExEvent(data []byte) (*UniversalEvent, error) {
	var err error
	var ue = &UniversalEvent{}
	ue.time = 0.0
	ue.metaType = META_NONE
	ue.message = portmidi.Event{0, int64(SYSEX), 0, 0, data}
	return ue, err
}

func validateSystemStatus(status StatusByte) error {
	var err error
	_, flag := systemStatusDataCount[status]
	if status == SYSEX || !flag {
		errmsg := "Expected non-sysex system status byte, got 0x%2X"
		err = fmt.Errorf(errmsg, byte(status))
	}
	return err
}

func MakeSystemEvent(status StatusByte) (*UniversalEvent, error) {
	var err = validateSystemStatus(status)
	var ue = &UniversalEvent{}
	if err != nil {
		return ue, err
	}
	ue.time = 0.0
	ue.metaType = META_NONE
	ue.message = portmidi.Event{0, int64(status), 0, 0, []byte{}}
	return ue, err
}
	
func MakeChannelEvent(status StatusByte, ch MIDIChannelNibble, data1 byte, data2 byte) (*UniversalEvent, error) {
	var err error
	var ue = &UniversalEvent{}
	if !isChannelStatus(byte(status)) {
		errmsg := "Expected MIDI channel status byte, got 0x%02X"
		err = fmt.Errorf(errmsg, byte(status))
		return ue, err
	}
	var cstatus int64 = int64(status) | int64(ch)
	ue.time = 0.0
	ue.metaType = META_NONE
	ue.message = portmidi.Event{0, cstatus, int64(data1), int64(data2), []byte{}}
	return ue, err
}
		
	
// TODO: Skeleton only, implement cases.
func (ue *UniversalEvent) String() string {
	acc := fmt.Sprintf("%8.4f ", ue.time)
	s := StatusByte(ue.message.Status)
	switch {
	case isChannelStatus(byte(s)):
		acc += "CHANNEL"
	case s == META:
		acc += "META"
	case s == SYSEX:
		acc += "SYSEX"
	default:
		acc += "SYSTEM"
	}
	return acc
}
	
