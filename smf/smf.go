package smf

import (
	"fmt"
	"os"
)

// chunkID indicates the chunk type as 4 bytes
//
type chunkID [4]byte

var (
	headerID chunkID = [4]byte{0x4d, 0x54, 0x68, 0x64}
	trackID chunkID = [4]byte{0x4d, 0x54, 0x72, 0x6B}
)

func (id *chunkID) String() string {
	acc := fmt.Sprintf("chunkID '%c%c%c%c'", id[0], id[1], id[2], id[3])
	return acc
}


// Chunk interface defines general methods of SMF chunk types.
// Currently there are only two possible chuck types: Header and Track
//
type Chunk interface {
	ID() chunkID
	Length() int
	Bytes() []byte
	Dump()
}

// SMF struct Represents a Standard MIDI File.
//
type SMF struct {
	filename string
	header *Header
	chunks []Chunk
}


// readChunkPreamble reads chunk ID and data length from current file position.
// 
func readChunkPreamble(f *os.File) (id [4]byte, length int, err error) {
	var buffer = make([]byte, 8)
	var n = 0
	n, err = f.Read(buffer)
	if n != 8 && err == nil {
		err = fmt.Errorf("ERROR: smf.readChunkPreamble, file did not contain expected number of bytes.")
		return id, length, err
	}
	if err != nil {
		return id, length, err
	}
	for i:=0; i<4; i++ {
		id[i] = buffer[i]
	}
	length, err = getLong(buffer, 4)
	return id, length, err
}



// ReadSMF opens MIDI file for input.
//
func ReadSMF(filename string) (*SMF, error) {
	var err error
	var smf = &SMF{}
	smf.filename = filename
	var header *Header
	var file *os.File
	file, err = os.Open(filename)
	if err != nil {
		return smf, err
	}
	defer file.Close()
	header, err = readHeader(file)
	if err != nil {
		return smf, err
	}
	smf.header = header
	for chunk := 0; chunk < header.chunkCount; chunk++ {
		var id [4]byte
		var length int
		id, length, err = readChunkPreamble(file)
		if err != nil {
			return smf, err
		}
		if id != trackID { // skip non-track chunks
			var buffer = make([]byte, length)
			_, err = file.Read(buffer)
			if err != nil {
				return smf, err
			}
		} else { // read track data
			trk := &Track{make([]byte, 0, 1024)}
			trk.bytes, err = readTrackBytes(file, length)
			if err != nil {
				return smf, err
			}
			smf.chunks = append(smf.chunks, trk)
		}
	}
	return smf, err
}


func (smf *SMF) Dump() {
	fmt.Println("SMF")
	fmt.Printf("\tfilename = \"%s\"\n", smf.filename)
	if smf.header != nil {
		smf.header.Dump()
	} else {
		fmt.Println("Header is nil")
		return
	}
	for i:=0; i<smf.header.chunkCount; i++ {
		cnk := smf.chunks[i]
		fmt.Printf("--------------- chunk %d\n", i)
		cnk.Dump()
	}

}


