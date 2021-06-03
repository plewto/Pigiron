package op

import (
	"fmt"
	
	writer "gitlab.com/gomidi/midi/writer"
	gomidi "gitlab.com/gomidi/midi"
	midi "github.com/plewto/pigiron/midi"
)


type DeviceWrapper interface {
	DeviceName() string
}


// MIDIOutput is an Operator wrapper for MIDI output devices.
// MIDIOutput implements the PigOp and DeviceWrapper interfaces.
//
// There may only be a single MIDIOutput for any one MIDI output device.
// Upon construct new MIDIutputs are cached.  Attempts to create another
// MIDIOutput for the same output device returns the cached object.
//
type MIDIOutput struct {
	Operator
	device gomidi.Out
	writer *writer.Writer
}

func (op *MIDIOutput) Close() {
	if op.device != nil {
		op.device.Close()
	}
}

var outputCache map[string]*MIDIOutput = make(map[string]*MIDIOutput)

func midiOutputExists(deviceName string) bool {
	_, flag := outputCache[deviceName]
	return flag
}

func newMIDIOutput(name string, deviceName string) (*MIDIOutput, error) {
	var err error
	var device gomidi.Out
	var op *MIDIOutput
	device, err = midi.GetOutputDevice(deviceName)
	if err == nil {
		op = new(MIDIOutput)
		initOperator(&op.Operator, "MIDIOutput", name, NoChannel)
		op.device = device
		err = op.device.Open()
		if err == nil {
			op.writer = writer.New(op.device)
		}
		assignName(op)
	}
	return op, err
}
		
func MakeMIDIOutput(name string, deviceSpec string) (*MIDIOutput, error) {
	var op *MIDIOutput
	var flag bool
	devname, err := midi.GetOutputDeviceName(deviceSpec)
	if err == nil {
		op, flag = outputCache[devname]
		if !flag {
			op, err = newMIDIOutput(name, devname)
			if err == nil {
				outputCache[devname] = op
			}
		}
	}
	return op, err
}
			
func (op *MIDIOutput) Info() string {
	s := op.commonInfo()
	s += fmt.Sprintf("\tdevice: %v\n", op.device)
	return s
}

func (op *MIDIOutput) DeviceName() string {
	if op.device != nil {
		return fmt.Sprintf("%v", op.device)
	} else {
		return "<nil>"
	}
}


func (op *MIDIOutput) Send(message gomidi.Message) {
	if op.MIDIEnabled() {
		if op.writer != nil {
			op.writer.Write(message)
		}
		op.distribute(message)
	}
}
