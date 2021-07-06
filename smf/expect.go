package smf

import (
	"fmt"
)

// expectID checks byte buffer for specific chunk ID.
// buffer - the data
// index - location in buffer to be checked.
// target - expect ID pattern
// Returns non-nil error if expected ID not found.
//
func expectID(buffer []byte, index int, target [4]byte) error {
	var err error
	for i, j := index, 0; j < 4; i, j = i+1, j+1 {
		if i > len(buffer) {
			err = fmt.Errorf("ERROR: smf.expectID,  index out of bounds, index = %d", i)
			return err
		}
		if buffer[i] != target[j] {
			id := ""
			for _, c := range target {
				id += fmt.Sprintf("%c", c)
			}
			err = fmt.Errorf("ERROR: smf.expectID, expected chunk id '%v' not found", id)
			return err
		}
	}
	return err
}

// getLong extracts 4-byte value from buffer starting at index.
//
func getLong(buffer []byte, index int) (int, error) {
	var err error
	if len(buffer) < index+4 {
		msg := "ERROR smf.getLong() index out of range: index = %d, buffer length = %d"
		err = fmt.Errorf(msg, index, len(buffer))
		return 0, err
	}
	acc := 0
	for i, j, shift := index, 0, 24; j < 4; i, j, shift = i+1, j+1, shift-8 {
		n := int(buffer[i])
		acc += int(n << shift)
	}
	return acc, err
}

// getShort extracts 2-byte value from buffer starting at index.
//
func getShort(buffer []byte, index int) (int, error) {
	var err error
	if len(buffer) < index+2 {
		msg := "ERROR smf.getShort() index out of range: index = %d, buffer length = %d"
		err = fmt.Errorf(msg, index, len(buffer))
		return 0, err
	}
	acc := 0
	for i, j, shift := index, 0, 8; j < 2; i, j, shift = i+1, j+1, shift-8 {
		n := int(buffer[i])
		acc += int(n << shift)
	}
	return acc, err
}

// getByte extracts byte from buffer at index.
//
func getByte(buffer []byte, index int) (int, error) {
	var err error
	if len(buffer) <= index {
		msg := "ERROR smf.getByte() index out of range: index = %d, buffer length = %d"
		err = fmt.Errorf(msg, index, len(buffer))
		return 0, err
	}
	return int(buffer[index]), err
}

// getVLQ extracts variable-length-value starting at index.
// Between 1 and 4 bytes are used.
//
func getVLQ(buffer []byte, index int)(*VLQ, error) {
	var err error
	var vlq *VLQ = new(VLQ)
	var maxCount = 4
	var acc = make([]byte, 0, maxCount)
	for {
		if index == maxCount {
			break
		}
		if index >= len(buffer) {
			msg := "ERROR smf.getVLQ index out of bounds, index = %d, buffer length = %d"
			err = fmt.Errorf(msg, index, len(buffer))
			return vlq, err
		}
		n := buffer[index]
		acc = append(acc, n)
		if n & 0x80 == 0 {
			break
		}
		index++
	}
	vlq.SetBytes(acc)
	return vlq, err
}
