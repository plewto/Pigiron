package midi

import (
	"fmt"
)

// MIDIChannel type is a single MIDI channel with range 0..16
//
type MIDIChannel byte 


// MIDIChannelNibble type is the lower 4-bit binary representation of a MIDI channel
// range 0..15
//
type MIDIChannelNibble byte


func ValidateMIDIChannel(c MIDIChannel) error {
	var err error
	if c < 1 || c > 16 {
		err = fmt.Errorf("Illegal MIDIChannel: %d", c)
	}
	return err
}

func ValidateMIDIChannelNibble(ci MIDIChannelNibble) error {
	var err error
	if ci < 0 || ci > 15 {
		err = fmt.Errorf("Illegal MIDIChannelNibble: %d", ci)
	}
	return err
}


type StatusByte byte
type MetaType byte


const (
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

	META_SEQUENCE_NUMBER MetaType = 0x00  
 	META_TEXT MetaType = 0x01
 	META_COPYRIGHT MetaType = 0x02
 	META_TRACK_NAME MetaType = 0x03
 	META_INSTRUMENT_NAME MetaType = 0x04
 	META_LYRIC MetaType = 0x05
 	META_MARKER MetaType = 0x06
 	META_CUEPOINT MetaType = 0x07
 	META_CHANNEL_PREFIX MetaType = 0x20
 	META_END_OF_TRACK MetaType = 0x2F
 	META_TEMPO MetaType = 0x51
 	META_SMPTE MetaType = 0x54
 	META_TIME_SIGNATURE MetaType = 0x58
 	META_KEY_SIGNATURE MetaType = 0x59    
 	META_SEQUENCER_EVENT MetaType = 0x7f
	META_NONE MetaType = 0xFF
)

var (
	// Maps StatusByte to string mnemonic.
	//
	statusMnemonics = map[StatusByte]string {
		NOTE_OFF: "OFF ",
		NOTE_ON: "ON  ",
		POLY_PRESSURE: "PRES",
		CONTROLLER: "CTRL",
		PROGRAM: "PROG",
		CHANNEL_PRESSURE: "CPRS",
		BEND: "BEND",
		META: "META",
		CLOCK: "CLK",
		START: "STRT",
		CONTINUE: "CONT",
		STOP: "STOP",
		ACTIVE_SNESING: "ASEN",
		SYSEX: "SYSX",
		END_SYSEX: "EOX ",
	}

	channelStatusDataCount = map[StatusByte]int {
		NOTE_OFF: 2,
		NOTE_ON: 2,
		POLY_PRESSURE: 2,
		CONTROLLER: 2,
		PROGRAM: 1,
		CHANNEL_PRESSURE: 1,
		BEND: 2,
	}

	// -1 indicates indeterminate
	systemStatusDataCount = map[StatusByte]int {
		CLOCK: 0,
		START: 0,
		CONTINUE: 0,
		STOP: 0,
		ACTIVE_SNESING: 0,
		SYSEX: -1,
		END_SYSEX: 0,
	}


	// metaTypeTable maps MetaType to string mnemonic.
	//
	metaMnemonics map[MetaType]string = map[MetaType]string {
		META_SEQUENCE_NUMBER: "SeqNum ",
		META_TEXT:            "Text   ",
		META_COPYRIGHT:       "CpyRite",
		META_TRACK_NAME:      "TrkName",
		META_INSTRUMENT_NAME: "InsName",
		META_LYRIC:           "Lyric  ",
		META_MARKER:          "Marker ",
		META_CUEPOINT:        "Cue    ",
		META_CHANNEL_PREFIX:  "ChanPre",
		META_END_OF_TRACK:    "EOT    ",
		META_TEMPO:           "Tempo  ",
		META_SMPTE:           "SMPTE  ",
		META_TIME_SIGNATURE:  "TSig   ",
		META_KEY_SIGNATURE:   "KSig   ",
		META_SEQUENCER_EVENT: "SeqEvnt",
		META_NONE:            "<none> ",
	}

	metaTextTypes map[MetaType]bool = map[MetaType]bool{
		META_TEXT: true,
		META_COPYRIGHT: true,
		META_TRACK_NAME: true,
		META_INSTRUMENT_NAME: true,
		META_LYRIC: true,
		META_MARKER: true,
		META_CUEPOINT: true,
	}
)



// NOTE: Use care with channel status, the lower 4-bits must be masked out
// 
func isStatusByte(s byte) bool {
	_, flag := statusMnemonics[StatusByte(s)]
	return flag
}

func isChannelStatus(s byte) bool {
	cs := byte(s & 0x0F)
	return isStatusByte(cs)
}


func isKeyedStatus(s byte) bool {
	sb := StatusByte(s)
	return sb == NOTE_OFF || sb == NOTE_ON || sb == POLY_PRESSURE
}


// isSystemStatus returns true iff argument is a system-message status byte.
//
func isSystemStatus(s byte) bool {
	_, flag := systemStatusDataCount[StatusByte(s)]
	return flag
}

// isMetaStatus returns true iff argument is the meta status byte.
//
func isMetaStatus(s byte) bool {
	return StatusByte(s) == META
}




func (s StatusByte) String() string {
	c, flag := statusMnemonics[s]
	if !flag {
		c = "?   "
	}
	return c
}

func (mt MetaType) String() string {
	c, flag := metaMnemonics[mt]
	if !flag {
		c = "?      "
	}
	return c
}


