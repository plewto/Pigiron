package op

import (
	"fmt"
	"os"
	"time"
	gomidi "gitlab.com/gomidi/midi/v2"
	goosc "github.com/hypebeast/go-osc/osc"
	"github.com/plewto/pigiron/midi"
	"github.com/plewto/pigiron/pigpath"
)

// Monitor is an Operator for real-time monitoring of MIDI events.
// It may optionally save events to a log file.
//
type Monitor struct {
	baseOperator
	excludeStatusFlags map[midi.StatusByte]bool
	enable bool
	logFilename string
	logFile *os.File
	previousTime time.Time
}

func newMonitor(name string) *Monitor {
	op := new(Monitor)
	op.excludeStatusFlags = map[midi.StatusByte]bool{}
	initOperator(&op.baseOperator, "Monitor", name, midi.MultiChannel)
	op.initLocalHandlers()
	op.SelectAllChannels()
	op.enable = true
	op.logFilename = ""
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
	fname := op.logFilename
	if fname == "" {
		fname = "<closed>"
	}
	acc += fmt.Sprintf("\nLog file : '%s'\n", fname)
	acc += "\n"
	return acc
}


func (op *Monitor) Close() {
	op.CloseLogFile()
}

func (op *Monitor) OpenLogFile(filename string) (truename string, err error) {
	op.logFilename = pigpath.SubSpecialDirectories(filename)
	op.logFile, err = os.OpenFile(op.logFilename, os.O_CREATE | os.O_WRONLY, 0644)
	if err != nil {
		msg := "Monitor %s, can not open log file: '%s'\nPrevious error: %v"
		err = fmt.Errorf(msg, op.Name(), op.logFilename, err)
		op.logFilename = ""
		return "", err
	}
	op.logFile.WriteString("# Δt in milliseconds")
	truename = op.logFilename
	op.previousTime = time.Now()
	return truename, err
}
	
func (op *Monitor) CloseLogFile() {
	if op.logFile != nil {
		fmt.Printf("Closing MIDI Monitor log file: %s\n", op.logFilename)
		op.logFile.Close()
		op.logFile = nil
	}
}


func (op *Monitor) monitorMessage(msg gomidi.Message) bool {
	if !op.enable {
		return false
	}
	st := midi.StatusByte(msg.Data[0])
	cmd := st & 0xF0
	ci := midi.MIDIChannelNibble(st & 0x0f)
	flag, exists := op.excludeStatusFlags[st]
	if flag && exists {
		return false
	}
	if midi.IsChannelStatus(cmd) {
		return op.ChannelIndexSelected(ci)
	}
	return true
}


func (op *Monitor) logEvent(s string) {
	if op.logFile != nil {
		op.logFile.WriteString(s)
	}
}


func (op *Monitor) print(msg gomidi.Message){
	if op.monitorMessage(msg) {
		now := time.Now()
		elapse := now.Sub(op.previousTime)
		op.previousTime = now
		s := formatEvent(elapse, msg)
		fmt.Print(s)
		op.logEvent(s)
	}
}


func (op *Monitor) Send(msg gomidi.Message) {
	op.distribute(msg)
	go op.print(msg)
}


func iMin(a int, b int) int {
	if a < b {
		return a
	} else {
		return b
	}
}


func formatSysex(msg gomidi.Message) string {
	data := msg.Data
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

func formatEvent(elapse time.Duration, msg gomidi.Message) string {
	tdelta := elapse.Milliseconds()
	st := midi.StatusByte(msg.Data[0])
	var acc = fmt.Sprintf("Δt %5d 0x%02X %s ", tdelta,  byte(st), st)
	if st >= 0xF0 {
		switch st {
		case 0xF0:
			acc += formatSysex(msg)
		default:
		}
		acc += "\n"
	} else {
		st, ci := st & 0xF0, st & 0x0F
		d1 := msg.Data[1]
		var d2 byte
		if len(msg.Data) > 2 {
			d2 = msg.Data[2]
		}
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

	// op name, open-logfile, filename
	// returns truename of logfile.
	//
	remoteOpenLogfile := func(msg *goosc.Message)([]string, error) {
		args, err := ExpectMsg("oss", msg)
		if err != nil {
			return empty, err
		}
		filename := args[2].S
		var truename string
		truename, err = op.OpenLogFile(filename)
		return []string{truename}, err
	}

	// op name, close-logfile
	//
	remoteCloseLogfile := func(msg *goosc.Message)([]string, error) {
		var err error
		op.CloseLogFile()
		return empty, err
	}

	// op name, q-logfile
	// Returns logfile filename, or '<closed>'
	//
	remoteQueryLogfile := func(msg *goosc.Message)([]string, error) {
		var err error
		var filename string
		if op.logFilename == "" {
			filename = "<closed>"
		} else {
			filename = op.logFilename
		}
		return []string{filename}, err
	}
	
	op.addCommandHandler("q-excluded-status", remoteQueryStatus)
	op.addCommandHandler("exclude-status", remoteBlockStatus)
	op.addCommandHandler("enable", remoteEnable)
	op.addCommandHandler("q-enabled", remoteQueryEnable)
	op.addCommandHandler("open-logfile", remoteOpenLogfile)
	op.addCommandHandler("close-logfile", remoteCloseLogfile)
	op.addCommandHandler("q-logfile", remoteQueryLogfile)
}


