package op

import (
	"C"
	"fmt"
	"time"
	
	"github.com/rakyll/portmidi"
	"github.com/plewto/pigiron/midi"
	"github.com/plewto/pigiron/config"
)

// MIDIOutput is an Operator wrapper for MIDI output devices.
// MIDIOutput implements the PigOp and DeviceWrapper interfaces.
//
// There may only be a single MIDIOutput for any one MIDI output device.
// Upon construct new MIDIutputs are cached.  Attempts to create another
// MIDIOutput for the same output device returns the cached object.
//
type MIDIOutput struct {
	baseOperator
	devID portmidi.DeviceID
	devInfo *portmidi.DeviceInfo
	stream *portmidi.Stream
}

var outputCache = make(map[portmidi.DeviceID]*MIDIOutput)


// ** Creates new MIDIOutput, does not cache it.
// ** Only called when cached version does not exists.
// **
func newMIDIOutput(name string,  devID portmidi.DeviceID) (*MIDIOutput, error) {
	var op *MIDIOutput
	var stream *portmidi.Stream
	var err error
	op = new(MIDIOutput)
	initOperator(&op.baseOperator, "MIDIOutput", name, midi.NoChannel)
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

// ** Factory creates new MIDIOutput or grabs it from the cache.
//
func NewMIDIOutput(name string, deviceSpec string) (*MIDIOutput, error) {
	var op *MIDIOutput
	var err error
	var devID portmidi.DeviceID
	var cached bool
	devID, err = midi.GetOutputID(deviceSpec)
	if err != nil {
		return op, err
	}
	if op, cached = outputCache[devID]; !cached {
		op, err = newMIDIOutput(name, devID)
		if err == nil {
			outputCache[devID] = op
			register(op)
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

func (op *MIDIOutput) Close() {
	op.Stream().Close()
}

func (op *MIDIOutput) Info() string {
	s := op.commonInfo()
	s += fmt.Sprintf("\tDevice ID   : %d\n", op.DeviceID())
	s += fmt.Sprintf("\tDevice Name : %s\n", op.DeviceName())
	return s
}


func (op *MIDIOutput) Send(event portmidi.Event) {
	if len(event.SysEx) > 0 {
		op.stream.WriteSysExBytes(portmidi.Time(), event.SysEx)
	} else {
		op.stream.WriteShort(event.Status, event.Data1, event.Data2)
	}
	op.distribute(event)
}

func (op *MIDIOutput) Panic() {
	var event portmidi.Event
	for ci:=0; ci<16; ci++ {
		time.Sleep(1 * time.Millisecond)
		st := int64(0x80 | ci)
		velocity := int64(0)
		for key := int64(0); key<128; key++ {
			if key % 16 == 0 {
				time.Sleep(1 * time.Millisecond)
			}
			event = portmidi.Event{
				Timestamp: 0,
				Status: st,
				Data1: key,
				Data2: velocity}
			op.Send(event)
		}
	}
	for _, child := range op.children() {
		child.Panic()
	}
}



