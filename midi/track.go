package midi


import (
	"fmt"
	"github.com/plewto/pigiron/pigerr"
)

var (
	trackID chunkID = [4]byte{0x4d, 0x54, 0x72, 0x6B}
)

type SMFTrack struct {
	events []*UniversalEvent
}

func (trk *SMFTrack) ID() chunkID {
	return trackID
}

func (trk *SMFTrack) Length() int {
	return len(trk.events)
}

func (trk *SMFTrack) Events() []*UniversalEvent {
	return trk.events
}

func (trk *SMFTrack) Dump() {
	fmt.Println("SMFTrack")
	fmt.Printf("  ID  : %s\n", trk.ID())
	fmt.Printf("  LEN : %d\n", trk.Length())
	for i, ev := range trk.Events() {
		fmt.Printf("  [%4d] %s\n", i, ev)
	}
}

func convertTrackBytes(bytes []byte) (track *SMFTrack, err error) {
	var events = make([]*UniversalEvent, 0, 1024)
	var runningStatus StatusByte = StatusByte(0)
	var runningChannel MIDIChannelNibble
	var index int = 0 // FOR error reports only.
	track = &SMFTrack{make([]*UniversalEvent, 0, 512)}

	// creates simple error 
	var error1 = func(text string) error {
		msg1 := fmt.Sprintf("convertTrackBytes() at index %d", index)
		return pigerr.New(msg1, text)
	}

	// creates compound error
	var error2 = func(err error, text string) error {
		msg1 := fmt.Sprintf("convertTrackBytes() at index %d", index)
		return pigerr.CompoundError(err, msg1, text)
	}
	
	for len(bytes) > 0 {
		var vlq *VLQ
		vlq, bytes, err = takeVLQ(bytes)
		if err != nil {
			errmsg := "Invalid delta-time"
			err = error2(err, errmsg)
			return
		}
		index += vlq.Length()
		deltaTime := vlq.Value()
		var status byte
		status, bytes, err = takeByte(bytes)
		if err != nil {
			errmsg := "Expected status byte or running-status data byte"
			err = error2(err, errmsg)
			return
		}
		index++
		switch {
		case status & 0x80 == 0: // using running status
			dataCount, _ := channelStatusDataCount[StatusByte(runningStatus)]
			if dataCount == 0 {
				errmsg := "Illegal runningStatus 0x%02X"
				err = error1(fmt.Sprintf(errmsg, byte(runningStatus)))
			}
			d1, d2 := byte(0), byte(0)
			d1, bytes, err = takeByte(bytes)
			if err != nil {
				errmsg := "error geting running-status first data byte"
				err = error2(err, errmsg)
				return
			}
			index++
			if dataCount == 2 {
				d2, bytes, err = takeByte(bytes)
				if err != nil {
					errmsg := "error getting running-status second data byte."
					err = error2(err, errmsg)
					return
				}
				index++
			}
			var ue *UniversalEvent
			ue, err = MakeChannelEvent(runningStatus, runningChannel, d1, d2)
			if err != nil {
				errmsg := "creating running-status message, runningStatus = 0x%02X, runningChannel = %d"
				err = error2(err, fmt.Sprintf(errmsg, byte(runningStatus), byte(runningChannel)))
				return 
			}
			ue.deltaTime = deltaTime
			events = append(events, ue)
		case isChannelStatus(status):
			runningStatus = StatusByte(status & 0xF0)
			runningChannel = MIDIChannelNibble(status & 0x0F)
			dataCount, _ := channelStatusDataCount[runningStatus]
			d1, d2 := byte(0), byte(0)
			d1, bytes, err = takeByte(bytes)
			if err != nil {
				errmsg := "creating channel message, status = 0x%02X, channel = %d"
				err = error2(err, fmt.Sprintf(errmsg, byte(runningStatus), byte(runningChannel)))
				return
			}
			index++
			if dataCount == 2 {
				d2, bytes, err = takeByte(bytes)
				if err != nil {
					errmsg := "getting channel message second data byte, status = 0x%02X, channel = %d"
					err = error2(err, fmt.Sprintf(errmsg, byte(runningStatus), byte(runningChannel)))
					return
				}
				index++
			}
			var ue *UniversalEvent
			ue, err = MakeChannelEvent(runningStatus, runningChannel, d1, d2)
			if err != nil {
				errmsg := "creating channel message, status = 0x%02X, channel = %d"
				err = error2(err, fmt.Sprintf(errmsg, byte(runningStatus), byte(runningChannel)))
				return
			}
			ue.deltaTime = deltaTime
			events = append(events, ue)
		case isSystemStatus(status):
			var ue *UniversalEvent
			runningStatus = StatusByte(0)
			if StatusByte(status) == SYSEX {
				// TODO build sysex message
			} else {
				ue, err = MakeSystemEvent(StatusByte(status))
				if err != nil {
					errmsg := "creatring sysmte-message, status = 0x%02X"
					err = error2(err, fmt.Sprintf(errmsg, byte(status)))
					return
				}
				index++
			}
			ue.deltaTime = deltaTime
			events = append(events, ue)
		case isMetaStatus(status):
			var mtype byte
			var ue *UniversalEvent
			mtype, bytes, err = takeByte(bytes)
			if err != nil {
				errmsg := "creating meta-type"
				err = error2(err, errmsg)
				return
			}
			if !isMetaType(mtype) {
				errmsg := "invalid metaType 0x%02X"
				err = error1(fmt.Sprintf(errmsg, mtype))
				return
			}
			index++
			var countVLQ *VLQ
			countVLQ, bytes, err = takeVLQ(bytes)
			if err != nil {
				errmsg1 := "invalid meta event count"
				err = error2(err, errmsg1)
				return
			}
			index += countVLQ.Length()
			var count = countVLQ.Value()
			switch byte(mtype) {
			case byte(META_END_OF_TRACK):
				ue, _ = MakeMetaEvent(META_END_OF_TRACK, []byte{})
				bytes = []byte{}
				ue.deltaTime = deltaTime
				events = append(events, ue)
				break
			default:
				data := bytes[0:count]
				bytes = bytes[count:]
				ue, err = MakeMetaEvent(MetaType(mtype), data)
				if err != nil {
					errmsg := "creating meta message, meta type = 0x%02X"
					err = error2(err, fmt.Sprintf(errmsg, byte(mtype)))
					return
				}
				index += count
				ue.deltaTime = deltaTime
				events = append(events, ue)
			}
		default:
			errmsg := "Unhandled switch case, status = 0x%02X"
			err = error1(fmt.Sprintf(errmsg, status))
			return
		}
	}
	track.events = events
	return
}
