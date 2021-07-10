package smf

import (
	"fmt"
	"github.com/rakyll/portmidi"
)

// Implements MIDIMessage
//
type SystemMessage struct {
	bytes []byte
}

func newSystemMessage(bytes []byte) (*SystemMessage, error) {
	var err error
	var sys *SystemMessage
	if len(bytes) < 1 {
		msg := "newSystemMessage, bytes slice is empty, must have at least a status byte"
		err = fmt.Errorf(msg)
		return sys, err
	}
	var status = StatusByte(bytes[0])
	if !isSystemStatus(byte(status)) {
		msg := "ERROR NewSystemMessage, illegal status byte 0x%x  '%s'"
		err = fmt.Errorf(msg, byte(status), status)
		return sys, err
	}
	sys = &SystemMessage{bytes}
	return sys, err
}

func NewClockMessage() (*SystemMessage, error) {
	bytes := []byte{byte(ClockStatus)}
	return newSystemMessage(bytes)
}

func NewStartMessage() (*SystemMessage, error) {
	bytes := []byte{byte(StartStatus)}
	return newSystemMessage(bytes)
}

func NewContinueMessage() (*SystemMessage, error) {
	bytes := []byte{byte(ContinueStatus)}
	return newSystemMessage(bytes)
}

func NewStopMessage() (*SystemMessage, error) {
	bytes := []byte{byte(StopStatus)}
	return newSystemMessage(bytes)
}

func NewEndSysexMessage() (*SystemMessage, error) {
	bytes := []byte{byte(EndSysexStatus)}
	return newSystemMessage(bytes)
}


// Do not include sysex status
func NewSysexMessage(data []byte) (*SystemMessage, error) {
	s := []byte{byte(SysexStatus)}
	bytes := append(s, data...)
	return newSystemMessage(bytes)
}


func (sys *SystemMessage) Status() StatusByte {
	return StatusByte(sys.bytes[0])
}

func (sys *SystemMessage) Bytes() []byte {
	return sys.bytes
}

func (sys *SystemMessage) IsSystemExclusive() bool {
	return sys.Status() == SysexStatus
}


func (sys *SystemMessage) Dump() {

	xDumpLine := func (bytes []byte, start int, width int) string {
		acc := fmt.Sprintf("[0x%04X] ", start)
		bcc := " : "
		for i, count := start, width; i < len(bytes) && count > 0; i, count = i+1, count-1 {
			b := bytes[i]
			acc += fmt.Sprintf("%2X ", b)
			if 32 <= b && b < 127 {
				bcc += fmt.Sprintf("%c", b)
			} else {
				bcc += "."
			}
		}
		pad := 9 + 3 * width
		for len(acc) < pad {
			acc += " "
		}
		return acc + bcc
	}

	xDumpSysex := func (bytes []byte) {
		fmt.Println("SystemExclusive Message")
		width := 8
		for i := 0; i < len(bytes); i += width {
			fmt.Println(xDumpLine(bytes, i, width))
		}
	}

	s := sys.Status()
	if s == SysexStatus {
		xDumpSysex(sys.bytes)
	} else {
		fmt.Printf("0x%02X '%s'\n", byte(s), s)
	}
}
	

		
func (sys *SystemMessage) ToPortmidiEvent() (portmidi.Event, error) {
	var err error
	var pm portmidi.Event
	time := portmidi.Timestamp(0)
	stat := int64(sys.Status())
	if sys.IsSystemExclusive() {
		// ISSURE: Should sysex status be included or excluded from payload? 
		payload := sys.Bytes()
		pm = portmidi.Event{time, stat, 0, 0, payload}
	} else {
		payload := []byte{}
		pm = portmidi.Event{time, stat, 0, 0, payload}
	}
	return pm, err
}
	
