package smf

/*
** track.go defines SMF track chunks.
**
*/

import (
 	"fmt"
	"os"
	"github.com/plewto/pigiron/midi"
 	gomidi "gitlab.com/gomidi/midi/v2"
)

var trackID chunkID = [4]byte{0x4d, 0x54, 0x72, 0x6B}

type Track struct {
	events []Event
}

func (trk *Track) ID() chunkID {
	return trackID
}

func (trk *Track) Length() int {
	acc := 0
	for _, ev := range trk.events {
		acc += ev.Length()
	}
	return acc
}

func (trk *Track) String() string {
	return fmt.Sprintf("Track (%d events)", len(trk.events))
}

func (trk *Track) Events() []Event {
	return trk.events
}


func readTrack(f *os.File) (track *Track, err error) {
	var id chunkID
	var length int
	id, length, err = readChunkPreamble(f)
	if err != nil {
		return
	}
	if !id.eq(trackID) {
		msg := "Expected track id '%s', got '%s'"
		err = fmt.Errorf(msg, trackID, id)
		return
	}
	var bytes = make([]byte, length)
	var readCount = 0
	readCount, err = f.Read(bytes)
	if readCount != length {
		msg := "Expected %d track bytes, read %d"
		err = fmt.Errorf(msg, length, readCount)
		return
	}
	if err != nil {
		msg := "smf.readTrack could not read track\n"
		msg += fmt.Sprintf("%s", err)
		err = fmt.Errorf(msg)
		return
	}
	track = new(Track)
	var index int
	index, err = track.convertEvents(bytes)
	if err != nil {
		errmsg := fmt.Sprintf("Error while converting smf track bytes to events.  index = %d", index)
		err = fmt.Errorf("%s\n%s", errmsg, err)
	}
	return
}



func (trk *Track) convertEvents(buffer []byte) (index int, err error) {
	var acc = make([]Event, 0, 1024)
	var runningStatus = midi.StatusByte(0)
	index = 0
	for index < len(buffer) {
		var vlq *VLQ
		vlq, index, err = ExpectVLQ(buffer, index)
		if err != nil {
			return
		}
		var deltaTime = uint64(vlq.Value())
		var b = buffer[index]
		var msgBytes []byte
		if b > 0x7F {   // new statys byte
			var st = midi.StatusByte(b)
			switch {
			case midi.IsChannelStatus(st):
				runningStatus = st
				msgBytes, index, err = ExpectChannelMessage(buffer, b, index)
			case b == byte(midi.SYSEX):
				runningStatus = midi.StatusByte(0)
				msgBytes, index, err = ExpectSysexMessage(buffer, index)
			case midi.IsSystemRealtimeStatus(st):
				runningStatus = midi.StatusByte(0)
				msgBytes, index, err = ExpectSystemMessage(buffer, index)
			case midi.IsMetaStatus(st):
				runningStatus = midi.StatusByte(0)
				msgBytes, index, err = ExpectMetaMessage(buffer, index)
				if err == nil && msgBytes[1] == byte(midi.META_END_OF_TRACK) {
					break
				}
			case runningStatus != 0:
				msgBytes, index, err = ExpectRunningStatus(buffer, byte(runningStatus), index)
			default:
				errmsg := "smf.Track.convertEvents switch default.\n"
				errmsg += "This should never happen, buffer index was %d"
				err = fmt.Errorf(errmsg, index)
			}
		} else { // assume running status
			if runningStatus == 0 {
				errmsg := "Expected running status at index %d"
				err = fmt.Errorf(errmsg, index)
				return
			}
			msgBytes, index, err = ExpectRunningStatus(buffer, byte(runningStatus), index)
		}
		if err != nil {
			return
		}
		acc = append(acc, Event{deltaTime, gomidi.NewMessage(msgBytes)})
	}
	events := make([]Event, len(acc), len(acc))
	for i, e := range acc {
		events[i] = e
	}
	trk.events = events
	return
}

func (trk *Track) Dump() string {
	var acc = fmt.Sprintf("Track  %d events\n", len(trk.events))
	var time = uint64(0)
	for i, evnt := range trk.events {
		acc += fmt.Sprintf("[%4d] t %8d %s\n", i, time, evnt.String())
		time += evnt.deltaTime
	}
	return acc
}
