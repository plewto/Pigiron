package midi

import (
	"fmt"
	gomidi "gitlab.com/gomidi/midi/v2"
)

const (
	TEMPO_CONSTANT float64 = 60000000
	MAX_TEMPO float64 = 300
)


// MakeTempoMessage creates new meta tempo-change message.
// tempo argument in BPM.
// Returns non-nil error if tempo is out of bounds.
//
func MakeTempoMessage(tempo float64) (msg gomidi.Message, err error) {
	var usec uint64
	var d0, d1, d2 byte
	if tempo <= 0 || tempo > MAX_TEMPO {
		errmsg := "Tempo out of bounds: %f"
		err = fmt.Errorf(errmsg, tempo)
		return
	}
	usec = uint64(TEMPO_CONSTANT/tempo)
	d0 = byte((usec & 0xFF0000) >> 16)
	d1 = byte((usec & 0x00FF00) >> 8)
	d2 = byte(usec & 0x0000FF)
	msg = gomidi.NewMessage([]byte{0xFF, 0x51, 0x03, d0, d1, d2})
	return
}

// MetaTempoMicroseconds return micro seconds per quarter not for meta temp message.
// The error return is non-nil if message is not a meta tempo-change.
//
func MetaTempoMicroseconds(msg gomidi.Message) (usec uint64, err error) {
	d := msg.Data
	if len(d) != 6 || d[0] != 0xFF || d[1] != 0x51 {
		errmsg := "%v is not a meta tempo message"
		err = fmt.Errorf(errmsg, msg)
		return  usec, err
	}
	 usec = 0
	for i, shift := 3, 16; i < 6; i, shift = i+1, shift-8 {
		 usec += uint64(d[i]) << shift
	}
	if usec == 0 {
		errmsg := "%v is malformed meta tempo message"
		err = fmt.Errorf(errmsg, d)
	}
	return  usec, err
}
	
// MetaTempBPM returns tempo in BPM for meta tempo-change message.
// The error return is non-nil if the message is not a meta tempo-change.
//
func MetaTempoBPM(msg gomidi.Message) (tempo float64, err error) {
	var usec uint64
	usec, err = MetaTempoMicroseconds(msg)
	if err != nil {
		tempo = 60.0
		return tempo, err
	}
	return TEMPO_CONSTANT/float64(usec), err
}


// IsTempoChange returns true iff message is a meta temp-change.
//
func IsTempoChange(msg gomidi.Message) bool {
	d := msg.Data
	return len(d) == 6 && d[0] = 0xFF && d[1] == 0x51
}
	
