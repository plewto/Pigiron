package midi

/*
** Defines scheme for tracking and resolving open MIDI notes.
**
*/

import (
	gomidi "gitlab.com/gomidi/midi/v2"
)


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

func (nqc *noteQueueChannel) updateCount(msg gomidi.Message) {
	switch {
	case IsNoteOff(msg):
		key := msg.Data[1]
		count := nqc.openCounts[key] -1
		if count < 0 {
			count = 0
		}
		nqc.openCounts[key] = count
	case IsNoteOn(msg):
		key := msg.Data[1]
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
func (nq *NoteQueue) Update(msg gomidi.Message) {
	d := msg.Data
	if len(d) > 0 {
		st := d[0] & 0xF0
		if st == 0x80 || st == 0x90 {
			ci := d[0] & 0x0F
			nq.channels[ci].updateCount(msg)
		}
	}
}

// nq.OpenCount() returns number of unresolved note-on events for given channel and key.
// The channel parameter has interval 0 <= ci <= 15.
//
func (nq *NoteQueue) OpenCount(ci byte, key byte) int {
	if ci < 0 || 16 < ci || key < 0 || 128 < key {
		return 0
	}
	nqc := nq.channels[ci] 
	return nqc.openCounts[key]
}

// nq.OffEvents() returns list of MIDI off events required to resolve all open notes.
//
func (nq *NoteQueue) OffEvents() []gomidi.Message {
	var acc = make([]gomidi.Message, 0, 48)
	for ci := byte(0); ci < byte(len(nq.channels)); ci++ {
		st := byte(0x80 | ci)
		for key := byte(0); key < 128; key++ {
			for n := 0; n < nq.OpenCount(ci, key); n++ {
				off := gomidi.NewMessage([]byte{st, key, 0})
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
