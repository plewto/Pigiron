package op

import (
	"fmt"

	"github.com/rakyll/portmidi"
	"github.com/plewto/pigiron/midi"
)

// Monitor is an Operator for real-time monitoring of MIDI events.
//
type Monitor struct {
	baseOperator
}

func newMonitor(name string) *Monitor {
	op := new(Monitor)
	initOperator(&op.baseOperator, "Monitor", name, midi.NoChannel)
	return op
}

func (op *Monitor) Send(event portmidi.Event) {
	op.distribute(event)
	// fmt.Printf(midi.EventToString(event))
	// fmt.Println()
	fmt.Print(formatEvent(event))
}

func iMin(a int, b int) int {
	if a < b {
		return a
	} else {
		return b
	}
}


func formatSysex(event portmidi.Event) string {
	data := event.SysEx
	ln := len(data)
	maxData := 16
	var postfix string
	if ln > maxData {
		postfix = "..."
	}
	acc := ""
	for i:=0; i < iMin(ln, maxData); i++ {
		acc += fmt.Sprintf("%02x ", data[i])
	}
	acc += postfix
	return acc
}


func formatEvent(event portmidi.Event) string {
	st := event.Status
	var acc = midi.StatusMnemonic(st)
	if st >= 0xF0 {
		switch st {
		case 0xF0:
			acc += formatSysex(event)
		default:
		}
		acc += "\n"
	} else {
		st, ci := st & 0xF0, st & 0x0F
		d1, d2 := event.Data1, event.Data2
		acc += fmt.Sprintf(" chan %2d ", ci+1)
		switch st {
		case 0x80:
			acc += fmt.Sprintf("key %3d %4s, vel %3d", d1, midi.KeyName(d1), d2)
		case 0x90:
			acc += fmt.Sprintf("key %3d %4s, vel %3d", d1, midi.KeyName(d1), d2)
		case 0xA0:
			acc += fmt.Sprintf("key %3d %4s, pressure %3d", d1, midi.KeyName(d1), d2)
		case 0xB0:
			acc += fmt.Sprintf("%3d %3d", d1, d2)
		case 0xC0:
			acc += fmt.Sprintf("%3d", d1)
		case 0xD0:
			acc += fmt.Sprintf("%3d", d1)
		case 0xE0:
			acc += fmt.Sprintf("%3d %3d", d1, d2)
		default:
			acc += " ? "
			
		}
		acc += "\n"
		
	}
	return acc
}
			
