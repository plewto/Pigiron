package smf

/*
** header.go defines MIDI file header chunk.
**
*/

import (
	"fmt"
	"os"
	"github.com/plewto/pigiron/pigerr"
	"github.com/plewto/pigiron/expect"
)

var headerID chunkID = [4]byte{0x4d, 0x54, 0x68, 0x64}


// SMFheadr strut implements Chunk interface for MIDI file headers.
//
type Header struct {
	format int
	trackCount int  // NOTE: trackCount may be greater then actual track count.
	division int
}

func (h *Header) String() string {
	msg := "smf.Header format: %d  trackCount: %d division: %d"
	return fmt.Sprintf(msg, h.format, h.trackCount, h.division)
}

func (h *Header) ID() chunkID {
	return headerID
}

// h.Format() return MIDI file format
//
func (h *Header) Format() int {
	return h.format
}


// h.Division() returns MIDI file clock division.
//
func (h *Header) Division() int {
	return h.division
}

func (h *Header) Length() int {
	return 6
}

// h.Dump() displays contents of MIDI header chunk.
//
func (h *Header) Dump() {
	fmt.Println("Header:")
	fmt.Printf("\tformat     : %4d\n", h.format)
	fmt.Printf("\tchuckCount : %4d\n", h.trackCount)
	fmt.Printf("\tdivision   : %4d\n", h.division)
}


// readHeader function reads MIDI file header chuck from file.
//
func readHeader(f *os.File) (header *Header, err error) {
	var id chunkID
	var length int
	id, length, err = readChunkPreamble(f)
	if err != nil {
		return
	}
	if !id.eq(headerID) {
		msg := "Expected header id '%s', got '%s'"
		err = fmt.Errorf(msg, headerID, id)
		return
	}
	if length != 6 {
		msg := "Unusual SMF header length, expected 6, got %d\n"
		err = fmt.Errorf(msg, length)
		return
	}

	var data = make([]byte, length)
	var count = 0
	count, err = f.Read(data)
	if count != length {
		msg := "SMF Header data count inconsistent, expected %d bytes, read %d"
		err = fmt.Errorf(msg, count, length)
		return
	}
	if err != nil {
		msg := "smf.readHeader could not read Header chunk\n"
		msg += fmt.Sprintf("%s", err)
		err = fmt.Errorf(msg)
		return
	}
	// DO NOT replace above lines with readRawChunk()
	//    It may not detect non-smf files and attmpt to read
	//    huge amountrs of data.
	//

	var format, trackCount, division int
	format, data, _ = expect.TakeShort(data)
	trackCount, data, _ = expect.TakeShort(data)
	division, _, _ = expect.TakeShort(data)
	header = &Header{format, trackCount, division}
	if format < 0 || 2 < format {
		dflt := 0
		errmsg := "MIDI file has unsupported format: %d, using default %d"
		pigerr.Warning(fmt.Sprintf(errmsg, format, dflt))
		header.format = dflt
	}
	if division < 24 || 960 < division {
		dflt := 24
		msg1 := "MIDI file has out of bounds clock division"
		msg2 := fmt.Sprintf("Expected division between 24 and 960, got %d", division)
		msg3 := fmt.Sprintf("Using default %d", dflt)
		pigerr.Warning(msg1, msg2, msg3)
		header.division = dflt
	}
	return
}
		
	


