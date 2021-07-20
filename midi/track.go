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

func tickDuration(division int, tempo float64) float64 {
	division = division & 0x7FFF
	if tempo == 0 {
		dflt := 60.0
		errmsg := "MIDI tempo is 0, using default %f"
		pigerr.NewWarning(fmt.Sprintf(errmsg, dflt))
		tempo = dflt
	}
	var qdur float64 = 60.0/tempo
	return qdur/float64(division)
}


func debugHexDump(index int, bytes []byte) {
	max := 8
	fmt.Printf("bytes [%4d]: ", index)
	j := max
	for i := 0; i < len(bytes); i = i+1 {
		j--
		if j == 0 {
			break
		}
		fmt.Printf("%02X ", bytes[i])
	}
	fmt.Printf("   length = %d\n", len(bytes))
}


func convertTrackBytes(division int, tempo float64, bytes []byte) (track *SMFTrack, err error) {
	var events = make([]*UniversalEvent, 0, 1024)
	var runningStatus StatusByte = StatusByte(0)
	var runningChannel MIDIChannelNibble
	var currentTime float64 = 0
	var index int = 0 // FOR error report only
	var tkDuration = tickDuration(division, tempo)
	track = &SMFTrack{make([]*UniversalEvent, 0, 512)}
	for len(bytes) > 0 {
		// debugHexDump(index, bytes)
		var vlq *VLQ
		vlq, bytes, err = takeVLQ(bytes)
		if err != nil {
			// TODO error message
			return
		}
		index += vlq.Length()
		deltaTime := vlq.Value()
		currentTime += float64(deltaTime) * tkDuration
		fmt.Printf("DEBUG delta %d, time %f\n", deltaTime, currentTime)
		var status byte
		status, bytes, err = takeByte(bytes)
		if err != nil {
			// TODO error message
			return
		}
		index++
		switch {
		case status & 0x80 == 0: // using running status
			dataCount, _ := channelStatusDataCount[StatusByte(runningStatus)]
			if dataCount == 0 {
				// TODO include index in error message
				errmsg := "running status  dataCount = 0"
				err = pigerr.New(errmsg)
				return
			}
			d1, d2 := byte(0), byte(0)
			d1, bytes, err = takeByte(bytes)
			if err != nil {
				// TODO error message
				return
			}
			index++
			if dataCount == 2 {
				d2, bytes, err = takeByte(bytes)
				if err != nil {
					// TODO error message
					return
				}
				index++
			}
			var ue *UniversalEvent
			ue, err = MakeChannelEvent(runningStatus, runningChannel, d1, d2)
			if err != nil {
				// TODO error message
				return
			}
			ue.time = currentTime
			events = append(events, ue)
		case isChannelStatus(status):
			runningStatus = StatusByte(status & 0xF0)
			runningChannel = MIDIChannelNibble(status & 0x0F)
			dataCount, _ := channelStatusDataCount[runningStatus]
			d1, d2 := byte(0), byte(0)
			d1, bytes, err = takeByte(bytes)
			if err != nil {
				// TODO error message
				return
			}
			index++
			if dataCount == 2 {
				d2, bytes, err = takeByte(bytes)
				if err != nil {
					// TODO error message
					return
				}
				index++
			}
			var ue *UniversalEvent
			ue, err = MakeChannelEvent(runningStatus, runningChannel, d1, d2)
			if err != nil {
				// TODO error message
				return
			}
			ue.time = currentTime
			events = append(events, ue)
		case isSystemStatus(status):
			var ue *UniversalEvent
			runningStatus = StatusByte(0)
			if StatusByte(status) == SYSEX {
				// TODO build sysex message
			} else {
				// system realtime message
				
				ue, err = MakeSystemEvent(StatusByte(status))
				if err != nil {
					// TODO error message
					return
				}
				index++
			}
			ue.time = currentTime
			events = append(events, ue)
		case isMetaStatus(status):
			var mtype byte
			var ue *UniversalEvent
			mtype, bytes, err = takeByte(bytes)
			if err != nil {
				// TODO error message
				return
			}
			if !isMetaType(mtype) {
				errmsg := "Expected meta type at index %d, got 0x%2X"
				err = pigerr.New(fmt.Sprintf(errmsg, index, mtype))
				return
			}
			index++
			var countVLQ *VLQ
			countVLQ, bytes, err = takeVLQ(bytes)
			if err != nil {
				// TODO error message
				return
			}
			index += countVLQ.Length()
			var count = countVLQ.Value()
			switch byte(mtype) {
			case byte(META_END_OF_TRACK):
				ue, _ = MakeMetaEvent(META_END_OF_TRACK, []byte{})
				bytes = []byte{}
				ue.time = currentTime
				events = append(events, ue)
				break
			default:
				data := bytes[0:count]
				bytes = bytes[count:]
				ue, err = MakeMetaEvent(MetaType(mtype), data)
				if err != nil {
					// TODO error message
					return
				}
				// TODO update tkDuration for tempo events
				index += count
				ue.time = currentTime
				events = append(events, ue)
			}
		default:
			errmsg := "Unhandle switch case for convertTrackBytes "
			errmsg += "index = %d, status = 0x%2X"
			err = pigerr.New(fmt.Sprintf(errmsg, index, status))
			return
		}
	}
	track.events = events
	return
}
