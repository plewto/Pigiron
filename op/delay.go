package op

// NOTE: Delay Operator is experimental and at times unstable


import (
	"fmt"
	"time"
	goosc "github.com/hypebeast/go-osc/osc"
	"github.com/rakyll/portmidi"
	"github.com/plewto/pigiron/midi"
	"github.com/plewto/pigiron/pattern"
	"github.com/plewto/pigiron/pigerr"
	
)


const (
	DELAY_PROGRAM_COUNT = 6
	MAX_DELAY_COUNT = 6
	MAX_DELAY_TIME = 1000
)


type delayProgram struct {
	table midi.DataTable
	timePattern pattern.Pattern
	delayCount int
	velocityShift int
}

func newDelayProgram() *delayProgram {
	dp := new(delayProgram)
	dp.table = *midi.NewDataTable()
	dp.Reset()
	return dp
}

func (dp *delayProgram) Reset() {
	dp.table.Reset()
	dp.timePattern = pattern.NewCycle([]int{100})
	dp.delayCount = 4
	dp.velocityShift = -8
}


// Delay is an experimental Operator which produces an echo like effect using.
// Delayed notes may be key-mapped with velocity scaling and variable delay
// times.  Up to 6 presets may be accessed via MIDI program-change.
// 
// Implements Operator, midi.Transform, midi.Program
//
type Delay struct {
	baseXformOperator
	channelIndex byte
	forwardProgramChange bool
	programs [DELAY_PROGRAM_COUNT]*delayProgram
	currentProgramSlot byte
}

func newDelay(name string) *Delay {
	op := new(Delay)
	initOperator(&op.baseOperator, "Delay", name, midi.SingleChannel)
	initXformOperator(&op.baseXformOperator)
	op.programs = [DELAY_PROGRAM_COUNT]*delayProgram{}
	for i := 0; i < DELAY_PROGRAM_COUNT; i++ {
		op.programs[i] = newDelayProgram()
	}
	op.initLocalHandlers()
	op.Reset()
	pigerr.Warning("Using experimental Delay Operator.",
		"Program may become unstable.")
	
	return op
}
	
func (op *Delay) Reset() {
	xbase := &op.baseXformOperator
	xbase.Reset()
	for _, dp := range(op.programs) {
		dp.Reset()
	}
	op.channelIndex = 0
	op.forwardProgramChange = true
	op.currentProgramSlot = 0
}

func (op *Delay) Info() string {
	s := op.commonInfo()
	s += fmt.Sprintf("\tMIDI Channel    : %2d\n", op.channelIndex + 1)
	s += fmt.Sprintf("\tCurrent Program : %2d\n", op.currentProgramSlot)
	s += fmt.Sprintf("\tForward program changes : %v\n", op.forwardProgramChange)
	s += fmt.Sprintf("%s\n", op.Dump())
	return s
}

// midi.ChannelSelector interface

func (op *Delay) ChannelMode() midi.ChannelMode {
	return midi.SingleChannel
}

func (op *Delay) EnableChannel(c midi.MIDIChannel, _ bool) error {
	var err error
	if c < 0 || 16 < c {
		msg := "Illegal MIDI channel %d"
		err = fmt.Errorf(msg, byte(c))
		return err
	}
	op.channelIndex = byte(c) - 1
	return err
}

func (op *Delay) SelectChannel(c midi.MIDIChannel) error  {
	return op.EnableChannel(c, true)
}
	
func (op *Delay) SelectedChannelIndexes() []midi.MIDIChannelNibble {
	acc := make([]midi.MIDIChannelNibble, 1)
	acc[0] = midi.MIDIChannelNibble(op.channelIndex)
	return acc
}

func (op *Delay) ChannelIndexSelected(ci midi.MIDIChannelNibble) bool  {
	return byte(ci) == op.channelIndex
}

func (op *Delay) DeselectAllChannels() {}
func (op *Delay) SelectAllChanels() {}
	

// midi.Program interface

func (op *Delay) ProgramRange() (floor byte, ceiling byte) {
	return 0, byte(DELAY_PROGRAM_COUNT)
}	

func (op *Delay) CurrentProgram() byte	 {
	return op.currentProgramSlot
}

func (op *Delay) ChangeProgram(event portmidi.Event) {
	prog := byte(event.Data1)
	if 0 <= prog || prog < DELAY_PROGRAM_COUNT {
		op.currentProgramSlot = prog
	}
	if op.forwardProgramChange {
		op.distribute(event)
	}
}

func (op *Delay) program() *delayProgram {
	return op.programs[op.currentProgramSlot]
}


