package midi

import (
	"fmt"
	"math"
	"github.com/rakyll/portmidi"
)


// UniversalEvent struct may represent a portmidi.Event or a meta event.
// If metaType is NOT_META, then portmidi message is used.
// If metaType != NOT_META, the SysEx byte slice of the portmidi message
// is co-opted for the meta message data.
//
type UniversalEvent struct {
	deltaTime int
	metaType MetaType
	message portmidi.Event
}

// ue.PortmidiEvent() returns portmidi.Event of UniversalEvent.
//
// This method should only be called if ue.MetaType() returns NOT_META (0xFF).
//
func (ue *UniversalEvent) PortmidiEvent() portmidi.Event {
	return ue.message
}

// ue.DeltaTime() returns the time difference for this event and the previous event.
// Time is specified as MIDI clock ticks.
//
func (ue *UniversalEvent) DeltaTime() int {
	return ue.deltaTime
}

// ue.IsMetaEvent() returns true if the UniversalEvent contains a Meta Message.
//
func (ue *UniversalEvent) IsMetaEvent() bool {
	return ue.metaType != NOT_META
}


// ue.IsChannelEvent() returns true if the UniversalEvent contains a MIDI channel message.
//
func (ue *UniversalEvent) IsChannelEvent() bool {
	s := byte(ue.message.Status)
	return isChannelStatus(s)
}

// ue.IsSystemEvent() returns true if the UniversalEvent contains a MIDI system real-time  message.
//
func (ue *UniversalEvent) IsSystemEvent() bool {
	s := byte(ue.message.Status)
	return isSystemStatus(s)
}

// us.MetaType() returns the type of meta message a UniversalMessage contains.
// A NOT_META (0xFF) indicates this is not a meta event.
//
func (ue *UniversalEvent) MetaType() MetaType {
	return ue.metaType
}

// validateMetaType() returns error if argument is not a valid MetaType.
//
func validateMetaType(mtype MetaType) error {
	var err error
	_, flag := metaMnemonics[mtype]
	if !flag {
		errmsg := "Expected valid MetaType, got 0x%2X"
		err = fmt.Errorf(errmsg, byte(mtype)) 
	}
	return err
}

// MakeMetaEvent() creates a new UniversalEvent with meta message
// mtype - The MetaType
// data  - data bytes.
//
// Returns non-nil error if mtype is not a valid MetaType.
//
func MakeMetaEvent(mtype MetaType, data []byte) (*UniversalEvent, error) {
	var err error
	var ue = &UniversalEvent{}
	err = validateMetaType(mtype)
	if err != nil {
		return ue, err
	}
	ue.deltaTime = 0.0
	ue.metaType = mtype
	ue.message = portmidi.Event{0, int64(META), 0, 0, data}
	return ue, err
}

// validateMetaTextType() returns error if argument is not a valid META text type.
//
func validateMetaTextType(mtype MetaType) error {
	var err error
	_, flag := metaTextTypes[mtype]
	if !flag {
		errmsg := "Expected valid text MetaType, got 0x%2X, using default 0x01"
		err = fmt.Errorf(errmsg, byte(mtype))
	}
	return err
}

// MakeMetaTextEvent() creates new UniversalEvent for text event.
// mtype - the meta text type.  If mtype id not a valid text type, type TEXT is used.
// text  - 
//
// Returns
//  1. *UniversalEvent (never nil)
//  2. non-nil error if mtype is not a valid text type.
//
func MakeMetaTextEvent(mtype MetaType, text string) (*UniversalEvent, error) {
	err := validateMetaTextType(mtype)
	if err != nil {
		mtype = META_TEXT
	}
	return MakeMetaEvent(mtype, []byte(text))
}

// MakeSysExEvent() creates new UniversalEvent for System Exclusive message.
// data - the message bytes.
//
func MakeSysExEvent(data []byte) (*UniversalEvent, error) {
	var err error
	var ue = &UniversalEvent{}
	ue.deltaTime = 0.0
	ue.metaType = NOT_META
	ue.message = portmidi.Event{0, int64(SYSEX), 0, 0, data}
	return ue, err
}


// validateSystemStatus() returns error if argument is not a valid system status byte.
//
func validateSystemStatus(status StatusByte) error {
	var err error
	_, flag := systemStatusDataCount[status]
	if status == SYSEX || !flag {
		errmsg := "Expected non-sysex system status byte, got 0x%2X"
		err = fmt.Errorf(errmsg, byte(status))
	}
	return err
}

// MakeSystemEvent() creates new UniversalEvent for MIDI system events.
// Do not use for SYSEX events, use the MakeSysExEvent instead.
//
// Returns
//   1. *UniversalEvent
//   2. err if status is not a valid (non SYSEX) system byte.
//
func MakeSystemEvent(status StatusByte) (*UniversalEvent, error) {
	var err = validateSystemStatus(status)
	var ue = &UniversalEvent{}
	if err != nil {
		return ue, err
	}
	ue.deltaTime = 0.0
	ue.metaType = NOT_META
	ue.message = portmidi.Event{0, int64(status), 0, 0, []byte{}}
	return ue, err
}


