package op

import (
	"C"
	"fmt"
	gomidi "gitlab.com/gomidi/midi/v2"
	goosc "github.com/hypebeast/go-osc/osc"
	"github.com/plewto/pigiron/backend"
	"github.com/plewto/pigiron/midi"
)

// MIDIOutput is an Operator wrapper for MIDI output devices.
// For each available MIDI device there may be only one corresponding MIDIOutput.
// If an attempt is made to create a MIDIOutput for a device which is already
// in use, the original MIDIOutput is returned.
//
type MIDIOutput struct {
	baseOperator
	port gomidi.Out
	
}

var outputCache = make(map[string]*MIDIOutput)


// ** Creates new MIDIOutput, does not cache it.
// ** Only called when cached version does not exists.
// **
func newMIDIOutput(name string,  port gomidi.Out) *MIDIOutput {
	var op *MIDIOutput
	op = new(MIDIOutput)
	initOperator(&op.baseOperator, "MIDIOutput", name, midi.NoChannel)
	op.port = port
	op.port.Open()
	op.addCommandHandler("q-device", op.remoteQueryDevice)
	outputCache[port.String()] = op
	register(op)
	return op
}

// ** Factory creates new MIDIOutput or grabs it from the cache.
//

func NewMIDIOutput(name string, deviceSpec string) (*MIDIOutput, error) {
	var err error
	var op *MIDIOutput
	var port gomidi.Out
	var cached bool
	port, err = backend.GetMIDIOutput(deviceSpec)
	if err != nil {
		return op, err
	}
	portName := port.String()
	if op, cached = outputCache[portName]; !cached {
		op = newMIDIOutput(name, port)
	}
	return op, err
}
		
	

func (op *MIDIOutput) String() string {
	msg := "%-12s name: \"%s\"  device: \"%s\""
	return fmt.Sprintf(msg, op.opType, op.name, op.port.String())
}

func (op *MIDIOutput) DeviceName() string {
	return op.port.String()
}

func (op *MIDIOutput) Info() string {
	s := op.commonInfo()
	s += fmt.Sprintf("\tDevice Name : %s\n", op.port)
	return s
}

func (op *MIDIOutput) Send(msg gomidi.Message) {
	if op.MIDIOutputEnabled() {
		op.port.Send(msg.Data)
		op.distribute(msg)
	}
}


func (op *MIDIOutput) Panic() {
	fmt.Println("ISSUE: WARNING: MIDIOutput.Paninc not implemented")
}


// op.remoteQueryDevice() extended osc handler for q-device
// osc /pig/op <name>, q-device
// osc returns wrapped port name.
//
func (op *MIDIOutput) remoteQueryDevice(_ *goosc.Message)([]string, error) {
	var err error
	id := fmt.Sprintf("%v", op.port)
	name := fmt.Sprintf("\"%s\"", op.port)
	return []string{id, name}, err
}

