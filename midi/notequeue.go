package midi

/*
** Defines scheme for tracking and resolving open MIDI notes.
**
*/

import "github.com/rakyll/portmidi"


// noteQueueChannel counts unresolved note-on events for a MIDI channel.
//
type noteQueueChannel struct {
	openCounts [128]int
}

func (nqc *noteQueueChannel) reset() {
	for i := 0; i < len(nqc.openCounts); i++ {
		nqc.openCounts[i] = 0
	}
}

// isNoteOff() returns true iff event status is 0x80 or is 0x90 with velocity 0.
//
func isNoteOff(event *portmidi.Event) bool {
	st := event.Status & 0xF0
	velocity := event.Data2
	return st == 0x80 || (st == 0x90 && velocity == 0)
}


// isNoteOn() returns true iff event status is 0x90 and velocity > 0.
//
func isNoteOn(event *portmidi.Event) bool {
	st := event.Status & 0xF0
	velocity := event.Data2
	return st == 0x90 && velocity > 0
}

// nqc.updateCount() increments/decrements specific open note count.
//
func (nqc *noteQueueChannel) updateCount(event *portmidi.Event) {
	switch {
	case isNoteOff(event):
		key := event.Data1
		count := nqc.openCounts[key] -1
		if count < 0 {
			count = 0
		}
		nqc.openCounts[key] = count
	case isNoteOn(event):
		key := event.Data1
		nqc.openCounts[key]++
	default:
		// ignore
	}
}

/*
** NoteQueue struct maintains a count of all unresolved note-on events.
**
*/
type NoteQueue struct {
	channels [16]noteQueueChannel
}

// nq.Reset() sets all note-counts to 0.
//
func (nq *NoteQueue) Reset() {
	for _, c := range nq.channels {
		c.reset()
	}
}


// nq.Update() increments/decrements count for specific note and channel.
// The note count is never negative.
//
func (nq *NoteQueue) Update(event *portmidi.Event) {
	st := event.Status & 0xF0
	if st == 0x80 || st == 0x90 {
		ci := event.Status & 0x0F
		nq.channels[ci].updateCount(event)
	}
}

// nq.OpenCount() returns number of unresolved note-on events for given channel and key.
// The channel parameter has interval 0 <= ci <= 15.
//
func (nq *NoteQueue) OpenCount(ci int, key int) int {
	if ci < 0 || 16 < ci || key < 0 || 128 < key {
		return 0
	}
	nqc := nq.channels[ci] 
	return nqc.openCounts[key]
}

// nq.OffEvents() returns list of MIDI off events required to resolve all open notes.
//
func (nq *NoteQueue) OffEvents() []*portmidi.Event {
	var acc = make([]*portmidi.Event, 0, 48)
	for ci := 0; ci < len(nq.channels); ci++ {
		st := int64(0x80 | ci)
		for key := int64(0); key < 128; key++ {
			for n := 0; n < nq.OpenCount(ci, int(key)); n++ {
				off := &portmidi.Event{0, st, key, 0, []byte{}}
				acc = append(acc, off)
			}
		}
	}
	return acc
}


// MakeNoteQueue() creates new instance of NoteQueue struct.
//
func MakeNoteQueue() *NoteQueue {
	clist := [16]noteQueueChannel{}
	for ci, _ := range clist {
		clist[ci] = noteQueueChannel{[128]int{}}
	}
	return &NoteQueue{clist}
}
