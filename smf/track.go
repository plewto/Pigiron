package smf

import (
	"os"
	"fmt"
)

// Track implements MIDI file track chunk
//
type Track struct {
	bytes []byte
}

func NewTrack() *Track {
	trk := &Track{make([]byte, 0)}
	return trk
}

func (t *Track) ID() chunkID {
	return trackID
}

func (t *Track) Length() int {
	return len(t.bytes)
}

func (t *Track) Bytes() []byte {
	return t.bytes
}

func (t *Track) Dump() {
	lineLength := 16

	xline := func(index int) (string, string) {
		acc := fmt.Sprintf("\t[%4x] ", index) // 9
		bcc := ""
		for i, j := index, 0; i < len(t.bytes) && j < lineLength; i, j = i+1, j+1 {
			b := t.bytes[i]
			acc += fmt.Sprintf("%02x ", b)  // 12 * 3
			if b < 32 || b > 127 {
				bcc = bcc + "-"
			} else {
				bcc = bcc + fmt.Sprintf("%c", b)
			}
		}
		return acc, bcc
	}

	pad := func(s string) string {
		width := 9 + lineLength * 3
		for len(s) < width {
			s = s + " "
		}
		return s + " : "
	}
			
	fmt.Printf("Chunk %s\n", t.ID())
	fmt.Printf("\tByte count = %d\n", t.Length())
	for i:=0; i < t.Length(); i += lineLength {
		acc, bcc := xline(i)
		fmt.Print(pad(acc))
		fmt.Println(bcc)
	}
}


// readTrackBytes reads length bytes from file.
// The file pointer is expected to point to the start of track data,
// immediately following the chunk length bytes.
// 
func readTrackBytes(f *os.File, length int) ([]byte, error) {
	var err error
	var buffer = make([]byte, length)
	var n = 0
	n, err = f.Read(buffer)
	if n != length {
		msg := "smf.readTrackBytes, file did not contain expected number of bytes: %d"
		err := exError(fmt.Sprintf(msg, length))
		return buffer, err
	}
	return buffer, err
}
