package op

import (
	"fmt"
	gomidi "gitlab.com/gomidi/midi/v2"
	goosc "github.com/hypebeast/go-osc/osc"
	"github.com/plewto/pigiron/backend"
	"github.com/plewto/pigiron/midi"
)


// MIDIInput is an Operator wrapper for MIDI input devices.
// For each available MIDI device there may be only one corresponding MIDIInput.
// If an attempt is made to create a MIDIInput for a device which is already
// in use, the original MIDIInput is returned.
//
type MIDIInput struct {
	baseOperator
	port gomidi.In
}

var inputCache = make(map[string]*MIDIInput)

func newMIDIInput(name string, port gomidi.In) *MIDIInput {
	op := new(MIDIInput)
	initOperator(&op.baseOperator, "MIDIInput", name, midi.NoChannel)
	op.addCommandHandler("q-device", op.remoteQueryDevice)
	return op
}

func NewMIDIInput(name string, deviceSpec string) (*MIDIInput, error) {
	var op *MIDIInput
	var err error
	var port gomidi.In
	var cached bool
	port, err = backend.GetMIDIInput(deviceSpec)
	if err != nil {
		return op, err
	}
	portName := port.String()
	if op, cached = inputCache[portName]; !cached {
		op = newMIDIInput(name, port)
		inputCache[portName] = op
		
		listener, err2 := gomidi.NewListener(op.port, func(msg gomidi.Message, dtime int64) {
			op.Send(msg)
		})
		if err2 != nil {
			msg1 := fmt.Sprintf("MIDIInput %s, can not listen to port %s\n", op.Name(), op.port)
			msg2 := fmt.Sprintf("%s", err2)
			return nil, fmt.Errorf(msg1 + msg2)
		}
		listener.StartListening()
		register(op)
	}
	return op, err
}


func (op *MIDIInput) String() string {
	msg := "%-12s name: \"%s\"  device: \"%s\""
	return fmt.Sprintf(msg, op.opType, op.name, op.port.String())
}

// op.DeviceName() returns name for wrapped portmidi device.
//
func (op *MIDIInput) DeviceName() string {
	return op.port.String()
}

func (op *MIDIInput) Info() string {
	s := op.commonInfo()
	s += fmt.Sprintf("\tDevice Name : %s\n", op.port)
	return s
}

// op.remoteQueryDevice() extended osc handler for q-device
// osc /pig/op <name>, q-device
// osc returns wrapped MIDI port name.
//
func (op *MIDIInput) remoteQueryDevice(_ *goosc.Message)([]string, error) {
	var err error
	name := fmt.Sprintf("\"%s\"", op.port)
	return []string{name}, err
}
	
