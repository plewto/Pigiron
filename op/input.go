package op

import (
	"fmt"
	
	"github.com/rakyll/portmidi"
	goosc "github.com/hypebeast/go-osc/osc"
	"github.com/plewto/pigiron/midi"
	"github.com/plewto/pigiron/config"
)


// MIDIInput is an Operator wrapper for MIDI input devices.
// For each available MIDI device there may be only one corresponding MIDIInput.
// If an attempt is made to create a MIDIInput for a device which is already
// in use, the original MIDIInput is returned.
//
type MIDIInput struct {
	baseOperator
	devID portmidi.DeviceID
	devInfo *portmidi.DeviceInfo
	stream *portmidi.Stream
}

var inputCache = make(map[portmidi.DeviceID]*MIDIInput)

func newMIDIInput(name string,  devID portmidi.DeviceID) (*MIDIInput, error) {
	var op *MIDIInput
	var stream *portmidi.Stream
	var err error
	op = new(MIDIInput)
	initOperator(&op.baseOperator, "MIDIInput", name, midi.NoChannel)
	op.devID = devID
	op.devInfo = portmidi.Info(devID)
	bufferSize := config.GlobalParameters.MIDIInputBufferSize
	stream, err = portmidi.NewInputStream(devID, bufferSize)
	if err != nil {
		return op, err
	}
	op.stream = stream
	op.addCommandHandler("q-device", op.remoteQueryDevice)
	return op, err
}

func NewMIDIInput(name string, deviceSpec string) (*MIDIInput, error) {
	var op *MIDIInput
	var err error
	var devID portmidi.DeviceID
	var cached bool
	devID, err = midi.GetInputID(deviceSpec)
	if err != nil {
		return op, err
	}
	if op, cached = inputCache[devID]; !cached {
		op, err = newMIDIInput(name, devID)
		if err == nil {
			inputCache[devID] = op
			register(op)
		}
	}
	return op, err
}


func notInputError(op *MIDIInput, err error) bool {
	if err != nil {
		fmt.Printf("%s %s\n", op, err)
		return false
	}
	return true
}


// ProcessInputs() polls all MIDIinputs to process weighting events.
//
func ProcessInputs() {
	for _, op :=range inputCache {
		flag, err := op.stream.Poll()
		if notInputError(op, err) {
			if flag {
				bufsize := int(config.GlobalParameters.MIDIInputBufferSize)
				events, err := op.stream.Read(bufsize)
				if notInputError(op, err) {
					for _, event := range events {
						op.Send(event)
					}
				}
			}
		}
	}
}

func (op *MIDIInput) String() string {
	msg := "%-12s name: \"%s\"  device: \"%s\""
	return fmt.Sprintf(msg, op.opType, op.name, op.DeviceName())
}


// op.DeviceID() returns the portmidi.DeviceID for the wrapped device.
//
func (op *MIDIInput) DeviceID() portmidi.DeviceID {
	return op.devID
}

// op.Stream() returns portmidi.Stream
//
func (op *MIDIInput) Stream() *portmidi.Stream {
	return op.stream
}

// op.DeviceName() returns name for wrapped portmidi device.
//
func (op *MIDIInput) DeviceName() string {
	return op.devInfo.Name
}

// op.IsOpen() returns true if wrapped portmidi device is open.
//
func (op *MIDIInput) IsOpen() bool {
	return op.devInfo.IsOpened
}

// op.Close() closes wrapped portmidi device.
// Close should only be called when the application terminates.
//
func (op *MIDIInput) Close() {
	op.Stream().Close()
}

func (op *MIDIInput) Info() string {
	s := op.commonInfo()
	s += fmt.Sprintf("\tDevice ID   : %d\n", op.DeviceID())
	s += fmt.Sprintf("\tDevice Name : %s\n", op.DeviceName())
	return s
}

// op.remoteQueryDevice() extended osc handler for q-device
// osc /pig/op <name>, q-device
// osc returns wrapped portmidi device name.
//
func (op *MIDIInput) remoteQueryDevice(_ *goosc.Message)([]string, error) {
	var err error
	id := fmt.Sprintf("%v", op.DeviceID())
	name := fmt.Sprintf("\"%s\"", op.DeviceName())
	return []string{id, name}, err
}
	
