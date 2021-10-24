package midi

import (
	"fmt"
	"testing"
	gomidi "gitlab.com/gomidi/midi/v2"
)


func TestTempo(t *testing.T) {
	var t60 gomidi.Message
	var err error
	var bpm float64
	t60, err = MakeTempoMessage(60.0)
	if err != nil {
		t.Fatalf(err.Error())
	}
	bpm, err = MetaTempoBPM(t60)
	if err != nil {
		t.Fatalf(err.Error())
	}
	if bpm != 60.0 {
		errmsg := "Expected tempo of 60 BPM, got %f"
		err = fmt.Errorf(errmsg, bpm)
		t.Fatalf(err.Error())
	}
	_, err = MakeTempoMessage(-1.0)
	if err == nil {
		errmsg := "Did not detect negative tempo"
		t.Fatalf(errmsg)
	}
	_, err = MakeTempoMessage(MAX_TEMPO+1)
	if err == nil {
		errmsg := "Did not detect out of bounds tempo"
		t.Fatalf(errmsg)
	}
}

