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
	if index >= len(buffer) {
		errmsg := "expectVLQ index out of bounds: %d"
		err = fmt.Errorf(errmsg, index)
		return
	}
	for i := index; i < len(buffer); i++ {
		count++
		if count > maxBytes {
			errmsg := "smf.expectVLQ, expected VLQ at index %d"
			err = fmt.Errorf(errmsg, index)
			return
		}
		n := buffer[i]
		acc = append(acc, n)
		if n & 0x80 == 0 {
			break
		}
	}
	vlq = NewVLQ(0)
	vlq.setBytes(acc)
	newIndex = index + count
	return
}

func expectDataByte(buffer []byte, index int) (value byte, err error) {
	if index >= len(buffer) {
		errmsg := "smf.expectDataByte index out of bounds at %d"
		err = fmt.Errorf(errmsg, index)
		return
	}
	value = buffer[index]
	if value > 0x7F {
		errmsg := "smf.expectDataByte, expected MIDI data byte at index %d, got 0x%02X"
		err = fmt.Errorf(errmsg, index, value)
	}
	return
}


func expectRunningStatus(buffer []byte, status byte, index int) (mdata []byte, newIndex int, err error) {
	count := midi.ChannelMessageDataCount(midi.StatusByte(status))
	if len(buffer) <= index+count {
		errmsg := "smf.expectRunningStatus index out of bounds %d, []byte length is %d"
		err = fmt.Errorf(errmsg, index, len(buffer))
		return
	}
	var d1, d2 byte
	switch count {
	case 1:
		d1, err = expectDataByte(buffer, index)
		if err != nil {
			return
		}
		mdata = []byte{status, d1}
		newIndex = index + 1
	case 2:
		d1, err = expectDataByte(buffer, index)
		if err != nil {
			return
		}
		d2, err = expectDataByte(buffer, index+1)
		if err != nil {
			return
		}
		mdata = []byte{status, d1, d2}
		newIndex = index + 2
	default:
		errmsg := "smf.expectRunningStatus swtich fallthrough. Status byte was 0x%02X"
		err = fmt.Errorf(errmsg, status)
	}
	return
}

func expectChannelMessage(buffer []byte, status byte, index int) (mdata []byte, newIndex int, err error) {
	mdata, newIndex, err = expectRunningStatus(buffer, status, index+1)
	return
}


func expectSysexMessage(buffer []byte, index int) (mdata []byte, newIndex int, err error) {
	var acc = make([]byte, 1, 1024)
	if index >= len(buffer) {
		errmsg := "smf.expectSysexMessage index out of bounds at %d, []byte length is %d"
		err = fmt.Errorf(errmsg, index, len(buffer))
		return
	}
	st := buffer[index]
	if st != 0xF0 {
		errmsg := "Expected sysex status 0xF0 at index %d, got 0x02X"
		err = fmt.Errorf(errmsg, index, st)
	}
	var b byte
	acc[0] = st
	index++
	for {
		if index >= len(buffer) {
			errmsg := "smf.expectSysexMessage index out of bounds at %d, []byte length is %d"
			err = fmt.Errorf(errmsg, index, len(buffer))
			return
		}
		b = buffer[index]
		switch {
		case b < 0x80:
			acc = append(acc, b)
		case b == 0xF7:
			acc = append(acc, b)
			index++
			mdata = acc[0 : len(acc)]
			newIndex = index
			return
		case midi.IsSystemRealtimeStatus(midi.StatusByte(b)):
			// ignore
		case b >= 0x80:
			errmsg := "Sysex aborted by invalid status byte 0x%02X at index %d"
			err = fmt.Errorf(errmsg, b, index)
			return
		default:
			errmsg := "Sysex message aborted by unexpected status byte 0x%02X at index %d"
			err = fmt.Errorf(errmsg, b, index)
			return
		}
		index++
	}
	return
}
	
		

// handles non-sysex system messages
// Theses are all a single byte in length
//
func expectSystemMessage(buffer []byte, index int) (mdata []byte, newIndex int, err error) {
	if index >= len(buffer) {
		errmsg := "smf.expectSystemMessage index %d out of bounds, []byte length is %d"
		err = fmt.Errorf(errmsg, index, len(buffer))
		return
	}
	st := midi.StatusByte(buffer[index])
	if !midi.IsSystemStatus(st) {
		errmsg := "Expected MIDI real time system message at index %d, got 0x%02X"
		err = fmt.Errorf(errmsg, index, st)
		return
	}
	mdata = []byte{byte(st)}
	newIndex = index+1
	return
}


func expectMetaMessage(buffer []byte, index int) (mdata []byte, newIndex int, err error) {
	if index >= len(buffer)-1 {
		errmsg := "smf.expectMetaMessage index %d out of bounds, []byte length is %d"
		err = fmt.Errorf(errmsg, index, len(buffer))
		return
	}
	st := midi.StatusByte(buffer[index])
	mt := midi.MetaType(buffer[index+1])
	if !midi.IsMetaStatus(st) || !midi.IsMetaType(mt) {
		errmsg := "Expected meta status 0xFF and valid meta type starting at index %d, "
		errmsg += "got 0x%02x and 0x$02x"
		err = fmt.Errorf(errmsg, index, byte(st), byte(mt))
	}
	acc := make([]byte, 2, 128)
	acc[0] = byte(st)
	acc[1] = byte(mt)
	index += 2
	var vlq *VLQ
	vlq, index, err = expectVLQ(buffer, index)
	if err != nil {
		return
	}
	for _, b := range vlq.Bytes() {
		acc = append(acc, b)
		index++
	}
	for j, count := index, 0; count < vlq.Value(); j, count = j+1, count+1 {
		if j >= len(buffer) {
			errmsg := "smf.expectMetaMesage index %d out of bounds, []byte length is %d"
			err = fmt.Errorf(errmsg, j, len(buffer))
			return
		}
		acc = append(acc, buffer[j])
		index++
	}
	mdata = acc[0 : len(acc)]
	newIndex = index
	return
}
	
	
		
	
	
	
		
	
