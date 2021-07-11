package smf

import (
	"fmt"
	// "github.com/plewto/pigiron/pigerror"
)


// expectChunkID checks byte buffer for specific chunk ID.
// buffer - the data
// index - location in buffer to be checked.
// target - expect ID pattern
// Returns non-nil error if expected ID not found.
//
func expectChunkID(buffer []byte, index int, target [4]byte) error {
	var err error
	for i, j := index, 0; j < 4; i, j = i+1, j+1 {
		if i > len(buffer) {
			msg := "smf.expectID index out of bounds, index = %d"
			err = fmt.Errorf(msg, i)
			return err
		}
		if buffer[i] != target[j] {
			id := ""
			for _, c := range target {
				id += fmt.Sprintf("%c", c)
			}
			msg := "smf.expectID, expected chunk id '%v' not found"
			err = fmt.Errorf(msg, id)
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
		msg := "smf.getShort() index out of range: index = %d, buffer length = %d"
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
func getByte(buffer []byte, index int) (byte, error) {
	var err error
	if len(buffer) <= index {
		msg := "smf.getByte() index out of range: index = %d, buffer length = %d"
		err = fmt.Errorf(msg, index, len(buffer))
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
	vlq.setBytes(acc)
	return vlq, err
}



// index points to first data byte
// 2nd (int) return points to new head.
//
func getRunningStatusMessage(buffer []byte, index int, st StatusByte, ch byte)(*ChannelMessage, int, error) {
	var startIndex = index
	var err error
	var cmsg *ChannelMessage
	var data1, data2 byte
	count, flag := channelStatusDataCount[StatusByte(st)]
	if !flag {
		msg := "smf.getRunningStatusMessage, non-channnel status: 0x%x '%s' at index %d"
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
		msg := "smf.GetRunningStatysMessage switch fall through, should never be here.\n"
		msg += fmt.Sprintf("status was 0x%x '%s', index was %d", st, st, index)
		err = fmt.Errorf(msg)
		panic(err)
	}
	var err2 error
	cmsg, err2 = NewChannelMessage(st, ch, data1, data2)
	if err2 != nil {
		msg := fmt.Sprintf("%s\n", err2)
		msg += fmt.Sprintf("smf.GetRunningStatusMessage, status was 0x%x '%s', index was %d", st, st, startIndex)
		err = fmt.Errorf(msg)
	}
	return cmsg, index, err
}
		
		
func getChannelMessage(buffer []byte, index int)(*ChannelMessage, int, error) {
	var startIndex = index
	var err error
	var cmsg *ChannelMessage
	var data1, data2 byte
	var sbyte byte
	sbyte, err = getByte(buffer, index)
	index++
	if err != nil {
		// msg := "smf.getChannelMessage index was %d"
		// err.Add(fmt.Sprintf(msg, startIndex))
		// return cmsg, startIndex, err
		msg := "%s\nsmf.getChannelMessage index was %d"
		err = fmt.Errorf(msg, err, startIndex)
		return cmsg, startIndex, err
	}
	st := StatusByte(sbyte & 0xF0)
	ch := sbyte & 0x0F
	count, flag := channelStatusDataCount[st]
	if !flag {
		msg := "smf.getRunningStatusMessage, non-channnel status: 0x%x '%s' at index %d"
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
		msg := "smf.GetChannelMessage switch fall through, should never be here.\n"
		msg += fmt.Sprintf("status was 0x%x '%s', index was %d", st, st, index)
		err = fmt.Errorf(msg)
		panic(err)
	}
	if err != nil {
		msg := "smf.getRunningStatusMessage, data index %d out of bounds\n"
		msg += fmt.Sprintf("%s\n", err)
		err = fmt.Errorf(msg, startIndex)
		return cmsg, startIndex, err
	}
	var err2 error
	cmsg, err2 = NewChannelMessage(st, ch, data1, data2)	
	if err2 != nil {
		msg := fmt.Sprintf("smf.GetRunningStatusMessage, status was 0x%x '%s', index was %d", st, st, startIndex)
		msg += fmt.Sprintf("\n%s", err2)
		err = fmt.Errorf(msg)
	}
	return cmsg, index, err
}



	
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
		end = findNextStatusByte(buffer, start+1)
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
		msg += fmt.Sprintf("%s", err2)
		err = fmt.Errorf(msg)
		return sys, start, err
	}
	return sys, end+1, err
}
	

