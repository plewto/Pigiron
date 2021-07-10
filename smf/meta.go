package smf


/*
 * Defines MetaMesasge implementation of MIDIMessage
 *
 */

import (
	"fmt"
	"github.com/rakyll/portmidi"
)


// MetaType defines valid meta message type values.
//
type MetaType byte

const (
 	MetaSequenceNumber MetaType = 0x00  
 	MetaText MetaType = 0x01
 	MetaCopyright MetaType = 0x02
 	MetaTrackName MetaType = 0x03
 	MetaInstrumentName MetaType = 0x04
 	MetaLyric MetaType = 0x05
 	MetaMarker MetaType = 0x06
 	MetaCuePoint MetaType = 0x07
 	MetaChannelPrefix MetaType = 0x20
 	MetaEndOfTrack MetaType = 0x2F
 	MetaTempo MetaType = 0x51
 	MetaSMPTE MetaType = 0x54
 	MetaTimeSignature MetaType = 0x58
 	MetaKeySignature MetaType = 0x59    
 	MetaSequencerEvent MetaType = 0x7f  
)

var (
	metaTextTypes map[MetaType]bool = map[MetaType]bool{
		MetaText: true,
		MetaCopyright: true,
		MetaTrackName: true,
		MetaInstrumentName: true,
		MetaLyric: true,
		MetaMarker: true,
		MetaCuePoint: true,
	}

	// metaTypeTable maps MetaType to string mnemonic.
	//
	metaTypeTable map[MetaType]string = map[MetaType]string {
		MetaSequenceNumber: "SEQ-NUMBER",
		MetaText: "TEXT",
		MetaCopyright: "COPYRIGHT",
		MetaTrackName: "TRACK-NAME",
		MetaInstrumentName: "INSTRUMENT-NAME",
		MetaLyric: "LYRIC",
		MetaMarker: "MARKER",
		MetaCuePoint: "CUE",
		MetaChannelPrefix: "CHANNEL-PREFIX",
		MetaEndOfTrack: "END-OF-TRACK",
		MetaTempo: "TEMPO",
		MetaSMPTE: "SMPTE",
		MetaTimeSignature: "TIME-SIG",
		MetaKeySignature: "KEY-SIG",
		MetaSequencerEvent: "SEQ-EVENT",
	}
)


func (t MetaType) String() string {
	s, flag := metaTypeTable[t]
	if !flag {
		s = "?"
	}
	return s
}

// isMetaType returns true if n is a valid meta type.
//
func isMetaType(n byte) bool {
	_, flag := metaTypeTable[MetaType(n)]
	return flag
}

// isMetaTextType returns true if n is one of the text based meta types.
//
func isMetaTextType(n byte) bool {
	_, flag := metaTextTypes[MetaType(n)]
	return flag
}


// Implements MIDIMessage
//
type MetaMessage struct {
	bytes []byte
}

// Status returns meta status byte (always 0xFF).
//
func (m *MetaMessage) Status() StatusByte {
	return StatusByte(m.bytes[0])
}

// Bytes returns the byte sequence corresponding to this MetaMessage.
//
func (m *MetaMessage) Bytes() []byte {
	return m.bytes
}

func (m *MetaMessage) String() string {
	if len(m.bytes) < 2 {
		return "MetaMessage ??"
	} else {
		return fmt.Sprintf("MetaMessage %s", MetaType(m.bytes[1]))
	}
}

func (m *MetaMessage) Dump() {
	if m == nil {
		fmt.Println("MetaMessage <nil>")
		return
	}
	fmt.Printf("%s\n", m)
	fmt.Printf("[ 0] 0x%02x   - status\n", m.bytes[0])
	fmt.Printf("[ 1] 0x%02x   - meta type '%s'\n", m.bytes[1], MetaType(m.bytes[1]))
	counter := 0
	for i, b := range m.bytes[2:] {
		fmt.Printf("[%2d] 0x%02x   - vlq-%d\n", i+2, b, i)
		counter++
		if b & 0x80 == 0 {
			break
		}
	}
	isText := isMetaTextType(m.bytes[1])
	start := counter + 2
	for i, b := range m.bytes[start:] {
		fmt.Printf("[%2d] 0x%02x   - data-%d", i+start, b, i)
		if isText {
			fmt.Printf(" '%c'", b)
		}
		fmt.Println()
	}
}
			
// MetaType returns byte indicating the type of this MetaMessage.
//
func (m *MetaMessage) MetaType() MetaType {
	return MetaType(m.bytes[1])
}

