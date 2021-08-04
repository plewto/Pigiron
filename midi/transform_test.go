package midi

import (
	"fmt"
	"testing"
	"github.com/rakyll/portmidi"
)

func TestTransform(t *testing.T) {
	fname := func(testNumber int) string {
		return fmt.Sprintf("\nDataTable test %d: ", testNumber)
	}
	
	dt := NewDataTable()
	for i := 0; i < dt.Length(); i++ {
		j, err := dt.Value(byte(i))
		if err != nil {
			msg := "%s unexpected error, index was %d, \nerr was %s"
			t.Fatalf(msg, fname(1), i, err)
		}
		if byte(i) != byte(j) {
			msg := "%s expected %d, got %d"
			t.Fatalf(msg, fname(1), i, j)
		}
	}

	_, err := dt.Value(byte(200))
	if err == nil {
		msg := "%s did not detect out of bounds index: 200"
		t.Fatalf(msg, fname(2))
	}

	shift := 100
	for i := 10; i < 20; i++ {
		err = dt.SetValue(byte(i), byte(i+shift))
		if err != nil {
			msg := "%s SetValue(%d) retunred unexpected error\nerr was %s"
			t.Fatalf(msg, fname(3), i+shift, err)
		}
	}

	for i := 10; i < 20; i++ {
		j, _ := dt.Value(byte(i))
		if j != byte(i+shift) {
			msg := "%s Value(%d), expected %d, got %d"
			t.Fatalf(msg, fname(4), i, i+shift, j)
		}
	}
	
}
			
func TestTransformBank(t *testing.T) {

	fname := func(testNumber int) string {
		return fmt.Sprintf("\nTransformBank test %d: ", testNumber)
	}

	program := func(channel int64, program int64) portmidi.Event {
		st := int64(PROGRAM) + (channel - 1)
		return portmidi.Event{0, st, program, 0, []byte{}}
	}

	
	shift := 10
	bank := NewTransformBank(2)
	bank.SelectChannel(1)
	event := program(1, 1)
	bank.ChangeProgram(event)
	for i := 0; i < 0x80; i++ {
		j := (i + shift) & 0x7F
		bank.SetValue(byte(i), byte(j))
	}

	event = program(1, 0)
	bank.ChangeProgram(event)
	for i := 0; i < bank.Length(); i++ {
		j, err := bank.Value(byte(i))
		if err != nil {
			msg := "%s, unexpected error index %d,\nerr %s"
			t.Fatalf(msg, fname(1), i, err)
		}
		if j != byte(i) {
			msg := "%s, Value(%d) expected %d, got %d"
			t.Fatalf(msg, fname(2), i, i, j)
		}
	}

	// Out of bounds program should be ignored.
	event = program(1, 100)
	bank.ChangeProgram(event)
	for i := 0; i < bank.Length(); i++ {
		j, err := bank.Value(byte(i))
		if err != nil {
			msg := "%s, unexpected error index %d,\nerr %s"
			t.Fatalf(msg, fname(3), i, err)
		}
		if j != byte(i) {
			msg := "%s, Value(%d) expected %d, got %d"
			t.Fatalf(msg, fname(4), i, i, j)
		}
	}

	// Switch to offset Transform
	event = program(1, 1)
	bank.ChangeProgram(event)
	for i := 0; i < bank.Length(); i++ {
		j, err := bank.Value(byte(i))
		if err != nil {
			msg := "%s, unexpected error index %d,\nerr %s"
			t.Fatalf(msg, fname(5), i, err)
		}
		x := (i + shift) & 0x7F
		if j != byte(x) {
			msg := "%s Value(%d) expected %d, got %d"
			t.Fatalf(msg, fname(6), i, x, j)
		}
	}

	// Check channel discrimination
	event = program(1, 0)
	bank.SelectChannel(1)
	bank.ChangeProgram(event)
	event = program(10, 1)
	bank.ChangeProgram(event)
	if bank.CurrentProgram() != 0 {
		msg := "%s, changed program on wrong MIDI channel"
		t.Fatalf(msg, fname(7))
	}
	bank.SelectChannel(10)
	bank.ChangeProgram(event)
	if bank.CurrentProgram() != 1 {
		msg := "%s did not change program on current MIDI channel."
		t.Fatalf(msg, fname(8))
	}

		
}
