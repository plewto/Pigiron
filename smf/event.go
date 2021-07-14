package smf

import (
	"fmt"
	"github.com/rakyll/portmidi"
)

type Event struct {
	time float64
	message *portmidi.Event  // only one field message or meta should be non-nil
	meta *MetaMessage
}

func portmidiSysExToString(pm *portmidi.Event) string {
	acc := "SysEx "
	maxbytes := 12
	data := pm.SysEx
	if len(data) == 0 {
		acc += "<malformed, no data bytes>"
		return acc
	}
	acc += "["
	for i := 0; i < len(data) && i < maxbytes; i++ {
		acc += fmt.Sprintf("%2X ", data[i])
	}
	if len(data) > maxbytes {
		acc += fmt.Sprintf("... %d more]", len(data) - maxbytes)
	} else {
		acc += "]"
	}
	return acc
}	

func portmidiToString(pm *portmidi.Event) string {
	acc := ""
	st := byte(pm.Status)
	mn, _ := statusTable[StatusByte(st)]
	switch {
	case len(pm.SysEx) > 0:
		acc += portmidiSysExToString(pm)
	case isSystemStatus(st):
		acc += fmt.Sprintf("SYS  %s", mn)
	case isChannelStatus(st):
		ch := (st & 0x0F) + 1
		d1, d2 := pm.Data1, pm.Data2
		acc += fmt.Sprintf("CHAN+ %2d %s %02X", ch, mn, d1)
		count, _ := channelStatusDataCount[StatusByte(st)]
		if count > 1 {
			acc += fmt.Sprintf(" %02X", d2)
		}
		if isKeyedStatus(st) {
			keys := []string{"C", "C#", "D", "D#", "E", "F", "F#", "G", "G#", "A", "A#", "B"}
			pc := d1 % 12
			oct := d1/12
			acc += fmt.Sprintf("   key %s%d", keys[pc], oct)
		}
		default:
		acc += "<malformed, invalid status type>"
	}
	return acc
}




func (e *Event) String() string {
	acc := fmt.Sprintf("t %8.4f ", e.time)
	if e.message != nil {
		acc += portmidiToString(e.message)
	} else {
		acc += fmt.Sprintf("%s", e.meta)
	}
	return acc
}


type EventList struct {
	tempo float64   // BPM
	division int    // clocks per quarter note
	events []*Event
}

func (e *EventList) tickDuration() float64 { // ISSUE: Implement
	return 0.01
}
	

func (e *EventList) getDeltaTime(vlq *VLQ) float64 {
	tick := e.tickDuration()
	return tick * float64(vlq.Value())
}

func (events *EventList) Dump() {
	tab := " "
	fmt.Println("EventList")
	fmt.Printf("%s tempo    : %f BPM\n", tab, events.tempo)
	fmt.Printf("%s division : %d\n", tab, events.division)
	fmt.Printf("%s event count : %d\n", tab, len(events.events))
	for i, ev := range events.events {
		fmt.Printf("%s [#%4d] %s\n", tab, i, ev)
	}
}



func createEventList(division int, bytes []byte) (*EventList, error) {
	var err error
	var eventList *EventList
	var acc []*Event = make([]*Event, 0, 1024)
	var runStat StatusByte = 0
	var runChan byte = 0
	var currentTime float64 = 0
	var index int
	for index = 0; index < len(bytes); {
		startIndex := index
		var vlq *VLQ
		vlq, err = getVLQ(bytes, index)
		if err != nil {
			errmsg := "smf createEventList expected VLQ at index %d"
			err = compoundError(err, fmt.Sprintf(errmsg, index))
			return eventList, err
		}
		//currentTime += ????   // TODO: update currentTime
		index += vlq.Length()
		if index >= len(bytes) {
			break
		}
		status := bytes[index]
		switch {
		case status & 0x80 == 0:  // Use running status
			fmt.Printf("DEBUG running status at %d\n", startIndex)
			var cmsg *ChannelMessage
			if !isChannelStatus(byte(runStat)) {
				errmsg := "smf createEventList running-status error. startIndex was %d\n"
				errmsg += "Expected runStat to be a ChannelStatus, go 0x%x"
				err = exError(fmt.Sprintf(errmsg, startIndex, runStat))
				return eventList, err
			}
			cmsg, index, err = getRunningStatus(bytes, index, runStat, runChan)
			if err != nil {
				errmsg := "smf createEventList running-status error. startIndex was %d\n"
				err = compoundError(err, fmt.Sprintf(errmsg, startIndex))
				return eventList, err
			}
			pmidi, _ := cmsg.ToPortmidiEvent()
			evnt := &Event{currentTime, &pmidi, nil}
			acc = append(acc, evnt)
		case isChannelStatus(status & 0xF0):
			var cmsg *ChannelMessage
			runStat = StatusByte(byte(status) & 0xF0)
			runChan = status & 0x0F
			cmsg, index, err = getChannelMessage(bytes, index)
			if err != nil {
				errmsg := "smf createEventList error in switch case isChannelMessage"
				err = compoundError(err, errmsg)
				return eventList, err
			}
			pmidi, _ := cmsg.ToPortmidiEvent()
			evnt := &Event{currentTime, &pmidi, nil}
			acc = append(acc, evnt)
		case isSystemStatus(status):
			var sys *SystemMessage
			runStat = 0
			runChan = 0
			sys, index, err = getSystemMessage(bytes, index)
			if err != nil {
				errmsg := "smf createEventList error in switch case isSystemStatus"
				errmsg += "startIndex was %d"
				err = compoundError(err, fmt.Sprintf(errmsg, startIndex))
				return eventList, err
			}
			pmidi, _ := sys.ToPortmidiEvent()
			evnt := &Event{currentTime, &pmidi, nil}
			acc = append(acc, evnt)
		case isMetaStatus(status):
			var meta *MetaMessage
			runStat = 0
			runChan = 0
			meta, index, err = getMetaMessage(bytes, index)
			if err != nil {
				errmsg := "smf createEventList error in switch case isMetaStatus"
				errmsg += "startIndex was %d"
				err = compoundError(err, fmt.Sprintf(errmsg, startIndex))
				return eventList, err
			}
			evnt := &Event{currentTime, nil, meta}
			acc = append(acc, evnt)
		default:
			errmsg := "smf createEventList expected status byte at index %d, got 0x%x"
			err = exError(fmt.Sprintf(errmsg, index, status))
			return eventList, err	
		} // end switch	
	} // end outer for
	eventList = &EventList{60, division, acc}  // TODO replace 60 with actual tempo
	return eventList, err
} // end createEventList




