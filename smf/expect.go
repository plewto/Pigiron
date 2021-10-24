package smf


/*
 * expect.go defines functions to convert byte array to MIDI Events.
 *
*/

import (
	"fmt"
	"github.com/plewto/pigiron/midi"
)

// const (
// 	KEYED_STATUS byte = 0x01
// 	CLOCK = byte 0xF8
// 	START = byte 0xFA
// 	CONTINUE = byte 0xFB
// 	STOP = byte 0xFC
// 	META = byte 0xFF
// )


// func channelMessageDataCount(st byte) int {
// 	hi := st & 0xF0
// 	if hi == 0xC0 || hi == 0xd0 {
// 		return 1
// 	} else {
// 		return 2
// 	}
// }

// func isChannelStatus(st byte) bool {
// 	return st == midi.KEYED_STATUS || (st & 0xF0) < 0xF0
// }
	

// func isSystemRealtimeStatus(st byte) bool {
// 	return st == midi.CLOCK || st == START || st == CONTINUE || st == STOP
// }

// func isMetaStatus(st byte) bool {
// 	return st == META
// }

// func isMetaType(mt byte) bool {
// 	switch {
// 	case 0x00 <= mt && mt <= 0x07:
// 		return true
// 	case mt == 0x20 || mt == 0x2F:
// 		return true
// 	case mt == 0x51 || mt == 0x54 || mt == 0x58 || mt == 0x59:
// 		return true
// 	case mt == 0x7F:
// 		return true
// 	default:
// 		return false
// 	}
// }
		
// func isSystemStatus(st byte) bool {
// 	return !(isMetaStatus(st) || isChannelStatus(st))
// }

func ExpectByte(buffer []byte, index int) (value byte, newIndex int, err error) {
	if index >= len(buffer) {
		errmsg := "ExpectByte, index out of bounds: %d"
		err = fmt.Errorf(errmsg, index)
		return
	}
	value = buffer[index]
	newIndex = index + 1
	return
}


func ExpectVLQ(buffer []byte, index int) (vlq *VLQ, newIndex int, err error) {
	var maxBytes = 4
	var acc = make([]byte, 0, maxBytes)
	count := 0
	if index >= len(buffer) {
		errmsg := "ExpectVLQ index out of bounds: %d"
		err = fmt.Errorf(errmsg, index)
		return
	}
	for i := index; i < len(buffer); i++ {
		count++
		if count > maxBytes {
			errmsg := "expect.ExpectVLQ, expected VLQ at index %d"
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

func ExpectDataByte(buffer []byte, index int) (value byte, err error) {
	if index >= len(buffer) {
		errmsg := "expect.ExpectDataByte index out of bounds at %d"
		err = fmt.Errorf(errmsg, index)
		return
	}
	value = buffer[index]
	if value > 0x7F {
		errmsg := "expect.ExpectDataByte, expected MIDI data byte at index %d, got 0x%02X"
		err = fmt.Errorf(errmsg, index, value)
	}
	return
}


func ExpectRunningStatus(buffer []byte, status byte, index int) (mdata []byte, newIndex int, err error) {
	count := midi.ChannelMessageDataCount(midi.StatusByte(status))
	if len(buffer) <= index+count {
		errmsg := "expect.ExpectRunningStatus index out of bounds %d, []byte length is %d"
		err = fmt.Errorf(errmsg, index, len(buffer))
		return
	}
	var d1, d2 byte
	switch count {
	case 1:
		d1, err = ExpectDataByte(buffer, index)
		if err != nil {
			return
		}
		mdata = []byte{status, d1}
		newIndex = index + 1
	case 2:
		d1, err = ExpectDataByte(buffer, index)
		if err != nil {
			return
		}
		d2, err = ExpectDataByte(buffer, index+1)
		if err != nil {
			return
		}
		mdata = []byte{status, d1, d2}
		newIndex = index + 2
	default:
		errmsg := "expect.ExpectRunningStatus swtich fallthrough. Status byte was 0x%02X"
		err = fmt.Errorf(errmsg, status)
	}
	return
}

func ExpectChannelMessage(buffer []byte, status byte, index int) (mdata []byte, newIndex int, err error) {
	mdata, newIndex, err = ExpectRunningStatus(buffer, status, index+1)
	return
}


func ExpectSysexMessage(buffer []byte, index int) (mdata []byte, newIndex int, err error) {
	var acc = make([]byte, 1, 1024)
	if index >= len(buffer) {
		errmsg := "expect.ExpectSysexMessage index out of bounds at %d, []byte length is %d"
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
			errmsg := "expect.ExpectSysexMessage index out of bounds at %d, []byte length is %d"
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
func ExpectSystemMessage(buffer []byte, index int) (mdata []byte, newIndex int, err error) {
	if index >= len(buffer) {
		errmsg := "expect.ExpectSystemMessage index %d out of bounds, []byte length is %d"
		err = fmt.Errorf(errmsg, index, len(buffer))
		return
	}
	st := buffer[index]
	if !midi.IsSystemStatus(midi.StatusByte(st)) {
		errmsg := "Expected MIDI real time system message at index %d, got 0x%02X"
		err = fmt.Errorf(errmsg, index, st)
		return
	}
	mdata = []byte{byte(st)}
	newIndex = index+1
	return
}


func ExpectMetaMessage(buffer []byte, index int) (mdata []byte, newIndex int, err error) {
	if index >= len(buffer)-1 {
		errmsg := "expect.ExpectMetaMessage index %d out of bounds, []byte length is %d"
		err = fmt.Errorf(errmsg, index, len(buffer))
		return
	}
	st := buffer[index]
	mt := buffer[index+1]
	if !midi.IsMetaStatus(midi.StatusByte(st)) || !midi.IsMetaType(midi.MetaType(mt)) {
		errmsg := "Expected meta status 0xFF and valid meta type starting at index %d, "
		errmsg += "got 0x%02x and 0x$02x"
		err = fmt.Errorf(errmsg, index, byte(st), byte(mt))
		return
	}
	acc := make([]byte, 2, 128)
	acc[0] = byte(st)
	acc[1] = byte(mt)
	index += 2
	var vlq *VLQ
	vlq, index, err = ExpectVLQ(buffer, index)
	if err != nil {
		return
	}
	for _, b := range vlq.Bytes() {
		acc = append(acc, b)
	}
	for j, count := index, 0; count < vlq.Value(); j, count = j+1, count+1 {
		if j >= len(buffer) {
			errmsg := "expect.ExpectMetaMesage index %d out of bounds, []byte length is %d"
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
	
	
		
	
	
	
		
	
