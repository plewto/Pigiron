package midi


import (
	"fmt"
)


// MIDIChannel type single MIDI channel, range 1..16 inclusive
//
type MIDIChannel byte


// ValidateMIDIChannel() returns non-nil error if channel is out of bounds.
// Valid channel range is [1,16].
//
func ValidateMIDIChannel(c MIDIChannel) error {
	var err error
	if c < 1 || c > 16 {
		err = fmt.Errorf("Illegal MIDIChannel: %d", c)
	}
	return err
}


func (c MIDIChannel) String() string {
	return fmt.Sprintf("CHAN: %02d", c)
}


// MIDIChannelNibble type, lower four bits of transmitted channel, range 0..15 inclusive
//
type MIDIChannelNibble byte

// ValidateMIDIChannelNibble returns non-nil error if channelNibble is out of bounds.
// Valid Channel Nibble range is [0,15].
//
func ValidateMIDIChannelNibble(ci MIDIChannelNibble) error {
	var err error
	if ci < 0 || ci > 15 {
		err = fmt.Errorf("Illegal MIDIChannelNibble: %d", ci)
	}
	return err
}

func (ci MIDIChannelNibble) String() string {
	return fmt.Sprintf("CHAN_NIB: %02d", ci)
}


// DataNumber indicates 1st or 2nd data byte from MIDI channel message.
//
type DataNumber int

const (
	DATA_1 DataNumber = iota
	DATA_2
)

func (d DataNumber) String() string {
	switch d {
	case DATA_1: return "DATA_1"
	case DATA_2: return "DATA_2"
	default: return "DATA_?"
	}
}


// StatusByte type represents a MIDI status byte.  
// Mask out lower 4-bits for channel status.
//
type StatusByte byte

const (
	NO_STATUS = 0x00    // Special case indicates no status is selected.
	KEYED_STATUS = 0x01 // Special case indicates both NOTE_OFF & NOTE_ON.
	NOTE_OFF StatusByte = 0x80
	NOTE_ON StatusByte = 0x90
	POLY_PRESSURE StatusByte = 0xA0
	CONTROLLER StatusByte = 0xB0
	PROGRAM StatusByte = 0xC0
	CHANNEL_PRESSURE StatusByte = 0xD0
	BEND StatusByte = 0xE0
	META StatusByte = 0xFF
	CLOCK StatusByte = 0xF8
	START StatusByte = 0xFA
	CONTINUE StatusByte = 0xFB
	STOP StatusByte = 0xFC
	ACTIVE_SNESING StatusByte = 0xFE
	SYSEX StatusByte = 0xF0
	END_SYSEX StatusByte = 0xF7
)

var statusMnemonics = map[StatusByte]string {
		NO_STATUS        : "NONE ",
		KEYED_STATUS     : "KEYED",
		NOTE_OFF         : "OFF  ",
		NOTE_ON          : "ON   ",
		POLY_PRESSURE    : "PPRS ",
		CONTROLLER       : "CTRL ",
		PROGRAM          : "PROG ",
		CHANNEL_PRESSURE : "CPRS ",
		BEND             : "BEND ",
		META             : "META ",
		CLOCK            : "CLCK ",
		START            : "STRT ",
		CONTINUE         : "CONT ",
		STOP             : "STOP ",
		ACTIVE_SNESING   : "ASNS ",
		SYSEX            : "SYEX ",
		END_SYSEX        : "EOX  "}


func IsChannelStatus(st StatusByte) bool {
	return st == KEYED_STATUS || (st & 0xF0) < 0xF0
}

func IsMetaStatus(st StatusByte) bool {
	return st == META
}

func IsSystemStatus(st StatusByte) bool {
	return !(IsMetaStatus(st) || IsChannelStatus(st))
}

func IsKeyedStatus(st StatusByte) bool {
	hi := st & 0xF0
	return hi == KEYED_STATUS || hi == NOTE_OFF || hi == NOTE_ON || hi == POLY_PRESSURE
}


func (st StatusByte) String() string {
	var rs string
	if IsChannelStatus(st) {
		rs, _ = statusMnemonics[st & 0xF0]
	} else {
		rs, _ = statusMnemonics[st]
	}
	if rs == "" {
		rs = "?STATUS¿"
	}
	return rs
}

func ChannelMessageDataCount(st StatusByte) int {
	hi := st & 0xF0
	if hi == 0xc0 || hi == 0xd0 {
		return 1
	} else {
		return 2
	}
}


// MetaType type represents a META message type.
//
type MetaType byte

const (
	META_SEQUENCE_NUMBER MetaType = 0x00  
 	META_TEXT MetaType = 0x01
 	META_COPYRIGHT MetaType = 0x02
 	META_TRACK_NAME MetaType = 0x03
 	META_INSTRUMENT_NAME MetaType = 0x04
 	META_LYRIC MetaType = 0x05
 	META_MARKER MetaType = 0x06
 	META_CUEPOINT MetaType = 0x07
 	META_CHANNEL_PREFIX MetaType = 0x20  // obsolete ?
 	META_END_OF_TRACK MetaType = 0x2F
 	META_TEMPO MetaType = 0x51
 	META_SMPTE MetaType = 0x54
 	META_TIME_SIGNATURE MetaType = 0x58
 	META_KEY_SIGNATURE MetaType = 0x59    
 	META_SEQUENCER_EVENT MetaType = 0x7f
	NOT_META MetaType = 0xFF
)

var metaMnemonics = map[MetaType]string {
	META_SEQUENCE_NUMBER : "SEQ_NUMBER",
 	META_TEXT            : "TEXT",
 	META_COPYRIGHT       : "COPYRIGHT",
 	META_TRACK_NAME      : "TRK_NAME", 
 	META_INSTRUMENT_NAME : "INST_NAME",
 	META_LYRIC           : "LYRIC",
 	META_MARKER          : "MARKER",
 	META_CUEPOINT        : "CUEPOINT",
 	META_CHANNEL_PREFIX  : "CHAN PREFIX",
 	META_END_OF_TRACK    : "EOT",
 	META_TEMPO           : "TEMPO",
 	META_SMPTE           : "SMPTE",
 	META_TIME_SIGNATURE  : "TIMESIG",
 	META_KEY_SIGNATURE   : "KEYSIG",
 	META_SEQUENCER_EVENT : "SEQ_EVENT",
	NOT_META             : "ERROR"}

func (mt MetaType) String() string {
	rs, flag := metaMnemonics[mt]
	if !flag {
		rs = "?META¿"
	}
	return rs
}




func StringRepMessage(data []byte) string {
	n := len(data)
	if n < 1 {
		return "MIDI Message data < 0 ?"
	}
	st := StatusByte(data[0])
	acc := st.String()
	if IsChannelStatus(st) {
		c := MIDIChannel(st & 0x0F) + 1
		acc += fmt.Sprintf(" %s", c)
	}
	if n > 1 {
		maxOut := 16
		acc += " DATA: "
		for i, d := range data[1:] {
			acc += fmt.Sprintf("%02X ", d)
			if i > maxOut {
				break
			}
		}
		if len(data)-1 > maxOut {
			acc += "..."
		}
		if IsKeyedStatus(st) {
			acc += fmt.Sprintf("  KEY: %s", KeyName(data[1]))
		}
	}
	return acc
}


