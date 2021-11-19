package midi

/*
** keynum.go defines string representations for MIDI key number.
**
*/

import (
	"fmt"
)

var keynames [128]string

func init() {
	base := [12]string {
		"C", "C#", "D", "D#", "E", "F",
		"F#", "G", "G#", "A", "A#", "B"}
	for i:=0; i<128; i++ {
		pc := i % 12
		oct := i / 12
		key := base[pc] + fmt.Sprintf("%d", oct)
		keynames[i] = key
		key = fmt.Sprintf("%4s", key)
	}
	
}


// KeyName() returns a string representation for a MIDI key number.
//
func KeyName(n byte) string {
	if 0 <= n && n < 128 {
		return keynames[n]
	} else {
		return "<ERR>"
	}
}
