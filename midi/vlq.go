package midi

/*
 * Defines MIDI variable-length quantity.
 * http://midi.teragonaudio.com/tech/midifile/vari.htm
 *
 */

import "fmt"

// VLQ implements MIDI file variable length quantity.
//
type VLQ struct {
	bytes []byte
}

func reverse(src []byte) []byte {
	acc := make([]byte, len(src))
	for i, j := len(src)-1, 0; j < len(src); i, j = i-1, j+1 {
		acc[j] = src[i]
	}
	return acc
}

func (vlq *VLQ) setBytes(bytes []byte) {
	vlq.bytes = bytes
}

func (vlq *VLQ) Bytes() []byte {
	return vlq.bytes
}

func (vlq *VLQ) SetValue(n int) {
	mask := 0x7f
	acc := make([]byte, 0, 4)
	acc = append(acc, byte(n & mask))
	n = n >> 7
	for n > 0 {
		acc = append(acc, byte(0x80 |(n & mask)))
		n = n >> 7
	}
	vlq.bytes = make([]byte, len(acc))
	vlq.bytes = reverse(acc)
}
	
func (vlq *VLQ) Value() int {
	acc := 0
	scale := 1
	for _, b := range reverse(vlq.bytes) {
		acc = acc + scale * int(b & 0x7f)
		scale *= 128
	}
	return acc
}

func (vlq *VLQ) String() string {
	s := "VLQ [ "
	for _, b := range vlq.bytes {
		s += fmt.Sprintf("%02x ", b)
	}
	s += fmt.Sprintf("] Value 0x%x", vlq.Value())
	return s
}


func (vlq *VLQ) Length() int {
	return len(vlq.bytes)
}

func NewVLQ(value int) *VLQ {
	vlq := &VLQ{}
	vlq.SetValue(value)
	return vlq
}