// MakeChannelEvent() creates new UniversalEvent for MIDI channel messages.
//
//   status - a channel message status, the lower 4-bits should be masked out.
//   ch - Channel number is lower 4-bit nibble, 0 <= ch <= 15.
//   data1, data2, the data bytes.
//
// NOTE_ON events with data2 = 0 are converted to NOTE_OFF.
//
// Returns
//   1. *UniversalEvent
//   2. non-nil error if status is not a valid channel status byte.
//
func MakeChannelEvent(status StatusByte, ch MIDIChannelNibble, data1 byte, data2 byte) (*UniversalEvent, error) {
	var err error
	var ue = &UniversalEvent{}
	if !isChannelStatus(byte(status)) {
		errmsg := "Expected MIDI channel status byte, got 0x%02X"
		err = fmt.Errorf(errmsg, byte(status))
		return ue, err
	}
	if status == NOTE_ON && data2 == 0 {
		status = NOTE_OFF
	}
	var cstatus int64 = int64(status) | int64(ch)
	ue.deltaTime = 0.0
	ue.metaType = NOT_META
	ue.message = portmidi.Event{0, cstatus, int64(data1), int64(data2), []byte{}}
	return ue, err
}

// MakeControllerEvent() creates UniversalEvent for MIDI control change message.
//
//  ch - MIDI channel index [0, 15]
//  controllerNumber - int [0, 127]
//  value - int [0, 127]
//
// All arguments are masked to be in valid range.
//
// Returns *UniversalEvent.
//
func MakeControllerEvent(ch MIDIChannelNibble, controllerNumber byte, value byte) *UniversalEvent {
	controllerNumber = controllerNumber & 0x7F
	value = value & 0x7F
	ev, _ := MakeChannelEvent(CONTROLLER, ch, controllerNumber, value)
	return ev
}

// bytesToString() returns hex-string representation of byte slice.
// The output may be truncated.
//
func (ue *UniversalEvent) bytesToString() string {
	maxOut := 8
	data := ue.message.SysEx
	acc := "["
	for i := 0; i < len(data) && i < maxOut; i++ {
		acc += fmt.Sprintf("%02X ", data[i])
	}
	if len(data) > maxOut {
		acc += fmt.Sprintf("... +%d more", len(data) - maxOut)
	}
	acc += "]"
	return acc
}

// MetaData() returns UniversalEvent meta data bytes.
// Returns error if the UniversalEvent is not a valid meta event.
//
func (ue *UniversalEvent) MetaData() ([]byte, error) {
	var err error
	var bytes []byte
	if !isMetaType(byte(ue.metaType)) {
		errmsg := "Non meta event passed to MetaData, status = 0x%02X, mtype = 0x%02X"
		err = fmt.Errorf(errmsg, ue.message.Status, byte(ue.metaType))
		return bytes, err
	}
	return ue.message.SysEx, err
}

// metaTempoMicroSeconds() returns the micro-second value of a META_TEMPO event.
// Returns error if event is not a tempo event or it is malformed.
//
func (ue *UniversalEvent) metaTempoMicroSeconds() (int64, error) {
	var err error
	st := StatusByte(ue.message.Status)
	mt := ue.metaType
	if st != META || mt != META_TEMPO {
		errmsg := "Non meta-tempo event passed to metaTempoMicroSeconds"
		errmsg += "status = %02X, mtype = %0x2X"
		err = fmt.Errorf(errmsg, byte(st), byte(mt))
		return 0, err
	}
	var acc int64 = 0
	data := ue.message.SysEx
	if len(data) != 3 {
		errmsg := "Malformed meta tempo event, expected 3 data bytes, got %d"
		err = fmt.Errorf(errmsg, len(data))
		return 0, err
	}
	for i, shift := 0, 16; i < 3; i, shift = i+1, shift-8 {
		acc += int64(data[i]) << shift
	}
	return acc, err
}

// MetaTempoBPM() returns the tempo in BPM for a META_TEMPO event.
// Returns error if event is not a tempo event or it is malformed.
// If an error is detected, the returned tempo defaults to 60.0
//
func (ue *UniversalEvent) MetaTempoBPM() (float64, error) {
	μ, err := ue.metaTempoMicroSeconds()
	if err != nil {
		return 60.0, err
	}
	if μ == 0 {
		errmsg := "Malformed meta tempo event, μ = 0"
		err = fmt.Errorf(errmsg)
		return 60.0, err
	}
	var k float64 = 60000000
	return k/float64(μ), err
}


