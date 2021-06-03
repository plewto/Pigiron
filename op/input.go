package op

import (
	"fmt"
	
	reader "gitlab.com/gomidi/midi/reader"
	gomidi "gitlab.com/gomidi/midi"
	midi "github.com/plewto/pigiron/midi"
)


// implements PigOp, DeviceWrapper
type MIDIInput struct {
	Operator
	device gomidi.In
	reader *reader.Reader
}

func (op *MIDIInput) Close() {
	if op.device != nil {
		op.device.Close()
	}
}

var inputCache map[string]*MIDIInput = make(map[string]*MIDIInput)

func midiInputExists(deviceName string) bool {
	_, flag := inputCache[deviceName]
	return flag
}

func newMIDIInput(name string, deviceName string) (*MIDIInput, error) {
	var err error
	var device gomidi.In
	var op *MIDIInput
	device, err = midi.GetInputDevice(deviceName)
	if err == nil {
		op = new(MIDIInput)
		initOperator(&op.Operator, "MIDIInput", name, NoChannel)
		op.device = device
		err = op.device.Open()
		if err == nil {
			fmt.Printf("MIDI input %v opened\n", op.device)
			op.reader = reader.New(
				reader.NoLogger(),
				reader.Each(func(pos *reader.Position, msg gomidi.Message) {
					op.Send(msg)
				}))
			err = op.reader.ListenTo(op.device)  // TODO: Handle error
		}
		assignName(op)
	}
	return op, err
}
		
func MakeMIDIInput(name string, deviceSpec string) (*MIDIInput, error) {
	var op *MIDIInput
	var flag bool
	devname, err := midi.GetInputDeviceName(deviceSpec)
	if err == nil {
		op, flag = inputCache[devname]
		if !flag {
			op, err = newMIDIInput(name, devname)
			if err == nil {
				inputCache[devname] = op
			}
		}
	}
	return op, err
}
			
func (op *MIDIInput) Info() string {
	s := op.commonInfo()
	s += fmt.Sprintf("\tdevice: %v\n", op.device)
	s += fmt.Sprintf("\treader: %T  is nil: %v\n", op.reader, op.reader == nil)
	return s
}


func (op *MIDIInput) DeviceName() string {
	if op.device != nil {
		return fmt.Sprintf("%v", op.device)
	} else {
		return "<nil>"
	}
}


