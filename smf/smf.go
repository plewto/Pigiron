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

// expectID checks byte buffer for specific chunk ID.
// buffer - the data
// index - location in buffer to be checked.
// target - expect ID pattern
// Returns non-nil error if expected ID not found.
//
func expectID(buffer []byte, index int, target [4]byte) error {
	var err error
	for i, j := index, 0; j < 4; i, j = i+1, j+1 {
		if i > len(buffer) {
			err = fmt.Errorf("ERROR: smf.expectID,  index out of bounds, index = %d", i)
			return err
		}
		if buffer[i] != target[j] {
			id := ""
			for _, c := range target {
				id += fmt.Sprintf("%c", c)
			}
			err = fmt.Errorf("ERROR: smf.expectID, expected chunk id '%v' not found", id)
			return err
		}
	}
	return err
}

// getLong extracts 4-byte value from buffer starting at index.
//
func getLong(buffer []byte, index int) (int, error) {
	var err error
	if len(buffer) < index+4 {
		msg := "ERROR smf.getLong() index out of range: index = %d, buffer length = %d"
		err = fmt.Errorf(msg, index, len(buffer))
		return 0, err
	}
	acc := 0
	for i, j, shift := index, 0, 24; j < 4; i, j, shift = i+1, j+1, shift-8 {
		n := int(buffer[i])
		acc += int(n << shift)
	}
	return acc, err
}

// getShort extracts 2-byte value from buffer starting at index.
//
func getShort(buffer []byte, index int) (int, error) {
	var err error
	if len(buffer) < index+2 {
		msg := "ERROR smf.getShort() index out of range: index = %d, buffer length = %d"
		err = fmt.Errorf(msg, index, len(buffer))
		return 0, err
	}
	acc := 0
	for i, j, shift := index, 0, 8; j < 2; i, j, shift = i+1, j+1, shift-8 {
		n := int(buffer[i])
		acc += int(n << shift)
	}
	return acc, err
}

// getByte extracts byte from buffer at index.
//
func getByte(buffer []byte, index int) (int, error) {
	var err error
	if len(buffer) <= index {
		msg := "ERROR smf.getByte() index out of range: index = %d, buffer length = %d"
		err = fmt.Errorf(msg, index, len(buffer))
		return 0, err
	}
	return int(buffer[index]), err
}

// getVLQ extracts variable-length-value starting at index.
// Between 1 and 4 bytes are used.
//
func getVLQ(buffer []byte, index int)(*VLQ, error) {
	var err error
	var vlq *VLQ = new(VLQ)
	var maxCount = 4
	var acc = make([]byte, 0, maxCount)
	for {
		if index == maxCount {
			break
		}
		if index >= len(buffer) {
			msg := "ERROR smf.getVLQ index out of bounds, index = %d, buffer length = %d"
			err = fmt.Errorf(msg, index, len(buffer))
			return vlq, err
		}
		n := buffer[index]
		acc = append(acc, n)
		if n & 0x80 == 0 {
			break
		}
		index++
	}
	vlq.SetBytes(acc)
	return vlq, err
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


