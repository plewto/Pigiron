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

func (m *MetaMessage) ToPortmidiEvent() (portmidi.Event, error) {
	var pme portmidi.Event
	msg := "MetaMessage can not be converted to portmidi.Event. "
	msg += "MetaType is 0x%x '%s'"
	err := fmt.Errorf(msg, byte(m.MetaType()), m.MetaType())
	return pme, err
}
	


func newMetaMessage(bytes []byte) (*MetaMessage, error) {
	var err error
	var meta *MetaMessage
	if len(bytes) < 2 {
		msg := "smf.newMetaMessage, bytes slice too small"
		err = fmt.Errorf(msg)
		return meta, err
	}
	st := bytes[0]
	mtype := bytes[1]
	if !isMetaStatus(st) || !isMetaType(mtype) {
		msg := "smf.newMetaMessage incorrect status 0x%x or type 0x%x bytes"
		err := fmt.Errorf(msg, st, mtype)
		return meta, err
	}
	// ISSUE: TODO validate data count
	meta = &MetaMessage{bytes}
	return meta, err
}
