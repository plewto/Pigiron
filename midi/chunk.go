package midi

import (
	"fmt"
	"os"
	"github.com/plewto/pigiron/pigerr"
)


// chunkID type represents a 4-byte type code.
//
type chunkID [4]byte

func (id *chunkID) String() string {
	s := "chunkID '%c%c%c%c'"
	return fmt.Sprintf(s, id[0], id[1], id[2], id[3])
}

// chunkID.eq returns true if two chunks are equivilent.
//
func (id chunkID) eq(other chunkID) bool {
	for i := 0; i < len(id); i++ {
		if id[i] != other[i] {
			return false
		}
	}
	return true
}


// Chunk interface defines common methods for MIDI file structures.
// Two MIDI file chunk types are defined:
//   1) header   id 'MThd'
//   2) track    id 'MTrk
//
type Chunk interface {
	ID() chunkID
	Length() int
	// Bytes() []byte
	// Dump()
}


// readChunkPreamble reads the next 8-bytes from an opern file as the start of a chunk.
// Returns:
//    id  - 4-byte chunkID
//    length - number of remaining bytes in the chunk
//    error - 
//    
func readChunkPreamble(f *os.File) (id chunkID, length int, err error) {
	var buffer = make([]byte, 8)
	var n = 0
	n, err = f.Read(buffer)
	if n != 8 && err == nil {
		errmsg := "smf.readChunkPreamble, file does not contain minimal number of byte."
		err = pigerr.New(errmsg)
		return id, length, err
	}
	if err != nil {
		err = pigerr.CompoundError(err, "smf.readChunkPreamble, can not read chunk")
		return id, length, err
	}
	for i := 0; i < 4; i++ {
		id[i] = buffer[i]
	}
	length, _, err = takeLong(buffer[4:])
	return id, length, err
}


func readRawChunk(f *os.File) (id chunkID, data []byte, err error) {
	var length int
	id, length, err = readChunkPreamble(f)
	if err != nil {
		return
	}
	data = make([]byte, length)
	var count int
	count, err = f.Read(data)
	if err != nil {
		errmsg := "smf.readRawChunk could not read chunk values."
		err = pigerr.CompoundError(err, errmsg)
		return
	}
	if count != length {
		errmsg := "smf.readRawChunk read value count inconsistenet.\n"
		errmsg += "Expected %d bytes, read %d"
		err = pigerr.New(fmt.Sprintf(errmsg, length, count))
		return
	}
	return
}



func msb(n int) byte {
	hi := (n & 0xFF00) >> 8
	return byte(hi)
}

func lsb(n int) byte {
	return byte(n & 0x00FF)
}

