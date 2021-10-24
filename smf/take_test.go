package smf

import (
	"testing"
	"fmt"
)

func TestTakeByte(t *testing.T) {
	fmt.Print("")
	buffer := []byte{1, 2, 3}
	value, newBuffer, err := TakeByte(buffer)
	if value != byte(1) {
		errmsg := "Expected TakeByte return of 1, got %d"
		t.Fatalf(errmsg, value)
	}
	if len(newBuffer) != 2 || newBuffer[0] != 2 || newBuffer[1] != 3 {
		errmsg := "TakeByte expected to return new byte buffer [2, 3], got %v"
		t.Fatalf(errmsg, newBuffer)
	}
	if err != nil {
		errmsg := "TakeByte returned incorrect error: %s"
		t.Fatalf(errmsg, err)
	}

	// empty buffer test
	buffer = []byte{}
	_, _, err = TakeByte(buffer)
	if err == nil {
		errmsg := "TakeByte did not detect empty buffer"
		t.Fatalf(errmsg)
	}
}


func TestTakeShort(t *testing.T) {
	buffer := []byte{1, 2, 3}
	value, newBuffer, err := TakeShort(buffer)
	expectedValue := int(1<<8 + 2)
	if value != expectedValue {
		errmsg := "TakeShort expected to return value %d, got %d"
		t.Fatalf(errmsg, expectedValue, value)
	}
	if len(newBuffer) != 1 || newBuffer[0] != 3 {
		errmsg := "TakeShort expected to return new byte buffer [3], got %v"
		t.Fatalf(errmsg, newBuffer)
	}
	if err != nil {
		errmsg := "TakeShort returned unexpected error: %s"
		t.Fatalf(errmsg, err)
	}
	// empty buffer test
	buffer = []byte{1}
	_, _, err = TakeShort(buffer)
	if err == nil {
		errmsg := "TakeShort did not detect empty buffer"
		t.Fatalf(errmsg)
	}
}

func TestTakeLong(t *testing.T) {
	buffer := []byte{1, 2, 3, 4}
	value, newBuffer, err := TakeLong(buffer)
	expectedValue := int(1<<24 + 2<<16 + 3<<8 + 4)
	if value != expectedValue {
		errmsg := "TakeLong expected to return value %d, got %d"
		t.Fatalf(errmsg, expectedValue, value)
	}
	if len(newBuffer) != 0 {
		errmsg := "TakeLong expected to return new byte buffer [], got %v"
		t.Fatalf(errmsg, newBuffer)
	}
	if err != nil {
		errmsg := "TakeLong returned unexpected error: %s"
		t.Fatalf(errmsg, err)
	}
	// empty buffer test
	buffer = []byte{1}
	_, _, err = TakeLong(buffer)
	if err == nil {
		errmsg := "TakeLong did not detect empty buffer"
		t.Fatalf(errmsg)
	}
}


func TestTakeVLQ(t *testing.T) {
	// test buffer contains 3-vlq values
	// 0x40 --> 0x40
	// 0x81, 0x00 --> 0x80
	// 0xFF, 0xFF, 0xFF, 0x7F --> 0x0FFFFFFF
	//
	buffer := []byte{0x40, 0x81, 0x00, 0xFF, 0xFF, 0xFF, 0x7F}

	// pass 1, single byte
	vlq, newBuffer, err := TakeVLQ(buffer)
	expected := 0x40
	if vlq.Value() != expected {
		errmsg := "TakeVLQ (pass 1) did note return expected value 0x%2X, got 0x%2X"
		t.Fatalf(errmsg, expected, vlq.Value())
	}
	if len(newBuffer) != len(buffer)-1 || newBuffer[0] != 0x81 {
		errmsg := "TakeVLQ (pass 1) did not return expected newBuffer\n"
		errmsg += "Expected %d byte buffer, and buffer[0] = 0x81, got %v"
		t.Fatalf(errmsg, len(buffer)-1, buffer)
	}
	if err != nil {
		errmsg := "TakeVLQ (pass 1) returned unexpected error, %s"
		t.Fatalf(errmsg, err)
	}

	// pass 2, 2 byte
	buffer = newBuffer
	vlq, newBuffer, err = TakeVLQ(buffer)
	expected = 0x80
	if vlq.Value() != expected {
		errmsg := "TakeVLQ (pass 2) did note return expected value 0x%2X, got 0x%2X"
		t.Fatalf(errmsg, expected, vlq.Value())
	}
	if len(newBuffer) != 4 || newBuffer[0] != 0xFF {
		errmsg := "TakeVLQ (pass 2) did not return expected newBuffer\n"
		errmsg += "Expected 4 byte buffer, and buffer[0] = 0xFF, got %v"
		t.Fatalf(errmsg, newBuffer)
	}
	if err != nil {
		errmsg := "TakeVLQ (pass 2) returned unexpected error, %s"
		t.Fatalf(errmsg, err)
	}

	// pass 3, 4-byte
	buffer = newBuffer
	vlq, newBuffer, err = TakeVLQ(buffer)
	expected = 0x0FFFFFFF
		if vlq.Value() != expected {
		errmsg := "TakeVLQ (pass 3) did note return expected value 0x%2X, got 0x%2X"
		t.Fatalf(errmsg, expected, vlq.Value())
	}
	if len(newBuffer) != 0 {
		errmsg := "TakeVLQ (pass 3) did not return expected newBuffer\n"
		errmsg += "Expected empty buffer, got %v"
		t.Fatalf(errmsg, newBuffer)
	}
	if err != nil {
		errmsg := "TakeVLQ (pass 3) returned unexpected error, %s"
		t.Fatalf(errmsg, err)
	}

	buffer = []byte{0xFF, 0xFF, 0xFF}
	_, _, err = TakeVLQ(buffer)
	if err == nil {
		errmsg := "TakeVLQ did not return error for malformed buffer"
		t.Fatalf(errmsg)
	}	
}
