package midi

import (
	"testing"
	"math/rand"
	"github.com/rakyll/portmidi"
)

func onEvent(ci int, key int) *portmidi.Event {
	st := int64(0x90 | ci)
	return &portmidi.Event{0, st, int64(key), int64(64), []byte{}}
}

func offEvent(ci int, key int) *portmidi.Event {
	var st int64
	if rand.Intn(100) > 50 {  // randomly use running status
		st = int64(0x80 | ci)
	} else {
		st = int64(0x90 | ci)
	}
	return &portmidi.Event{0, st, int64(key), int64(0), []byte{}}
}

func TestNoteQueue(t *testing.T) {
	nq := MakeNoteQueue()
	for ci := 0; ci < 15; ci++ {
		for key := 0; key < 8; key++ {
			for n := 0; n < ci+1; n++ {
				nq.Update(onEvent(ci, key))
			}
			for n := 0; n < ci; n++ {
				nq.Update(offEvent(ci, key))
			}
			diff := nq.OpenCount(ci, key)
			if diff != 1 {
				msg := "Expected note count 1, ci = %d, key = %d, got count %d"
				t.Fatalf(msg, ci, key, diff)
			}
		}
	}
	// check floor value, open count should never be less then 0.
	for i := 0; i < 100; i++ {
		nq.Update(offEvent(0, 0))
	}
	if nq.OpenCount(0, 0) != 0 {
		msg := "Negative open note count: %d"
		t.Fatalf(msg, nq.OpenCount(0, 0))
	}
}
