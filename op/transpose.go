package op

import (
	"fmt"
	goosc "github.com/hypebeast/go-osc/osc"
	"github.com/rakyll/portmidi"	
	"github.com/plewto/pigiron/midi"
)

// Transposer is an Operator which selectivlty modifies MIDI data bytes.
// Either data byte of any channel message may be modified.
// Transposer implements the Operator and midi.Transform interfaces.
//
type Transposer struct {
	baseXformOperator
	status midi.StatusByte
	dataNumber midi.DataNumber
}

func newTransposer (name string) *Transposer {
	op := new(Transposer)
	op.status = midi.KEYED_STATUS
	op.dataNumber = midi.DATA_1
	initOperator(&op.baseOperator, "Transposer", name, midi.NoChannel)
	initXformOperator(&op.baseXformOperator)
	op.initLocalHandlers()
	op.Reset()
	return op
}

func (op *Transposer) Reset() {
	xbase := &op.baseXformOperator
	xbase.Reset()
	op.status = midi.KEYED_STATUS
	op.dataNumber = midi.DATA_1
}

func (op *Transposer) Info() string {
	s := op.commonInfo()
	s += fmt.Sprintf("\tStatus    : 0x%02X  %s\n", byte(op.status), op.status)
	s += fmt.Sprintf("\tData byte : %s\n", op.dataNumber)
	s += fmt.Sprintf("%s\n", op.Dump())
	return s
}

func (op *Transposer) Send(event portmidi.Event) {
	st := midi.StatusByte(event.Status & 0xF0)
	keyed := (st == midi.NOTE_OFF || st == midi.NOTE_ON)
	if st == op.status || ((op.status == midi.KEYED_STATUS) && keyed) {
		if op.dataNumber == midi.DATA_2 {
			data := byte(event.Data2)
			data, _ = op.Value(data)
			event.Data2 = int64(data)
		} else {
			data := byte(event.Data1)
			data, _ = op.Value(data)
			event.Data1 = int64(data)
		}
	}
	op.distribute(event)
}
		
func (op *Transposer) Length() int {
	return op.Length()
}


func (op *Transposer) initLocalHandlers() {

	channelStats := map[int64]string{0x00 : "DISABLED",
		0x01 : "KEYED",
		0x80 : "NOTE_OFF",
		0x90 : "NOTE_ON",
		0xA0 : "POLY_PRESSURE",
		0xB0 : "CONTROLLER",
		0xC0 : "PROGRAM",
		0xD0 : "MONO_PRESSURE",
		0xE0 : "PITCH_BEND"}
		
		
	
	isChannelStatus := func(n int64) error {
		var err error
		_, flag := channelStats[n]
		if !flag {
			msg := "Expected valid status value to Transposer, got 0x%02X"
			err = fmt.Errorf(msg, n)
		}
		return err
	}

	
	// cmd op name, select-status, status
	// osc pig/op name, set-status, status
	//
	// status 0x00 - DISABLE
	//        0x01 - KEY_OFF & KEY_ON 
	//        0x80 - KEY_OFF
	//        0x90 - KEY_ON
	//        0xA0 - POLY_PRESSURE
	//        0xB0 - CONTROLLER
	//        0xC0 - PROGRAM
	//        0xD0 - MONO_PRESSURE
	//        0xE0 - PITCH_BEND
	//
	remoteSetStatus := func(msg *goosc.Message)([]string, error) {
		args, err := ExpectMsg("osi", msg)
		if err != nil {
			return empty, err
		}
		st := args[2].I
		err = isChannelStatus(st)
		if err != nil {
			return empty, err
		}
		op.status = midi.StatusByte(st)
		return empty, err
	}

	// cmd op name, q-status
	// osc /pig/op name, q-status
	//
	// Returns: Selected MIDI status byte.
	//
	remoteQueryStatus := func(msg *goosc.Message)([]string, error) {
		_, err := ExpectMsg("os", msg)
		if err != nil {
			return empty, err
		}
		st := fmt.Sprintf("0x%02X", byte(op.status))
		return []string{st}, err
	}

	// cmd op name, select-data-byte, n       where n = 1 or 2
	// osc /pig/op name, select-data, n
	//
	remoteSelectDataByte := func(msg *goosc.Message)([]string, error) {
		args, err := ExpectMsg("osi", msg)
		if err != nil {
			return empty, err
		}
		n := args[2].I
		switch n {
		case 1:
			op.dataNumber = midi.DATA_1
		case 2:
			op.dataNumber = midi.DATA_2
		default:
			msg := "Expected data byte 1 or 2, got %d"
			err = fmt.Errorf(msg, n)
		}
		return empty, err
	}

	// cmd op name, q-data-byte
	remoteQuerySelectedDataByte := func(msg *goosc.Message)([]string, error) {
		var err error
		var s string
		switch op.dataNumber {
		case midi.DATA_1: s = "1"
		case midi.DATA_2: s = "2"
		default: s = "?"
		}
		return []string{s}, err
	}
		
	op.addCommandHandler("select-status", remoteSetStatus)
	op.addCommandHandler("q-status", remoteQueryStatus)
	op.addCommandHandler("select-data-byte", remoteSelectDataByte)
	op.addCommandHandler("q-data-byte", remoteQuerySelectedDataByte)
}
		
		
