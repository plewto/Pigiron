package op

import (
	"fmt"
	
	"github.com/rakyll/portmidi"
	"github.com/plewto/pigiron/midi"
	"github.com/plewto/pigiron/config"
)


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
		fmt.Println("%s %s\n", op, err)
		return false
	}
	return true
}


//* Read all MIDIInput and process events.
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

func (op *MIDIInput) DeviceID() portmidi.DeviceID {
	return op.devID
}

func (op *MIDIInput) Stream() *portmidi.Stream {
	return op.stream
}

func (op *MIDIInput) DeviceName() string {
	return op.devInfo.Name
}

func (op *MIDIInput) IsOpen() bool {
	return op.devInfo.IsOpened
}

func (op *MIDIInput) Close() {
	op.Stream().Close()
}

func (op *MIDIInput) Info() string {
	s := op.commonInfo()
	s += fmt.Sprintf("\tDevice ID   : %d\n", op.DeviceID())
	s += fmt.Sprintf("\tDevice Name : %s\n", op.DeviceName())
	return s
}