func (op *Delay) loop(event portmidi.Event) {
	program := op.program()
	keytab := program.table
	tpat := program.timePattern
	vshift := int64(program.velocityShift)
	count := program.delayCount
	if count > MAX_DELAY_COUNT {
		count = MAX_DELAY_COUNT
	}
	// isNoteOff := event.Status & 0xF0 == 0x80
	for i := 0; i < count; i++ {
		delay := time.Duration(tpat.Next())
		time.Sleep(delay * time.Millisecond)
		key, _ := keytab.Value(byte(event.Data1))
		vel := (event.Data2 + vshift)
		switch {
		case vel > 0x7F: vel = 0x7F
		case vel < 0: vel = 0x00
		// case isNoteOff: vel = 0x00
		default:
			if vel == 0 {  // note_on with velocity 0
				break
			}
		}
		event.Data1 = int64(key)
		event.Data2 = int64(vel)
		op.distribute(event)
	}
}

func (op *Delay) Send(event portmidi.Event) {
	st := byte(event.Status)
	cmd, ci := st & 0xF0, st & 0x0F
	if ci != op.channelIndex {
		op.distribute(event)
	} else {
		switch cmd {
		case 0x80:
			go op.loop(event)
			op.distribute(event)
		case 0x90:
			go op.loop(event)
			op.distribute(event)
		case 0xC0:
			op.ChangeProgram(event)
		default:
			op.distribute(event)
		}
	}
}


func (op *Delay) initLocalHandlers() {

	remoteQueryDelayPattern := func(msg *goosc.Message)([]string, error) {
		var err error
		values := op.program().timePattern.Values()
		acc := make([]string, len(values))
		for i := 0; i < len(values); i++ {
			acc[i] = fmt.Sprintf("%d", values[i])
		}
		return acc, err
	}

	remoteQueryRepeat := func(msg *goosc.Message)([]string, error) {
		var err error
		count := op.program().delayCount
		acc := []string{fmt.Sprintf("%d", count)}
		return acc, err
	}

	remoteQueryVelocityShift := func(msg *goosc.Message)([]string, error) {
		var err error
		vshift := op.program().velocityShift
		acc := []string{fmt.Sprintf("%d", vshift)}
		return acc, err
	}

	// op name set-delay-pattern i, [i, i, ...] 
	remoteSetDelayPattern := func(msg *goosc.Message)([]string, error) {
		extra := len(msg.Arguments) - 3
		template := "osi"
		for i := 0; i < extra; i++ {
			template += "i"
		}
		args, err := ExpectMsg(template, msg)
		if err != nil {
			return empty, err
		}
		values := make([]int, extra+1)
		for i := 0; i < len(values); i++ {
			v := int(args[i+2].I)
			switch {
			case v < 0: v = 0
			case v > MAX_DELAY_TIME: v = MAX_DELAY_TIME
			}
			values[i] = v
		}
		op.program().timePattern = pattern.NewCycle(values)
		return empty, err
	}

	remoteSetRepeat := func(msg *goosc.Message)([]string, error) {
		args, err := ExpectMsg("osi", msg)
		if err != nil {
			return empty, err
		}
		count := int(args[2].I)
		switch {
		case count < 0: count = 1
		case count > MAX_DELAY_COUNT: count = MAX_DELAY_COUNT
		}
		op.program().delayCount = count
		return empty, err
	}

	remoteSetVelocityShift := func(msg *goosc.Message)([]string, error) {
		args, err := ExpectMsg("osi", msg)
		if err != nil {
			return empty, err
		}
		vshift := int(args[2].I)
		switch {
		case vshift < -64: vshift = -64
		case vshift > 64: vshift = 64
		}
		op.program().velocityShift = vshift
		return empty, err
	}

	remoteUseProgram := func(msg *goosc.Message)([]string, error) {
		args, err := ExpectMsg("osi", msg)
		if err != nil {
			return empty, err
		}
		p := int(args[2].I) & 0x7F
		if p >= 0 && p < len(op.programs) {
			op.currentProgramSlot = byte(p)
		}
		return empty, err
	}

	remoteQueryProgramNumber := func(msg *goosc.Message)([]string, error) {
		var err error
		acc := []string{fmt.Sprintf("%d", op.currentProgramSlot)}
		return acc, err
	}
			
	
	op.addCommandHandler("q-delay-pattern", remoteQueryDelayPattern)
	op.addCommandHandler("q-repeat-count", remoteQueryRepeat)
	op.addCommandHandler("q-velocity-shift", remoteQueryVelocityShift)
	op.addCommandHandler("q-program-number", remoteQueryProgramNumber)
	op.addCommandHandler("set-delay-pattern", remoteSetDelayPattern)
	op.addCommandHandler("set-repeat-count", remoteSetRepeat)
	op.addCommandHandler("set-velocity-shift", remoteSetVelocityShift)
	op.addCommandHandler("use-program", remoteUseProgram)
	
}
	
