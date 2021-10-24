package smf

/*
** take.go defines functions for extracting values from byte slices.
** Includes miscellaneous utility functions.
*/

import (
	"fmt"
	"github.com/plewto/pigiron/midi"
	"github.com/plewto/pigiron/pigerr"
)


// msb(n) function returns upper byte of 16-bit value.
//
func msb(n int) byte {
	hi := (n & 0xFF00) >> 8
	return byte(hi)
}

// lsb(n) function returns the lower byte of 16-bit value.
//
func lsb(n int) byte {
	return byte(n & 0x00FF)
}

// reverse order of byte slice.
//
func reverse(src []byte) []byte {
	acc := make([]byte, len(src))
	for i, j := len(src)-1, 0; j < len(src); i, j = i-1, j+1 {
		acc[j] = src[i]
	}
	return acc
}


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
		
// TakeByte() returns first byte in buffer.
// Returns:
//   value - the first byte of buffer
//   newBuffer - slice of buffer starting after first byte.
//   error - non-nil if buffer is empty.
//
func TakeByte(buffer []byte) (value byte, newBuffer []byte, err error) {
	err = requireBufferLength(buffer, 1)
	if err != nil {
		return 0, []byte{}, err
	}
	return buffer[0], buffer[1:], err
}

func TakeStatusByte(buffer []byte) (value midi.StatusByte, newBuffer []byte, err error) {
	var bvalue byte
	bvalue, newBuffer, err = TakeByte(buffer)
	return midi.StatusByte(bvalue), newBuffer, err
}


// TakeShort() returns first two buffer bytes as 16-bit int.
// Returns:
//   value - 16-bit 'short' value
//   newBuffer - slice of buffer starting at index 2.
//   error - non-nil if buffer length less then 2.
//
func TakeShort(buffer []byte) (value int, newBuffer []byte, err error) {
	err = requireBufferLength(buffer, 2)
	if err != nil {
		return 0, []byte{}, err
	}
	b1, b2 := int(buffer[0]), int(buffer[1])
	value = b1 << 8 | b2
	return value, buffer[2:], err
}


// TakeLong() returns first four buffer bytes as 32-bit int.
// Returns:
//    value - 32-bit 'long' value
//    newBuffer - slice of buffer starting after index 4.
//    error - non-nil if buffer is not at least 4-bytes long.
//
func TakeLong(buffer []byte) (value int, newBuffer []byte, err error) {
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

// TakeVLQ() returns variable length value from start of buffer.
// The maximum number of bytes consumed is 4.
// Returns:
//    vlq - the 'value'
//    newBuffer - slice of buffer after final vlq byte.
//    error - non-nil if vlq is not terminated after reading 4 bytes.
//
func TakeVLQ(buffer []byte) (vlq *VLQ, newBuffer []byte, err error) {
	vlq = new(VLQ)
	var maxBytes = 4
	var acc = make([]byte, 0, maxBytes)
	for i := 0; i < maxBytes; i++ {
		if i >= len(buffer) {
			errmsg := "expect.TakeVLQ index out of bounds, "
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

