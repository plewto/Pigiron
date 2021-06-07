package op

import (
	"fmt"
	
	portmidi "github.com/rakyll/portmidi"
	midi "github.com/plewto/pigiron/midi"
	config "github.com/plewto/pigiron/config"
)

// MIDIOutput is an Operator wrapper for MIDI output devices.
// MIDIOutput implements the PigOp and DeviceWrapper interfaces.
//
// There may only be a single MIDIOutput for any one MIDI output device.
// Upon construct new MIDIutputs are cached.  Attempts to create another
// MIDIOutput for the same output device returns the cached object.
//
type MIDIOutput struct {
	Operator
	devID portmidi.DeviceID
	devInfo *portmidi.DeviceInfo
	stream *portmidi.Stream
}

var outputCache map[portmidi.DeviceID]*MIDIOutput = make(map[portmidi.DeviceID]*MIDIOutput)


func midiOutputExists(devID portmidi.DeviceID) bool {
	_, flag := outputCache[devID]
	return flag
}
	

func newMIDIOutput(name string,  devID portmidi.DeviceID) (*MIDIOutput, error) {
	var op *MIDIOutput
	var stream *portmidi.Stream
	var err error
	op = new(MIDIOutput)
	initOperator(&op.Operator, "MIDIOutput", name, midi.NoChannel)
	op.devID = devID
	op.devInfo = portmidi.Info(devID)
	bufferSize := config.MIDIOutputBufferSize
	latency := config.MIDIOutputLatency
	stream, err = portmidi.NewOutputStream(devID, bufferSize, latency)
	if err != nil {
		return op, err
	}
	op.stream = stream
	return op, err
}
	
func registerNewMIDIOutput(name string, deviceSpec string) (*MIDIOutput, error) {
	var op *MIDIOutput
	var err error
	var devID portmidi.DeviceID
	devID, err = midi.GetOutputID(deviceSpec)
	if err != nil {
		return op, err
	}
	if midiOutputExists(devID) {
		op = outputCache[devID]
	} else {
		op, err = newMIDIOutput(name, devID)
		if err == nil {
			outputCache[devID] = op
		}
	}
	return op, err
}
			
func (op *MIDIOutput) DeviceID() portmidi.DeviceID {
	return op.devID
}

func (op *MIDIOutput) Stream() *portmidi.Stream {
	return op.stream
}

func (op *MIDIOutput) DeviceName() string {
	return op.devInfo.Name
}

func (op *MIDIOutput) IsOpen() bool {
	return op.devInfo.IsOpened
}

func (op *MIDIOutput) Close() error {
	return op.Stream().Close()
}

func (op *MIDIOutput) Info() string {
	s := op.commonInfo()
	s += fmt.Sprintf("\tDevice ID   : %d\n", op.DeviceID())
	s += fmt.Sprintf("\tDevice Name : %s\n", op.DeviceName())
	return s
}

// TODO Implement MIDIOutput send





