package expect

import (
	"fmt"
	"testing"
)

var testBuffer = []byte{
	0x00,				          // [  0] vlq 0x00
	0x7F,				          // [  1] vlq 0x7f
	0x81, 0x00,			          // [  2] vlq 0x80
	0xff, 0xff, 0x7f,		          // [  4] vlq 0x1FFFFF (2097151 d)
	0xff, 0xff, 0xff, 0x7f,		          // [  7] max vlq 0x0FFFFFFF (268435455 d)
	0xff, 0xff, 0xff, 0xff, 0xff,             // [ 11] malformed vlq
	0x00, 0x7f,			          // [ 16] valid MIDI data bytes
	0x80,				          // [ 18] invalid MIDI data byte
	0x80, 0x00, 0x7f,                         // [ 19] valid note off
	0x80, 0x00, 0x80,		          // [ 22] invalid note off
	0xC0, 0x01,			          // [ 25] program change
	0xf0, 0x00, 0x01, 0x02, 0xF7,             // [ 27] valid sysex message
	0xf0, 0x00, 0x80,                         // [ 32] invalid sysex message
	0xf8,				          // [ 35] clock
	0xff, 0x00, 0x02, 0x00, 0x01,             // [ 36] meta sequence number
	0xff, 0x01, 0x04, 0x41, 0x42, 0x43, 0x44, // [ 41] meta text 'ABCD'
	0xff, 0x60, 0x00,			  // [ 48] malformed meta message
	0xff, 0x2f, 0x00}			  // [ 51] meta end-of-track
	
	
func TestExpectByte(t *testing.T) {
	fmt.Println("TestExpectByte")
	index := 4
	value, newIndex, err := expectByte(testBuffer, index)
	if err != nil {
		t.Fatalf("Unexpected err: %s", err)
	}
	if newIndex != index+1 {
		errmsg := "Expected newIndex = %d, got %d"
		t.Fatalf(errmsg, index+1, newIndex)
	}
	if value != testBuffer[index] {
		errmsg := "Expected byte value 0x%02X, got 0x%02X"
		t.Fatalf(errmsg, testBuffer[index], value)
	}
	_, _, err = expectByte(testBuffer, len(testBuffer))
	if err == nil {
		errmsg := "Did not detect out of bounds index"
		t.Fatalf(errmsg)
	}
}


func TestExpectVLQ(t *testing.T) {
	fmt.Println("TestExpectVLQ")
	var index, newIndex int
	var vlq *VLQ
	var err error
	
	var validate = func(expectValue int, expectNewIndex int, errr error) {
		if errr != nil {
			fmt.Printf("FAIL: TestExpectVLQ index = %d\n", index)
			t.Fatalf("Got unexpected error: %s", errr)
		}
		if newIndex != expectNewIndex {
			fmt.Printf("FAIL: TestExpectVLQ index = %d\n", index)
			errmsg := "Expected newIndex of %d, got %d"
			t.Fatalf(errmsg, expectNewIndex, newIndex)
		}
		if vlq.Value() != expectValue {
			fmt.Printf("FAIL: TestExpectVLQ index = %d\n", index)
			errmsg := "Expected VLQ value of %d, got %d"
			t.Fatalf(errmsg, expectValue, vlq.Value())
		}
	}

	index = 0
	vlq, newIndex, err = expectVLQ(testBuffer, index)
	validate(0, 1, err)

	index = newIndex
	vlq, newIndex, err = expectVLQ(testBuffer, index)
	validate(0x7f, 2, err)

	index = newIndex
	vlq, newIndex, err = expectVLQ(testBuffer, index)
	validate(0x80, 4, err)

	index = newIndex
	vlq, newIndex, err = expectVLQ(testBuffer, index)
	validate(0x1FFFFF, 7, err)

	index = newIndex
	vlq, newIndex, err = expectVLQ(testBuffer, index)
	validate(0x0FFFFFFF, 11, err)

	index = newIndex
	_, _, err = expectVLQ(testBuffer, index)
	if err == nil {
		errmsg := "Did not detect malformed VLQ at index 11"
		t.Fatalf(errmsg)
	}

	index = len(testBuffer)
	_, _, err = expectVLQ(testBuffer, index)
	if err == nil {
		errmsg := "Did not detect out of bounds index: %d"
		t.Fatalf(errmsg, index)
	}
}

func TestExpectDataByte(t *testing.T) {
	fmt.Println("TestExpectDataByte")
	value, err := expectDataByte(testBuffer, 0)
	if err != nil {
		errmsg := "Unexpected error at index 0: %s"
		t.Fatalf(errmsg, err)
	}
	if value != testBuffer[0] {
		errmsg := "Expected data byte 0x00, got 0x%02X"
		t.Fatalf(errmsg, value)
	}
	_, err = expectDataByte(testBuffer, 2)
	if err == nil {
		errmsg := "Did not detect non-data byte at index 2"
		t.Fatalf(errmsg)
	}
	_, err = expectDataByte(testBuffer, len(testBuffer))
	if err == nil {
		errmsg := "Did not detect out of bounds index: %d"
		t.Fatalf(errmsg, len(testBuffer))
	}
}

