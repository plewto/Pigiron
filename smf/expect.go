package smf

import (
	"fmt"
)


// expectChunkID checks byte buffer for specific chunk ID.
// buffer - smf bytes.
// index - location in buffer to be checked.
// target - expect ID pattern
// Returns non-nil error if expected ID not found.
//
func expectChunkID(buffer []byte, index int, target [4]byte) error {
	var err error
	for i, j := index, 0; j < 4; i, j = i+1, j+1 {
		if i > len(buffer) {
			msg := "smf.expectID index out of bounds, index = %d"
			err = exError(msg)
			return err
		}
		if buffer[i] != target[j] {
			id := ""
			for _, c := range target {
				id += fmt.Sprintf("%c", c)
			}
			msg := fmt.Sprintf("smf.expectID, expected chunk id '%v' not found", id)
			err = exError(msg)
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
		msg := "smf.getLong index out of range: index = %d, buffer length = %d"
		msg = fmt.Sprintf(msg, index, len(buffer))
		err = exError(msg)
		return 0, err
	}
	acc := 0
	for i, j, shift := index, 0, 24; j < 4; i, j, shift = i+1, j+1, shift-8 {
		n := int(buffer[i])
		acc += int(n << shift)
	}
	return acc, err
}

func get3Bytes(buffer []byte, index int)  (int, error) {
	var err error
	if index > len(buffer) + 3 {
		msg := "smf.get3Bytes index out of range: index = %d, buffer length = %d"
		err = exError(fmt.Sprintf(msg, index, len(buffer)))
		return 0, err
	}
	acc := 0
		for i, j, shift := index, 0, 16; j < 3; i, j, shift = i+1, j+1, shift-8 {
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
		msg := "smf.getShort() index out of range: index = %d, buffer length = %d"
		msg = fmt.Sprintf(msg, index, len(buffer))
		err = exError(msg)
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
func getByte(buffer []byte, index int) (byte, error) {
	var err error
	if len(buffer) <= index {
		msg := "smf.getByte() index out of range: index = %d, buffer length = %d"
		msg = fmt.Sprintf(msg, index, len(buffer))
		err = exError(msg)
		return 0, err
	}
	return buffer[index], err
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
			msg := "smf.getVLQ index out of bounds, index = %d, buffer length = %d"
			msg = fmt.Sprintf(msg, index, len(buffer))
			err = exError(msg)
			return vlq, err
		}
		n := buffer[index]
		acc = append(acc, n)
		if n & 0x80 == 0 {
			break
		}
		index++
	}
	vlq.setBytes(acc)
	return vlq, err
}


// getRunningStatuscreates MIDI message from smf track data using 'running status'.
// buffer - smf track data
// index - index should point to first data byte.
// st - MIDI status byte
// ch - MIDI channel byte
//
// The 2nd return value is an updated index pointing to the start of the 
// following event's delta-time.
//
func getRunningStatus(buffer []byte, index int, st StatusByte, ch byte)(*ChannelMessage, int, error) {
	var startIndex = index
	var err error
	var cmsg *ChannelMessage
	var data1, data2 byte
	count, flag := channelStatusDataCount[StatusByte(st)]
	if !flag {
		msg := "smf getRunningStatus, runing-statys error\n"
		msg += "Expected non-status byte at index %d, got 0x%x"
		err = exError(fmt.Sprintf(msg, index, st))
		return cmsg, index, err
	}
	switch count {
	case 1:
		data1, err = getByte(buffer, index)
		index++
	case 2:
		data1, _ = getByte(buffer, index)
		data2, err = getByte(buffer, index+1)
		index += 2
	default:
		msg := "smf getRunningStatus,  switch fell through, should never be here.\n"
		msg += fmt.Sprintf("status was 0x%x '%s', index was %d", st, st, index)
		err = fmt.Errorf(msg)
		panic(err)
	}
	var err2 error
	cmsg, err2 = NewChannelMessage(st, ch, data1, data2)
	if err2 != nil {
		msg := "smf getRunningStatus, status was 0x%x '%s', index was %d"
		msg = fmt.Sprintf(msg, st, st, startIndex)
		err = compoundError(err2, msg)
	}
	return cmsg, index, err
}
		
// getChannelMessage creates MIDI ChannelMessage from smf track data.
// buffer - smf track data
// index - location of message status byte.
//
// The 2nd return value is an updated index pointing to the start of the 
// following event's delta-time.
//
func getChannelMessage(buffer []byte, index int)(*ChannelMessage, int, error) {
	var startIndex = index
	var err error
	var cmsg *ChannelMessage
	var data1, data2 byte
	var sbyte byte
	sbyte, err = getByte(buffer, index)
	index++
	if err != nil {
		msg := "%s\nsmf.getChannelMessage index was %d"
		err = fmt.Errorf(msg, err, startIndex)
		return cmsg, startIndex, err
	}
	st := StatusByte(sbyte & 0xF0)
	ch := sbyte & 0x0F
	count, flag := channelStatusDataCount[st]
	if !flag {
		msg := "smf getRunningStatus, non-channel status: 0x%x '%s' at index %d"
		err = fmt.Errorf(msg, st, st, index)
	        return cmsg, index, err
	}
	switch count {
	case 1:
		data1, err = getByte(buffer, index)
		index++
	case 2:
		data1, _ = getByte(buffer, index)
		data2, err = getByte(buffer, index+1)
		index += 2
	default:
		msg := "smf getChannelMessage switch fall through, should never be here.\n"
		msg += fmt.Sprintf("status was 0x%x '%s', index was %d", st, st, index)
		err = fmt.Errorf(msg)
		panic(err)
	}
	if err != nil {
		msg := "smf getChannelMessage, data index %d out of bounds\n"
		msg += fmt.Sprintf("%s\n", err)
		err = fmt.Errorf(msg, startIndex)
		return cmsg, startIndex, err
	}
	var err2 error
	cmsg, err2 = NewChannelMessage(st, ch, data1, data2)	
	if err2 != nil {
		msg := "smf.GetRunningStatusMessage, status was 0x%x '%s', index was %d"
		msg = fmt.Sprintf(msg, st, st, startIndex)
		err = compoundError(err2, msg)
	}
	return cmsg, index, err
}


// findNextStatusByte locates the next status byte after the starting location.
// buffer - smf track data
// start - index where search commences,  start should not point to a status byte.
//
// returns index of first status byte (bit-7 set) following start.
// 
func findNextStatusByte(buffer []byte, start int) int {
	index := start
	var b byte
	for index < len(buffer) {
		b = buffer[index]
		if b & 0x80 == 0x80 {
			break
		}
		index++
	}
	if index > len(buffer) {
		index = len(buffer)
	}
	return index
}

// getSystemMessage creates a MIDI system-message from smf track bytes
// buffer - smf track data
// index - location of system message status byte
//
// The 2nd return value is an updated index pointing to the start of the 
// following event's delta-time.
//
func getSystemMessage(buffer []byte, index int)(*SystemMessage, int, error) {
	var start = index
	var end = start
	var err error
	var sys *SystemMessage
	var status byte
	status, err = getByte(buffer, index)
	if err != nil {
		msg := "smf.getSystemMessage\n"
		msg += fmt.Sprintf("%s", err)
		err = fmt.Errorf(msg)
		return sys, start, err
	}
	if !isSystemStatus(status) {
		msg := "smf.getSystemMessage expected system status byte at index %d, got 0x%02x, '%s'"
		err = fmt.Errorf(msg, index, status, StatusByte(status))
		return sys, start, err
	}
	count, _ := systemStatusDataCount[StatusByte(status)]
	if count == -1 {
		// assume sysex
		end = findNextStatusByte(buffer, start+1) + 1
	} else {
		// assume 0
		end = start+1
	}
	if start > end || end > len(buffer) {
		msg := "smf.getSystemMessage bytes slice indexes are wrong.\n"
		msg += "index = %d,  start = %d,  end = %d,  buffer length is %d"
		err = fmt.Errorf(msg, index, start, end, len(buffer))
		return sys, start, err
	}
	var bytes []byte = buffer[start:end]
	var err2 error
	sys, err2 = newSystemMessage(bytes)
	if err2 != nil {
		msg := fmt.Sprintf("smf.getSystemMessage index = %d\n", index)
		err = compoundError(err2, msg)
		return sys, start, err
	}
	return sys, end, err
}


// getMetaMessage creates MIDI meta message from smf track bytes.
// buffer - smf track data
// index - location of meta status byte
//
// The 2nd return value is an updated index pointing to the start of the 
// following event's delta-time.
//
func getMetaMessage(buffer []byte, index int) (*MetaMessage, int, error) {
	startIndex := index
	var err error
	var meta *MetaMessage
	var st, mtype byte
	var vlq *VLQ
	st, err = getByte(buffer, index)
	if err != nil {
		msg := fmt.Sprintf("smf.getMetaMessage (status byte) index = %d", index)
		err = compoundError(err, msg)
		return meta, startIndex, err
	}
	index++
	mtype, err = getByte(buffer, index)
	index++
	if err != nil {
		msg := fmt.Sprintf("smf.getMetaMessage (meta-type byte) index = %d", index)
		err = compoundError(err, msg)
		return meta, startIndex, err
	}
	if StatusByte(st) != MetaStatus {
		msg := "smf.getMetaMessage, expected Meta status byte 0xFF, got 0x%x"
		err = fmt.Errorf(msg, st)
		return meta, index, err
	}
	if !isMetaType(mtype) {
		msg := "smf.getMetaMessage, expected valid MetaType, got 0x%x"
		err = fmt.Errorf(msg, mtype)
		return meta, index, err
	}
	vlq, err = getVLQ(buffer, index)
	if err != nil {
		msg := fmt.Sprintf("smf.getMetaMessage, index = %d", index)
		err = compoundError(err, msg)
		return meta, index, err
	}
	dataCount := vlq.Value()
	index += vlq.Length()
	start, end := index, index+dataCount
	if start > end || end > len(buffer) {
		msg := "smf.getMetaMessage, data slice index out of bounds.  "
		msg += "start = %d, end = %s, buffer length is %d"
		err = fmt.Errorf(msg, start, end, len(buffer))
		return meta, index, err
	}
	data := buffer[start:end]
	meta, err = newMetaMessage(MetaType(mtype), data)
	if err != nil {
		msg := fmt.Sprintf("smf.getMetaMessage, index = %d", startIndex)
		err = compoundError(err, msg)
		return meta, index, err
	}
	
	return meta, end, err
}

