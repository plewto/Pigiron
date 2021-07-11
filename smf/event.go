package smf

import (
	"fmt"
)

type Event struct {
	time float64
	message MIDIMessage
}

type EventList struct {
	tempo float64   // BPM
	division int64  // clocks per quarter note
	events []*Event
}

func (e *EventList) tickDuration() float64 { // ISSUE: Implement
	return 0.01
}
	

func (e *EventList) getDeltaTime(vlq *VLQ) float64 {
	tick := e.tickDuration()
	return tick * float64(vlq.Value())
}
	


func CreateEventList(division int64, bytes []byte) (*EventList, error) {
	var err error
	var eventList *EventList
	var acc []*Event = make([]*Event, 0, 512)
	var useRunningStatus bool = false
	var runningChannel byte = 0
	var runningStatus StatusByte = 0
	var currentTime float64 = 0

	for index := 0; index < len(bytes);  {
		var vlq *VLQ
		vlq, err = getVLQ(bytes, index)
		if err != nil {
			msg := "smf.CreateEventList expected VLQ at index %d\n%s"
			err = fmt.Errorf(msg, index, err)
			return eventList, err
		}
		delta := eventList.getDeltaTime(vlq)
		currentTime += delta
		index += vlq.Length()
		if index > len(bytes) {
			break
		}
		status := bytes[index]
		switch {
		case !isStatusByte(status):
			if !useRunningStatus {
				msg := "smf.CreateEventList expected running status at index %d"
				err = fmt.Errorf(msg, index)
				return eventList, err
			}
			var cmsg *ChannelMessage
			rs, rc  := runningStatus, runningChannel
			cmsg, index, err = getRunningStatusMessage(bytes, index, rs, rc)
			if err != nil {
				return eventList, err
			}
			evnt := &Event{currentTime, cmsg}
			acc = append(acc, evnt)
		case isChannelStatus(status & 0xF0):
			// useRunningStatus = true
			// runningStatus = int(status & 0xF0)
			// runningChannel = int(status & 0x0F)
			// dataCount, _ := channelStatusByteCount[StatusByte(runningStatus)]
			// var data1, data2 byte
			// if dataCount == 3 {
			// 	data1, data2 = bytes[index+1], bytes[index+2]  // issue index not checked
			// 	index += 2
			// } else {
			// 	data1 = bytes[index+1] // issue index not checked
			// 	index += 1
			// }
			// rs, rc := StatusByte(runningStatus), byte(runningChannel)
			// chanmsg, _ := NewChannelMessage(rs, rc, data1, data2)
			// evnt := &Event{currentTime, chanmsg}
			// acc = append(acc, evnt)
		case isSystemStatus(status):
			// useRunningStatus = false
			// index++
			// b := bytes[index] // issue index not checked
			// bcc := make([]byte, 0, 16)
			// for index < len(bytes) && (b & 0x80 == 0) {
			// 	bcc = append(bcc, b)
			// 	index++
			// }
			// smsg, _ := newSystemMessage(bcc)
			// evnt := &Event{currentTime, smsg}
			// acc = append(acc, evnt)
		case isMetaStatus(status):
			useRunningStatus = false
			// handle meta message
		default:
			msg := "smf.CreateEventList expected staus byte at index %d, got 0x%x"
			err = fmt.Errorf(msg, index, status)
			return eventList, err
		}
		
				
	} // end for index
	eventList.events = acc
	return eventList, err
} // end func