// metaTimesig() returns the numerator and denominator for a META_TIME_SIGNATURE event.
// Returns error if event is not a time-signature or it is malformed.
// If an error is detected the resulting values default to 4/4.
//
func (ue *UniversalEvent) metaTimesig() (num byte, den byte, err error) {
	st := StatusByte(ue.message.Status)
	mt := ue.metaType
	if st != META || mt != META_TIME_SIGNATURE {
		errmsg := "Non time-signature event passed to metaTimesig"
		errmsg += "status = 0x%02X, mtype = 0x%02X, using default 4/4"
		err = fmt.Errorf(errmsg, byte(st), byte(mt))
		return 4, 4, err
	}
	bytes := ue.message.SysEx
	if len(bytes) != 4 {
		errmsg := "Malformed meta time-signature, expected 4 bytes got %d, "
		errmsg += "using default 4/4"
		err = fmt.Errorf(errmsg, len(bytes))
		return 4, 4, err
	}
	num = bytes[0]
	den = byte(math.Pow(2, float64(bytes[1])))
	return num, den, err
}

func (ue *UniversalEvent) String() string {
	s := StatusByte(ue.message.Status)
	acc := fmt.Sprintf("[∆t %8d] ", ue.deltaTime)
	switch {
	case isChannelStatus(byte(s)):
		c := byte(s & 0x0F) + 1
		s := StatusByte(byte(s) & 0xF0)
		d1, d2 := ue.message.Data1, ue.message.Data2
		acc += fmt.Sprintf("%s  chan: %02d  data: %02X %02X", s, c, byte(d1), byte(d2))
		if isKeyedStatus(byte(s)) {
			acc += fmt.Sprintf("  KEY: %s", KeyName(d1))
		}
	case s == META:
		mt := ue.metaType
		acc += fmt.Sprintf("META %s ", ue.metaType)
		acc += ue.bytesToString()
		switch {
		case IsMetaTextType(byte(mt)):
			acc += fmt.Sprintf("  \"%s\"", string(ue.message.SysEx))
		case mt == META_TEMPO:
			μ, _ := ue.metaTempoMicroSeconds()
			bpm, err := ue.MetaTempoBPM()
			if err != nil {
				acc += " <malformed tempo event, using 60.0>"
			} else {
				acc += fmt.Sprintf(" %d μsec, %8.4f BPM", μ, bpm)
			}
		case mt == META_SMPTE:
			// pass
		case mt == META_TIME_SIGNATURE:
			num, den, err := ue.metaTimesig()
			if err != nil {
				acc += " <malformed timesig, using default 4/4>"
			} else {
				acc += fmt.Sprintf(" %d/%d", num, den)
			}
					
		default:
			// ignore
		}
	case s == SYSEX:
		acc += fmt.Sprintf("SYSEX ")
		acc += ue.bytesToString()
	default:
		// assume non-sysex system event
		acc += fmt.Sprintf("%s", s)
	}
	return acc
}

// BytesToEvents() converts byte slice to list of MIDI events.
// Do not use running-status
// 
func BytesToEvents(bytes []byte)(events []*UniversalEvent, err error) {
	events = make([]*UniversalEvent, 0, 8)
	for len(bytes) > 0 {
		st := bytes[0]
		switch {
		case isChannelStatus(byte(st)):
			var event *UniversalEvent
			var count int
			cmd, ci := StatusByte(st & 0xF0), MIDIChannelNibble(st & 0x0F)
			count, _ = channelStatusDataCount[cmd]
			if len(bytes) < count + 1 {
				errmsg := "Expected %d data bytes for %s, got %d"
				err = fmt.Errorf(errmsg, count, StatusByte(st), count-1)
				return
			}
			var d1, d2 byte = 0, 0
			d1 = bytes[1]
			if count == 2 {
				d2 = bytes[2]
				bytes = bytes[3:]
			} else {
				bytes = bytes[2:]
			}
			event, err = MakeChannelEvent(cmd, ci, d1, d2)
			if err != nil {
				errmsg := "BytesToEvents could not create ChannelMessage: status = 0x%X"
				errmsg += "Original error was: %s"
				err = fmt.Errorf(errmsg, st, err)
				return
			}
			events = append(events, event)
		case st == byte(SYSEX):
			var event *UniversalEvent
			var pointer = 1
			var by byte = 0
			for pointer < len(bytes) && by != byte(END_SYSEX) {
				pointer++
			}
			data := bytes[1:pointer]
			event, err = MakeSysExEvent(data)
			if err != nil {
				errmsg := "BytesToEvents could not create SysexMessage: status = 0x%X"
				errmsg += "Original error was: %s"
				err = fmt.Errorf(errmsg, st, err)
				return
			}
			events = append(events, event)
			bytes = bytes[pointer:]
		case isStatusByte(st):
			var event *UniversalEvent
			event, err = MakeSystemEvent(StatusByte(st))
			if err != nil {
				errmsg := "BytesToEvents could not create System: status = 0x%X"
				errmsg += "Original error was: %s"
				err = fmt.Errorf(errmsg, st, err)
				return
			}
			events = append(events, event)
			bytes = bytes[1:]
		default:
			errmsg := "midi.ByteToEvent Unhandled switch case, status was 0x%X"
			err = fmt.Errorf(errmsg, int(st))
			return
		}
	}
	
	return
}
