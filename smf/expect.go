package smf


/*
 * expect.go defines functions to convert byte array to MIDI Events.
 *
*/

import (
	"fmt"
	"github.com/plewto/pigiron/midi"
)

func expectByte(buffer []byte, index int) (value byte, newIndex int, err error) {
	if index >= len(buffer) {
		errmsg := "expectByte, index out of bounds: %d"
		err = fmt.Errorf(errmsg, index)
		return
	}
	value = buffer[index]
	newIndex = index + 1
	return
}


func expectVLQ(buffer []byte, index int) (vlq *VLQ, newIndex int, err error) {
	var maxBytes = 4
	var acc = make([]byte, 0, maxBytes)
	count := 0
	for i := index; i < len(buffer); i++ {
		count++
		if count > maxBytes {
			errmsg := "Expected VLQ at index %d"
			err = fmt.Errorf(errmsg, index)
			return
		}
		n := buffer[i]
		acc = append(acc, n)
		if n & 0x80 == 0 {
			break
		}
	}
	vlq.setBytes(acc)
	newIndex = index + count
	return
}
	


// stch must contain both hi and low (command & channel) status-byte nibbles.
// index -> first data byte of message
//
func expectChannelMessage(buffer []byte, stch byte, index int) (mdata []byte, newIndex int, err error) {
	count := midi.ChannelMessageDataCount(midi.StatusByte(stch))
	if len(buffer) < index + count {
		errmsg := "expectChannelMessage index out of bounds. index = %d, status = 0x%02X"
		err = fmt.Errorf(errmsg, index, stch)
		return
	}

	var assertDataByte = func(d byte, i int) error {
		var derr error
		if d > 0x7F {
			errmsg := "Expected MIDI data byte at index %d, got 0x%02X"
			derr = fmt.Errorf(errmsg, i, d)
		}
		return derr
	}
	
	if count == 1 {
		d1 := buffer[index]
		err = assertDataByte(d1, index)
		if err != nil {
			return
		}
		mdata = []byte{stch, d1}
		newIndex = index + 1
		return
	} else { // assume 2 data bytes
		d1 := buffer[index]
		d2 := buffer[index + 1]
		err = assertDataByte(d1, index)
		if err != nil {
			return
		}
		err = assertDataByte(d2, index)
		if err != nil {
			return
		}
		mdata = []byte{stch, d1, d2}
		newIndex = index + 2
		return
	}
}
		
// index -> first data bytre after status
//
func expectSysexMessage(bytes []byte, index int) (mdata []byte, newIndex int, err error) {
	var b byte
	newIndex = index
	for b & 0xfF > 0x7f && newIndex < len(bytes) {
		b = bytes[newIndex]
		newIndex++
	}
	mdata = bytes[index-1:newIndex]
	return
}

// index -> first byte after status
//
func expectSystemMessage(bytes []byte, index int) (mdata []byte, newIndex int, err error) {
	st := bytes[index-1]
	newIndex = index
	mdata = []byte{st}
	return
}


// index -> metaType byte
//
func expectMetaMessage(bytes []byte, index int) (mdata []byte, newIndex int, err error) {
	var i = index
	var mtype byte
	if len(bytes) <= i {
		errmsg := "expectMetaMessage index out of bounds, start-index = %d"
		err = fmt.Errorf(errmsg, index-1)
		return
	}
	mtype = byte(i)
	if !midi.IsMetaType(midi.MetaType(mtype)) {
		errmsg := "Expected Meta type bytes at index %d, got 0x%02X"
		err = fmt.Errorf(errmsg, i, mtype)
		return
	}

	i++
	var vlq *VLQ
	vlq, i, err = expectVLQ(bytes, i)
	if err != nil {
		errmsg := fmt.Sprintf("Expected meta VLQ at index %d", i-1)
		err = fmt.Errorf("%s\n%s", errmsg, err)
		return
	}
	mdata = make([]byte, 2, 2+vlq.Length()+vlq.Value())
	mdata[0] = 0xFF
	mdata[1] = mtype
	for _, b := range vlq.Bytes() {
		mdata = append(mdata, b)
		i++
	}
	for j, count := i, 0; count < vlq.Value(); j, count = j+1, count+1 {
		if j >= len(bytes) {
			errmsg := "expectMetaMessage index out of bounds, start-index = %d"
			err = fmt.Errorf(errmsg, index)
			return
		}
		mdata = append(mdata, bytes[j])
		i++
	}
	newIndex = i
	return
}
