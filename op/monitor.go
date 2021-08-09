package op

import (
	"fmt"
	goosc "github.com/hypebeast/go-osc/osc"
	"github.com/rakyll/portmidi"
	"github.com/plewto/pigiron/midi"
)

// Monitor is an Operator for real-time monitoring of MIDI events.
//
type Monitor struct {
	baseOperator
	excludeStatusFlags map[midi.StatusByte]bool
	enable bool
}

func newMonitor(name string) *Monitor {
	op := new(Monitor)
	op.excludeStatusFlags = map[midi.StatusByte]bool{}
	initOperator(&op.baseOperator, "Monitor", name, midi.MultiChannel)
	op.initLocalHandlers()
	op.SelectAllChannels()
	op.enable = true
	return op
}

func (op *Monitor) Reset() {
	(&op.baseOperator).Reset()
	op.SelectAllChannels()
	for key, _ := range op.excludeStatusFlags {
		op.excludeStatusFlags[key] = false
	}
	op.enable = true
}

func (op *Monitor) Info() string {
	acc := op.commonInfo()
	acc += fmt.Sprintf("Enabled : %v\n", op.enable)
	acc += fmt.Sprintf("Excluded Status : ")
	for key, flag := range op.excludeStatusFlags {
		if flag {
			acc += fmt.Sprintf("0x%02x ", byte(key))
		}
	}
	acc += "\n"
	return acc
}


func (op *Monitor) monitorEvent(event portmidi.Event) bool {
	if !op.enable {
		return false
	}
	st := midi.StatusByte(event.Status)
	cmd := st & midi.StatusByte(0xF0)
	ci := midi.MIDIChannelNibble(st & 0x0F)
	flag, exists := op.excludeStatusFlags[cmd]
	if flag && exists {
		return false
	}
	if midi.IsChannelStatus(byte(cmd)) {
		return op.ChannelIndexSelected(ci)
	}
	return true
}

func (op *Monitor) Send(event portmidi.Event) {
	op.distribute(event)
	if op.monitorEvent(event) {
		fmt.Print(formatEvent(event))
	}
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
	st := midi.StatusByte(event.Status)
	//var acc = "MON 0x%02X " + (st & 0xF0).String()
	var acc = fmt.Sprintf("MON 0x%02X %s ", byte(st), st)
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
			
func (op *Monitor) initLocalHandlers() {

	// op name, q-excluded-status
	// --> list of blocked status bytes
	//
	remoteQueryStatus := func(msg *goosc.Message)([]string, error) {
		var err error
		acc := make([]string, 0, len(op.excludeStatusFlags))
		for k, flag := range op.excludeStatusFlags {
			if flag {
				acc = append(acc, fmt.Sprintf("0x%02X", byte(k)))
			}
		}
		return acc, err
	}

	// op name, exclude-status, st, bool
	// if bool true exclude status st from output.
	// if false, remove from excluded list.
	//
	remoteBlockStatus := func(msg *goosc.Message)([]string, error) {
		args, err := ExpectMsg("osib", msg)
		if err != nil {
			return empty, err
		}
		s := midi.StatusByte(args[2].I)
		flag := args[3].B
		op.excludeStatusFlags[s] = flag
		return empty, err
	}

	// op name, enable, bool
	// if true enable event printing.
	//
	remoteEnable := func(msg *goosc.Message)([]string, error) {
		args, err := ExpectMsg("osb", msg)
		if err != nil {
			return empty, err
		}
		op.enable = args[2].B
		return empty, err
	}

	// op name, q-enabled
	// returns enable flag
	//
	remoteQueryEnable := func(msg *goosc.Message)([]string, error) {
		var err error
		acc := []string{fmt.Sprintf("%v", op.enable)}
		return acc, err
	}
	
	op.addCommandHandler("q-excluded-status", remoteQueryStatus)
	op.addCommandHandler("exclude-status", remoteBlockStatus)
	op.addCommandHandler("enable", remoteEnable)
	op.addCommandHandler("q-enabled", remoteQueryEnable)
}
