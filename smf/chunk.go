package smf

import (
	"fmt"
	"os"
)

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


// readChunkPreamble reads chunk ID and data length from current file position.
// 
func readChunkPreamble(f *os.File) (id [4]byte, length int, err error) {
	var buffer = make([]byte, 8)
	var n = 0
	n, err = f.Read(buffer)
	if n != 8 && err == nil {
		msg := "smf.readChunkPreamble, file does not contain minimal number of bytes."
		err = fmt.Errorf(msg)
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
