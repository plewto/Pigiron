package midi

/*
** chunk.go defines a generalized MIDI file chunk structure.
** In practice there are two type of chunks:
**   1. Header  id = 'MThd'
**   2. Track   id = 'MTrk'
**
*/

import (
	"fmt"
	"os"
	"github.com/plewto/pigiron/pigerr"
)


// chunkID type represents a 4-byte chunk-type code.
//
type chunkID [4]byte

func (id *chunkID) String() string {
	s := "chunkID '%c%c%c%c'"
	return fmt.Sprintf(s, id[0], id[1], id[2], id[3])
}

// chunkID.eq returns true if two chunks are equivalent.
//
func (id chunkID) eq(other chunkID) bool {
	for i := 0; i < len(id); i++ {
		if id[i] != other[i] {
			return false
		}
	}
	return true
}


// Chunk interface defines common methods for MIDI file chunks.
//
type Chunk interface {
	ID() chunkID
	Length() int
}

// readChunkPreamble(f *osFile) reads the next 8-bytes from an open file as the start of a chunk.
//
// Returns:
//    id     - 4-byte chunkID
//    length - number of remaining bytes in the chunk
//    error  - no-nil if the data dose not look like the start of a chunk.
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

// readRawChunk(f *osFile) reads chuck data from open file.
// The file's read pointer should be positioned at the start of the chunk's 
// 4-byte id.
//
// Returns:
//     id    - 4-byte chunk id
//     data  - chucks contents
//     error - non-nil if the chunk could not be read.
//
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
		errmsg := "smf.readRawChunk read value count inconsistent.\n"
		errmsg += "Expected %d bytes, read %d"
		err = pigerr.New(fmt.Sprintf(errmsg, length, count))
		return
	}
	return
}


// msb(n) function returns upper byte of 16-bit value.
//
func msb(n int) byte {
	hi := (n & 0xFF00) >> 8
	return byte(hi)
}

// lsb(n) function returns the lower byte of 16-bit value.
//
func lsb(n int) byte {
	return byte(n & 0x00FF)
}

