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





func convertTrackBytes(bytes []byte) (track *SMFTrack, err error) {
	var events = make([]*UniversalEvent, 0, 1024)
	var runningStatus StatusByte = StatusByte(0)
	var runningChannel MIDIChannelNibble
	var index int = 0 // FOR error reports only.
	track = &SMFTrack{make([]*UniversalEvent, 0, 512)}
	for len(bytes) > 0 {
		var vlq *VLQ
		vlq, bytes, err = takeVLQ(bytes)
		if err != nil {
			errmsg := "convertTrackBytes takeVLQ() error  at index %d"
			err = pigerr.CompoundError(err, fmt.Sprintf(errmsg, index))
			return
		}
		index += vlq.Length()
		deltaTime := vlq.Value()
		var status byte
		status, bytes, err = takeByte(bytes)
		if err != nil {
			errmsg := "convertTrackBytes expected status byte "
			errmsg += "(or running status) at index %d"
			err = pigerr.CompoundError(err, fmt.Sprintf(errmsg, index))
			return
		}
		index++
		switch {
		case status & 0x80 == 0: // using running status
			dataCount, _ := channelStatusDataCount[StatusByte(runningStatus)]
			if dataCount == 0 {
				errmsg1 := "convertTrackBytes running status dataCount = 0"
				errmsg2 := "runningStatus = 0x%02X %s, index = %d"
				err = pigerr.New(errmsg1, fmt.Sprintf(errmsg2, byte(runningStatus), runningStatus, index))
				return
			}
			d1, d2 := byte(0), byte(0)
			d1, bytes, err = takeByte(bytes)
			if err != nil {
				errmsg := "convertTrackBytes get running status data byte 1 error, index = %d"
				err = pigerr.CompoundError(err, fmt.Sprintf(errmsg, index))
				return
			}
			index++
			if dataCount == 2 {
				d2, bytes, err = takeByte(bytes)
				if err != nil {
					errmsg := "convertTrackBytes get running status data byte 2 error, index = %d"
					err = pigerr.CompoundError(err, fmt.Sprintf(errmsg, index))
					return
				}
				index++
			}
			var ue *UniversalEvent
			ue, err = MakeChannelEvent(runningStatus, runningChannel, d1, d2)
			if err != nil {
				errmsg1 := "convertTrackByets error while creating running-status message "
				errmsg2 := "runningStatus = 0x%02X  runningChannel = %d,  index = %d"
				err = pigerr.CompoundError(err, errmsg1, fmt.Sprintf(errmsg2, byte(runningStatus), byte(runningChannel), index))
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
				errmsg1 := "convertTrackBytes error while creating Channel Message"
				errmsg2 := "status = 0x%02X, channel = %d,  index = %d"
				err = pigerr.CompoundError(err, errmsg1, fmt.Sprintf(errmsg2, byte(runningStatus), byte(runningChannel), index))
				return
			}
			index++
			if dataCount == 2 {
				d2, bytes, err = takeByte(bytes)
				if err != nil {
					errmsg1 := "convertTrackBytes error while reading 2nd Channel Message data byte."
					errmsg2 := "status = 0x%02X, channel = %d,  index = %d"
					err = pigerr.CompoundError(err, errmsg1, fmt.Sprintf(errmsg2, byte(runningStatus), byte(runningChannel), index))
					return
				}
				index++
			}
			var ue *UniversalEvent
			ue, err = MakeChannelEvent(runningStatus, runningChannel, d1, d2)
			if err != nil {
				errmsg1 := "convertTrackBytes error while crearting Chanbnel Message"
				errmsg2 := "status = 0x%02X, channel = %d,  index = %d"
				err = pigerr.CompoundError(err, errmsg1, fmt.Sprintf(errmsg2, byte(runningStatus), byte(runningChannel), index))
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
					errmsg1 := "convertTrackBytes error while creating System Message"
					errmsg2 := "status = 0x%02X,  index = %d"
					err = pigerr.CompoundError(err, errmsg1, fmt.Sprintf(errmsg2, byte(status), index))
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
				errmsg1 := "convertTrackBytes error getting meta-type byte, index = %d"
				err = pigerr.CompoundError(err, fmt.Sprintf(errmsg1, index))
				return
			}
			if !isMetaType(mtype) {
				errmsg := "convertTrackBytes Expected meta type at index %d, got 0x%2X"
				err = pigerr.New(fmt.Sprintf(errmsg, index, mtype))
				return
			}
			index++
			var countVLQ *VLQ
			countVLQ, bytes, err = takeVLQ(bytes)
			if err != nil {
				errmsg1 := "convertTrackBytes error geting meta event count, index = %d"
				err = pigerr.CompoundError(err, fmt.Sprintf(errmsg1, index))
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
					errmsg1 := "convertTrackBytes error making Meta Message"
					errmsg2 := "meta type = 0x%02X '%s',  index = %d"
					err = pigerr.CompoundError(err, errmsg1, fmt.Sprintf(errmsg2, byte(mtype), mtype, index))
					return
				}
				index += count
				ue.deltaTime = deltaTime
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