// Data() returns only the data bytes for this MetaMessage.
// The result is a sub-sequence of the Bytes() value and excludes the 
// status, meta-type and vlq length bytes.
//
func (m *MetaMessage) Data() []byte {
	index := 2
	bytes := m.bytes
	for index < len(bytes) {
		b := bytes[index]
		index++
		if b & 0x80 == 0 {
			break
		}
	}
	return m.bytes[index:]
}

// ToPortmidiEvent for MetaMessage is included to satisfy the MIDIMessage interface.
// It is not possible to convert meta messages to portmidi events.
// Always returns a non-nil error.
//
func (m *MetaMessage) ToPortmidiEvent() (portmidi.Event, error) {
	var dummy portmidi.Event
	msg := "MetaMessage %s can not be converted to portmidi.Event"
	err := fmt.Errorf(msg, m.MetaType())
	return dummy, err
}

// newMetaMessage returns pointer to new instance of MetaMessage
// mtype - must be a valid MetaType
// data - the data bytes.
// status, meta-type and vlq count bytes are included automatically.
//
func newMetaMessage(mtype MetaType, data []byte) *MetaMessage {
	head := []byte{byte(MetaStatus), byte(mtype)}
	vlq := NewVLQ(len(data))
	bytes := append(head, vlq.Bytes()...)
	mm := &MetaMessage{append(bytes, data...)}
	return mm
}

// NewMetaSequenceNumber returns pointer to new MetaMessage
// n - the sequence number
//
func NewMetaSequenceNumber(n byte) *MetaMessage {
	data := []byte{n}
	return newMetaMessage(MetaSequenceNumber, data)
}

// NewMetaText returns pointer to new MetaMessage for holding text.
// mtype may be any text-based MetaType (Text, Copyright, TrackName,
// InstrumentName, Lyric, Marker, or Cue).  Non-valid MetaTypes are
// silently converted to MetaText.
//
func NewMetaText(mtype MetaType, text string) *MetaMessage {
	if !isMetaTextType(byte(mtype)) {
		mtype = MetaText
	}
	data := []byte(text)
	return newMetaMessage(mtype, data)
}


// NewMetaChannelPrefix returns pointer to MetaMessage of type ChannelPrefix
// prefix - the channel prefix
//
func NewMetaChannelPrefix(prefix byte) *MetaMessage {
	data := []byte{prefix}
	return newMetaMessage(MetaChannelPrefix, data)
}

var (
	eot = newMetaMessage(MetaEndOfTrack, []byte{})
)

// NewMetaEndOfTrack returns pointer to MetaMessage of type EndOfTrack.
// This function always returns the same object.
//
func NewMetaEndOfTrack() *MetaMessage {
	return eot
}

// NewMetaTempo returns pointer to MetaMessage of type MetaTempo.
// The 4-byte argument is the tempo is microseconds per quarter note.
//
func NewMetaTempo(usec [4]byte) *MetaMessage {
	return newMetaMessage(MetaTempo, usec[:])
}

// NewMetaSMPTE returns pointer to MetaMessage of type MetaSMPTE.
// SMPTE type is included for completeness and not otherwise supported.
func NewMetaSMPTE(values [5]byte) *MetaMessage {
	return newMetaMessage(MetaSMPTE, values[:])
}

// NewMetaTimeSignature returns pointer to MetaMessage of type MetaTimesignature
// num - beats per measure
// dnom - beat value 2->half, 4->quarter 8->eighth etc...
// optional args 1 - MIDI clocks per metronome click, default 24
// optional args 2 - notated 32nd notes per quarter note, default 8
//
func NewMetaTimeSignature(num byte, dnom byte, args... byte) *MetaMessage {
	cpc, n32 := byte(24), byte(8)
	if len(args) > 0 {
		cpc = args[0]
	}
	if len(args) > 1 {
		n32 = args[1]
	}
	return newMetaMessage(MetaTimeSignature, []byte{num, dnom, cpc, n32})
}

// NewMetaKeySignature returns pointer to MetaMessage of type MetaKeySignature.
// key - the number of flats or sharps.  ISSUE: How are these represented.
// minor - bool, if true key is minor, otherwise it is major.
//
// KeySignature is included for completeness and not otherwise supported.
//
func NewMetaKeySignature(key byte, minor bool) *MetaMessage {
	var flag = byte(0)
	if minor {
		flag = 1
	}
	return newMetaMessage(MetaKeySignature, []byte{key, flag})
}

// NewSequencerSpecifcMessage returns pointer to MetaMessage
// SequencerSpecific is included for completeness and not otherwise supported.
//
func NewSequencerSpecificMessage(data []byte) *MetaMessage {
	return newMetaMessage(MetaSequencerEvent, data)
}

