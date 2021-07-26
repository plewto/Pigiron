package midi

/*
** take.go defines functions for extracting values from byte slices.
**
*/

import (
	"fmt"
	"github.com/plewto/pigiron/pigerr"
)


// requireBufferLength() returns error if buffer is not at least count bytes long.
//
func requireBufferLength(buffer []byte, count int) error {
	var err error
	if len(buffer) < count {
		errmsg := "Require byte buffer of at least %d bytes, got %d"
		err = pigerr.New(fmt.Sprintf(errmsg, count, len(buffer)))
	}
	return err
}
		
// takeByte() returns first byte in buffer.
// Returns:
//   value - the first byte of the buffer
//   newBuffer - slice of buffer starting after first byte.
//   error - non-nil if buffer is empty.
//
func takeByte(buffer []byte) (value byte, newBuffer []byte, err error) {
	err = requireBufferLength(buffer, 1)
	if err != nil {
		return 0, []byte{}, err
	}
	return buffer[0], buffer[1:], err
}


// takeShort() returns first two buffer bytes as 16-bit int.
// Returns:
//   value - 16-bit 'short' value
//   newBuffer - slice of buffer starting at index 2
//   error - non-nil if buffer length is less then 2.
//
func takeShort(buffer []byte) (value int, newBuffer []byte, err error) {
	err = requireBufferLength(buffer, 2)
	if err != nil {
		return 0, []byte{}, err
	}
	b1, b2 := int(buffer[0]), int(buffer[1])
	value = b1 << 8 | b2
	return value, buffer[2:], err
}


// takeLong() returns first four buffer bytes as 32-bit int.
// Returns:
//    value - 32-bit 'long' value
//    newBuffer - slice of buffer starting after index 4
//    error - non-nil if buffer is not at least 4-bytes long.
//
func takeLong(buffer []byte) (value int, newBuffer []byte, err error) {
	err = requireBufferLength(buffer, 4)
	if err != nil {
		return 0, []byte{}, err
	}
	value = 0
	for i, shift := 0, 24; i < 4; i, shift = i+1, shift-8 {
		n := int(buffer[i])
		value += n << shift
	}
	return value, buffer[4:], err
}

// takeVLQ() returns variable length value from start of buffer.
// The maximum number of bytes consumed is 4.
// Returns:
//    vlq - the 'value'
//    newBuffer - slice of buffer after final vlq byte.
//    error - non-nil if vlq is not terminated after reading 4 bytes.
//
func takeVLQ(buffer[] byte) (vlq *VLQ, newBuffer []byte, err error) {
	vlq = new(VLQ)
	var maxBytes = 4
	var acc = make([]byte, 0, maxBytes)
	for i := 0; i < maxBytes; i++ {
		if i >= len(buffer) {
			errmsg := "smf.takeVLQ index out of bounds, "
			errmsg += "index = %d, buffer length = %d"
			err = pigerr.New(fmt.Sprintf(errmsg, i, len(buffer)))
			return vlq, []byte{}, err
		}
		n := buffer[i]
		acc = append(acc, n)
		if n & 0x80 == 0 {
			break
		}
	}
	vlq.setBytes(acc)
	return vlq, buffer[len(acc):], err
}
