package midi

import (
	"testing"
	"math/rand"
	gomidi "gitlab.com/gomidi/midi/v2"
)



func onEvent(ci byte, key byte) gomidi.Message {
	st := byte(0x90 | ci)
	return gomidi.NewMessage([]byte{st, key, 64})
}


func offEvent(ci byte, key byte) gomidi.Message {
	var st byte
	if rand.Intn(100) > 50 {
		st = 0x80 | ci
	} else {
		st = 0x90 | ci
	}
	return gomidi.NewMessage([]byte{st, key, 0})
}

func TestNoteQueue(t *testing.T) {
	nq := MakeNoteQueue()
	for ci := byte(0); ci < 15; ci++ {
		for key := byte(0); key < 8; key++ {
			for n := byte(0); n < ci+1; n++ {
				nq.Update(onEvent(ci, key))
			}
			for n := byte(0); n < ci; n++ {
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
