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

var TRACK_ID chunkID = [4]byte{0x4d, 0x54, 0x72, 0x6B}

type Track struct {
	events []Event
}

func (trk *Track) ID() chunkID {
	return TRACK_ID
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

func (trk *Track) Dump() {
	fmt.Println("Track:")
	for i, ev := range trk.events {
		fmt.Printf("[%5d] %s\n", i, ev.String())
	}
}
		
func ReadTrack(f *os.File) (track *Track, err error) {
	var id chunkID
	var length int
	id, length, err = readChunkPreamble(f)
	if err != nil {
		return
	}
	if !id.eq(TRACK_ID) {
		msg := "Expected track id '%s', got '%s'"
		err = fmt.Errorf(msg, TRACK_ID, id)
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
		msg := "smf.ReadTrack could not read track\n"
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

func (trk *Track) convertEvents(bytes []byte) (index int, err error) {
	var acc = make([]Event, 0, 1024)
	var runningStatus = midi.StatusByte(0)
	index = 0
	for index < len(bytes) {
		var vlq *VLQ
		vlq, index, err = expectVLQ(bytes, index)
		if err != nil {
			return
		}
		var deltaTime = vlq.Value()
		var b byte
		b, index, err = expectByte(bytes, index)
		status := midi.StatusByte(b)
		if err != nil {
			return
		}
		var msgdata []byte
		if !midi.IsChannelStatus(status) { 
			if !midi.IsChannelStatus(runningStatus) {
				errmsg := "Expected running status at index %d"
				err = fmt.Errorf(errmsg, index-1)
				return
			}
			msgdata, index, err = expectChannelMessage(bytes, byte(runningStatus), index-1)
			if err != nil {
				return
			}
		} else {
			switch {
			case midi.IsChannelStatus(status):
				msgdata, index, err = expectChannelMessage(bytes, byte(status), index)
				if err != nil {
					return
				}
				runningStatus = status
			case status == midi.SYSEX:
				msgdata, index, err = expectSysexMessage(bytes, index)
				if err != nil {
					return
				}
				runningStatus = midi.StatusByte(0)
			case midi.IsSystemStatus(status):
				msgdata, index, err = expectSystemMessage(bytes, index)
				if err != nil {
					return
				}
				runningStatus = midi.StatusByte(0)
			case midi.IsMetaStatus(status):
				msgdata, index, err = expectMetaMessage(bytes, index)
				if err != nil {
					return
				}
			default:
				errmsg := "Unhandled switch case. index = %d, status = 0x%02X"
				err = fmt.Errorf(errmsg, index, status)
				return
			}
			
		}
		acc = append(acc, Event{uint64(deltaTime), gomidi.NewMessage(msgdata)})
	} 
	events := make([]Event, len(acc), len(acc))
	for i, e := range acc {
		events[i] = e
	}
	trk.events = events
	return
	
}
	
	
