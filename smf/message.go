package smf

/*
 * Defines general MIDI messages.
 *
 */

import (
	"github.com/rakyll/portmidi"	
)

// StatusByte type represents MIDI message status byte.
// 
type StatusByte byte

const (
	NoteOffStatus StatusByte = 0x80
	NoteOnStatus StatusByte = 0x90
	PolyPressureStatus StatusByte = 0xA0
	ControllerStatus StatusByte = 0xB0
	ProgramStatus StatusByte = 0xC0
	ChannelPressureStatus StatusByte = 0xD0
	BendStatus StatusByte = 0xE0
	MetaStatus StatusByte = 0xFF
	ClockStatus StatusByte = 0xF8
	StartStatus StatusByte = 0xFA
	ContinueStatus StatusByte = 0xFB
	StopStatus StatusByte = 0xFC
	ActiveSnesingStatus StatusByte = 0xFE
	SysexStatus StatusByte = 0xF0
	EndSysexStatus StatusByte = 0xF7
)

var (
	// statusTable maps StatusByte to string mnemonic.
	//
	statusTable = map[StatusByte]string {
		NoteOffStatus: "OFF ",
		NoteOnStatus: "ON  ",
		PolyPressureStatus: "PRES",
		ControllerStatus: "CTRL",
		ProgramStatus: "PROG",
		ChannelPressureStatus: "CPRS",
		BendStatus: "BEND",
		MetaStatus: "META",
		ClockStatus: "CLK",
		StartStatus: "STRT",
		ContinueStatus: "CONT",
		StopStatus: "STOP",
		ActiveSnesingStatus: "ASEN",
		SysexStatus: "SYSX",
		EndSysexStatus: "EOX ",
	}

	channelStatusTable = map[StatusByte]bool {
		NoteOffStatus: true,
		NoteOnStatus: true,
		PolyPressureStatus: true,
		ControllerStatus: true,
		ProgramStatus: true,
		ChannelPressureStatus: true,
		BendStatus: true,
	}

	channelStatusDataCount = map[StatusByte]int {
		NoteOffStatus: 2,
		NoteOnStatus: 2,
		PolyPressureStatus: 2,
		ControllerStatus: 2,
		ProgramStatus: 1,
		ChannelPressureStatus: 1,
		BendStatus: 2,
	}
	
	systemStatusTable = map[StatusByte]bool {
		ClockStatus: true,
		StartStatus: true,
		ContinueStatus: true,
		StopStatus: true,
		ActiveSnesingStatus: true,
		SysexStatus: true,
		EndSysexStatus: true,
	}

	// -1 indicates indeterminate
	systemStatusDataCount = map[StatusByte]int {
		ClockStatus: 0,
		StartStatus: 0,
		ContinueStatus: 0,
		StopStatus: 0,
		ActiveSnesingStatus: 0,
		SysexStatus: -1,
		EndSysexStatus: 0,
	}
	
)

// isStatusByte returns true iff argument is a valid MIDI status byte.
//
func isStatusByte(s byte) bool {
	_, flag := statusTable[StatusByte(s)]
	return flag
}

// isChannelStatus returns true iff argument is a channel-based status byte.
//
func isChannelStatus(s byte) bool {
	_, flag := channelStatusTable[StatusByte(s)]
	return flag
}

// isSystemStatus returns true iff argument is a system-message status byte.
//
func isSystemStatus(s byte) bool {
	_, flag := systemStatusTable[StatusByte(s)]
	return flag
}

// isMetaStatus returns true iff argument is the meta status byte.
//
func isMetaStatus(s byte) bool {
	return StatusByte(s) == MetaStatus
}


func (s StatusByte) String() string {
	c, flag := statusTable[s]
	if !flag {
		c = "????"
	}
	return c
}


// MIDIMessage interface universal MIDI message methods.
// There are at least three implementing structure:
//   ChannelMessage
//   SystemMessage
//   MetaMessage
//
type MIDIMessage interface {
	Status() StatusByte
	Bytes() []byte
	Dump()
	ToPortmidiEvent() (portmidi.Event, error)
}
	
