package midi

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
	
func KeyName(n int64) string {
	if 0 <= n && n < 128 {
		return keynames[n]
	} else {
		return "<ERR>"
	}
}
