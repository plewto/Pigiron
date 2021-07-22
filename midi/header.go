
package midi

import (
	"fmt"
	"os"
	"github.com/plewto/pigiron/pigerr"
)

var (
	headerID chunkID = [4]byte{0x4d, 0x54, 0x68, 0x64}
)

type SMFHeader struct {
	format int
	trackCount int  // NOTE: trackCount may be greater then actual track count.
	division int
}

func (h *SMFHeader) ID() chunkID {
	return headerID
}

func (h *SMFHeader) Division() int {
	return h.division
}

func (h *SMFHeader) Length() int {
	return 6
}

func (h *SMFHeader) Bytes() []byte {
	var acc = make([]byte, 14)
	for i, b := range h.ID() {
		acc[i] = byte(b)
	}
	for i := 4; i < 7; i++ {
		acc[i] = 0
	}
	acc[7] = 6
	acc[ 8], acc[ 9] = msb(h.format), lsb(h.format)
	acc[10], acc[11] = msb(h.trackCount), lsb(h.trackCount)
	acc[12], acc[13] = msb(h.division), lsb(h.division)
	return acc
}



func (h *SMFHeader) Dump() {
	bytes := h.Bytes()
	var dumpLine = func(index int, count int, tag string) {
		fmt.Printf("  [%3d] %-8s : ", index, tag)
		for i, j := index, count; i < len(bytes) && j > 0; i, j = i+1, j-1 {
			fmt.Printf("%02X ", bytes[i])
		}
		fmt.Println()
	}
	fmt.Println("SMFHeader")
	dumpLine( 0, 4, "ID")
	dumpLine( 4, 4, "Length")
	dumpLine( 8, 2, "Format")
	dumpLine(10, 2, "Tracks")
	dumpLine(12, 2, "Division")
}


	
func readSMFHeader(f *os.File) (header *SMFHeader, err error) {
	var id chunkID
	var length int
	id, length, err = readChunkPreamble(f)
	if err != nil {
		return
	}
	if !id.eq(headerID) {
		errmsg := "Expected header chunk id '%s', got '%s'"
		err = pigerr.New(fmt.Sprintf(errmsg, headerID, id))
		return
	}
	var data = make([]byte, length)
	var count = 0
	count, err = f.Read(data)
	if count != length {
		errmsg := "smf.readRawChunk read value count inconsistent.\n"
		errmsg += "Expected %d bytes, read %d"
		err = pigerr.New(fmt.Sprintf(errmsg, length, count))
		return
	}
	// DO NOT replace above lines with readRawChunk()
	// It may not detect non-smf files and may attempt to 
	// load huge amounts of data, freezing the application.
	// id, data, err = readRawChunk(f)
	//
	if err != nil {
		errmsg := "readSMFHeader Could not read SMF header chunk"
		err = pigerr.CompoundError(err, errmsg)
		return
	}
	if !id.eq(headerID) {
		errmsg := "readSMFHeader encountered wrong chunk id, expected %s, got %s"
		err = pigerr.New(fmt.Sprintf(errmsg, headerID, id))
		return
	}
	if len(data) != 6 {
		errmsg := "readSMFHeader expected header chunk length of 6 byte, got %d"
		err = pigerr.New(fmt.Sprintf(errmsg, len(data)))
		return
	}
	var format, trackCount, division int
	format, data, _ = takeShort(data)
	trackCount, data, _ = takeShort(data)
	division, _, _ = takeShort(data)
	header = &SMFHeader{format, trackCount, division}
	if format < 0 || 2 < format {
		dflt := 0
		errmsg := "MIDI file has unsupported format value: %d, using default %d"
		pigerr.Warning(fmt.Sprintf(errmsg, format, dflt))
		header.format = dflt
	}
	if division < 24 || 960 < division {
		dflt := 24
		errmsg1 := "MIDI file has out of bounds clock division"
		errmsg2 := fmt.Sprintf("Expected division between 24 and 960, got %d\n", division)
		errmsg3 := fmt.Sprintf("Using default %d\n", dflt)
		pigerr.Warning(errmsg1, errmsg2, errmsg3)
		header.division = dflt
	}
	return
}	
	
	
	
