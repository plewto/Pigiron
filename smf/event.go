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
	var runChan byte = 0
	var runStat StatusByte = 0
	var currentTime float64 = 0 
	for index := 0; index < len(bytes);  {
		var vlq *VLQ
		vlq, err = getVLQ(bytes, index)
		if err != nil {
			errmsg := "smf.CreateEventList expected VLQ at index %d"
			err = compoundError(err, fmt.Sprintf(errmsg, index))
			return eventList, err
		}
		delta := eventList.getDeltaTime(vlq)
		currentTime += delta
		index += vlq.Length()
		if index > len(bytes) {
			break
		}
		status := bytes[index]
		startIndex := index
		switch {
		case !isStatusByte(status):
			var cmsg *ChannelMessage
			errmsg := "smf.CreateEventList running-status error\n"
			errmsg += "Expected non-status byte at index %d, got 0x%x"
			errmsg = fmt.Sprintf(errmsg, startIndex, byte(status))
			if !useRunningStatus {
				err = exError(errmsg)
				return eventList, err
			}
			cmsg, index, err = getRunningStatusMessage(bytes, startIndex, runStat, runChan)
			if err != nil {
				err = compoundError(err, errmsg)
				return eventList, err
			}
			evnt := &Event{currentTime, cmsg}
			acc = append(acc, evnt)
		case isChannelStatus(status & 0xF0):
			var cmsg *ChannelMessage
			useRunningStatus = true
			runStat = StatusByte(status & 0xF0)
			runChan = byte(status & 0x0F)
			cmsg, index, err = getChannelMessage(bytes, startIndex)
			if err != nil {
				errmsg := "smf.CreateEventList error in switch case isChannelMessage"
				err = compoundError(err, errmsg)
				return eventList, err
			}
			evnt := &Event{currentTime, cmsg}
			acc = append(acc, evnt)
		case isSystemStatus(status):
			var sys *SystemMessage
			useRunningStatus = false
			index++  // ISSUE is index correct
			sys, index, err = getSystemMessage(bytes, startIndex)
			if err != nil {
				errmsg := "smf.CreateEventList error in switch case isSystemStatus"
				err = compoundError(err, errmsg)
				return eventList, err
			}
			// ISSUE filter non-suported types? 
			evnt := &Event{currentTime, sys}
			acc = append(acc, evnt)
		case isMetaStatus(status):
			var meta *MetaMessage
			useRunningStatus = false
			meta, index, err = getMetaMessage(bytes, startIndex)
			if err != nil {
				errmsg := "smf.CreateEventList error in switch case isMetaStatus"
				err = compoundError(err, errmsg)
				return eventList, err
			}
			// ISSUE fitler non-suported types?
			// ISSUE update tempo
			evnt := &Event{currentTime, meta}
			acc = append(acc, evnt)
		default:
			errmsg := "smf.CreateEventList expected staus byte at index %d, got 0x%x"
			err = exError(fmt.Sprintf(errmsg, index, status))
			return eventList, err
		}
				
	} // end for index
	eventList.events = acc
	return eventList, err
} // end func

