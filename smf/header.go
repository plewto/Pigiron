package smf

import (
	"fmt"
	"os"
)

// Header implements MIDI file header Chunk.
//
type Header struct {
	format int
	chunkCount int
	division int
}

func (h *Header) ID() chunkID {
	return headerID
}

func (h *Header) Length() int {
	return 6
}

func (h *Header) Bytes() []byte {
	lsb := func(n int) byte {
		return byte(n & 0xFF)
	}
	msb := func(n int) byte {
		return byte(n >> 8)
	}
	acc := make([]byte, 14)
	for i, b := range headerID {
		acc[i] = b
	}
	for i, b := range []byte{0, 0, 0, 6} {
		acc[i+4] = b
	}
	acc[ 8], acc[ 9] = msb(h.format), lsb(h.format)
	acc[10], acc[11] = msb(h.chunkCount), lsb(h.chunkCount)
	acc[12], acc[13] = msb(h.division), lsb(h.division)
	return acc
}

// readHeader constructs Header from values at current file position.
//
func readHeader(f *os.File) (*Header, error) {
	var err error
	var n int
	var buffer = make([]byte, 14)
	var header *Header
	n, err = f.Read(buffer)
	if err != nil {
		return header, err
	}
	if n != 14 {
		msg := "Expected smf Header of length 14 bytes, got %d"
		err = exError(fmt.Sprintf(msg, n))
		return header, err
	}
	err = expectChunkID(buffer, 0, headerID)
	if err !=  nil {
		return header, err
	}
	header = &Header{}
	header.format, _ = getShort(buffer, 8)        // TODO validate
	header.chunkCount, _ = getShort(buffer, 10)   // TODO validate
	header.division, _ = getShort(buffer, 12)     // TODO validate
	return header, err
}


func (h *Header) Dump() {
	msb := func(n int) byte {
		return byte(n >> 8) & 0x7F
	}
	lsb := func(n int) byte {
		return byte(n & 0x7f)
	}
	fmt.Println("smf.Header")
	fmt.Printf("\t[ 0] id     : ")
	acc := ""
	for _, b := range headerID {
		fmt.Printf("%02x ", b)
		acc += fmt.Sprintf("%c", b)
	}
	fmt.Printf(" -> %s\n", acc)
	fmt.Printf("\t[ 4] Length : 00 00 00 06  -> 6\n")
	n := h.format
	fmt.Printf("\t[ 8] Format : %02x %02x        -> %d\n", msb(n), lsb(n), n)
	n = h.chunkCount
	fmt.Printf("\t[10] Tracks : %02x %02x        -> %d\n", msb(n), lsb(n), n)
	n = h.division
	fmt.Printf("\t[12] Div    : %02x %02x        -> %d\n", msb(n), lsb(n), n)
}
	