// Returns non-nil error if template and index are not identical.
//
func matchArray(index int, template []byte, sample []byte) (err error) {
	errmsg := "At index %d, Expected bytes %v, got %v"
	if len(sample) != len(template) {
		err = fmt.Errorf(errmsg, index, template, sample)
		return
	}
	for i, a := range template {
		b := sample[i]
		if a != b {
			err = fmt.Errorf(errmsg, index, template, sample)
			return
		}
	}
	return
}


func TestExpectRunningStatus(t *testing.T) {
	fmt.Println("TestExpectRunningStatus")
	index := 20
	mdata, newIndex, err := expectRunningStatus(testBuffer,0x80, index)
	if err != nil {
		errmsg := "Got unexpected error: %s"
		t.Fatalf(errmsg, err)
	}
	if newIndex != index+2 {
		errmsg := "Expected newIndex %2, got %d"
		t.Fatalf(errmsg, index+2, newIndex)
	}
	err = matchArray(index, []byte{0x80, 0x00, 0x7f}, mdata)
	if err != nil {
		t.Fatalf("%s", err)
	}

	_, _, err = expectRunningStatus(testBuffer, 0x80, 23)
	if err == nil {
		errmsg := "Did not detect malformed running-status at index 23"
		t.Fatalf(errmsg)
	}

	// single byte running status
	mdata, newIndex, err = expectRunningStatus(testBuffer, 0xc0, 26)
	if err != nil {
		errmsg := "Got unexpected error: %s"
		t.Fatalf(errmsg, err)
	}
	if newIndex != 27 {
		errmsg := "Expected newIndex 27, got %d"
		t.Fatalf(errmsg, newIndex)
	}
	err = matchArray(26, []byte{0xC0, 0x01}, mdata)
	if err != nil {
		t.Fatalf("%s", err)
	}

	// check index bounds
	_, _, err = expectRunningStatus(testBuffer, 0x80, 53)
	if err == nil {
		errmsg := "Did not detect out of bounds index"
		t.Fatalf(errmsg)
	}	
}
	
func TestExpectChannelMessage(t *testing.T) {
	fmt.Println("TestExpectChannelMessage")
	index := 19
	mdata, newIndex, err := expectChannelMessage(testBuffer, 0x80, index)
	if err != nil {
		errmsg := "Got unexpected error: %s"
		t.Fatalf(errmsg, err)
	}
	if newIndex != index + 3 {
		errmsg := "Expected newIndex %d, got %d"
		t.Fatalf(errmsg, index+3, newIndex)
	}
	err = matchArray(index, []byte{0x80, 0x00, 0x7f}, mdata)
	if err != nil {
		t.Fatalf("%s", err)
	}
}


func TestExpectSysexMessage(t *testing.T) {
	fmt.Println("TestExpectSysexMessage")
	index := 27
	mdata, newIndex, err := expectSysexMessage(testBuffer, index)
	if err != nil {
		errmsg := "Got unexpected error: %s"
		t.Fatalf(errmsg, err)
	}
	if newIndex != 32 {
		errmsg := "Expected newIndex 32, got %d"
		t.Fatalf(errmsg, index+3, newIndex)
	}
	err = matchArray(index, []byte{0xf0, 0x00, 0x01, 0x02, 0xF7}, mdata)
	if err != nil {
		t.Fatalf("%s", err)
	}

	index = 32
	_, _, err = expectSysexMessage(testBuffer, index)
	if err == nil {
		t.Fatalf("Did not detect malformed sysex message starting at index %d", index)
	}

	index = 28
	_, _, err = expectSysexMessage(testBuffer, index)
	if err == nil {
		t.Fatalf("Did not detect missing sysex status at index %d", index)
	}
}


func TestExpectSystemMessage(t *testing.T) {
	fmt.Println("TestExpectSystemMessage")
	index := 35
	mdata, newIndex, err := expectSystemMessage(testBuffer, index)
	if err != nil {
		t.Fatalf("Got unexpected errorL %s", err)
	}
	if newIndex != index+1 {
		errmsg := "Expected newIndex %d, got %d"
		t.Fatalf(errmsg, index+1, newIndex)
	}
	err = matchArray(index, []byte{0xF8}, mdata)
	if err != nil {
		t.Fatalf("%s", err)
	}

	index = 18
	_, _, err = expectSystemMessage(testBuffer, index)
	if err == nil {
		t.Fatalf("Did not detect non-system message at index %d", index)
	}
}


func TestExpectMetaMessage(t *testing.T) {
	fmt.Println("TestExpectMetaMessage")
	var err error
	var index, newIndex int
	var mdata []byte

	index = 41
	mdata, newIndex, err = expectMetaMessage(testBuffer, index)
		if err != nil {
		t.Fatalf("Got unexpected errorL %s", err)
	}
	if newIndex != 48 {
		errmsg := "Expected newIndex 48, got %d"
		t.Fatalf(errmsg, newIndex)
	}
	err = matchArray(index, []byte{0xff, 0x01, 0x04, 0x41, 0x42, 0x43, 0x44}, mdata)
	if err != nil {
		t.Fatalf("%s", err)
	}

	index = 48
	_, _, err = expectMetaMessage(testBuffer, index)
	if err == nil {
		t.Fatalf("Did not detect malformed meta message at index %d", index)
	}
}
