package smf

import (
	"github.com/rakyll/portmidi"	
)


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
	statusTable = map[StatusByte]string {
		NoteOffStatus: "OFF ",
		NoteOnStatus: "ON  ",
		PolyPressureStatus: "PRES",
		ControllerStatus: "CTRL",
		ProgramStatus: "PROG",
		ChannelPressureStatus: "CPRS",
		BendStatus: "BEND",
		MetaStatus: "META",
		ClockStatus: "CLOK",
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

	channelStatusByteCount = map[StatusByte]int {
		NoteOffStatus: 3,
		NoteOnStatus: 3,
		PolyPressureStatus: 3,
		ControllerStatus: 3,
		ProgramStatus: 2,
		ChannelPressureStatus: 2,
		BendStatus: 3,
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
	
)


func isStatusByte(s byte) bool {
	_, flag := statusTable[StatusByte(s)]
	return flag
}

func isChannelStatus(s byte) bool {
	_, flag := channelStatusTable[StatusByte(s)]
	return flag
}

func isSystemStatus(s byte) bool {
	_, flag := systemStatusTable[StatusByte(s)]
	return flag
}

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

type MIDIMessage interface {
	Status() StatusByte
	Bytes() []byte
	Dump()
	ToPortmidiEvent() (portmidi.Event, error)
}
	
// type RealtimeMIDIMessage interface {
// 	MIDIMessage
//	
// }

