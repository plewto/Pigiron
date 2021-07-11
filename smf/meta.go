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
	mtype MetaType
	data []byte
}

// Status returns meta status byte (always 0xFF).
//
func (m *MetaMessage) Status() StatusByte {
	return MetaStatus
}

// Bytes returns the byte sequence corresponding to this MetaMessage.
//
func (m *MetaMessage) Bytes() []byte { 
	prefix := []byte{byte(MetaStatus), byte(m.mtype)}
	vlq := NewVLQ(len(m.data))
	acc := append(prefix, vlq.Bytes()...)
	acc = append(acc, m.data...)
	return acc
}

func (m *MetaMessage) String() string {
	return fmt.Sprintf("MetaMessage %s", m.mtype)
}

func (m *MetaMessage) Dump() {
	fmt.Printf("MetaMessage '%s'\n", m.mtype)
	fmt.Printf("[ 0] 0x%02x    - status", byte(MetaStatus))
	fmt.Printf("[ 1] 0x%02x    - type '%s'\n", byte(m.mtype), m.mtype)
	bytes := m.Bytes()
	counter := 0
	for i, b := range bytes[2:] {
		fmt.Printf("[%2d] 0x%02x   - VLQ-%d\n", i+2, b, counter)
		counter++
		if b & 0x80 == 0 {
			break
		}
	}
	offset := counter + 2
	counter = 0
	mtype := byte(m.MetaType())
	for i, b := range bytes[offset:] {
		fmt.Printf("[%2d] 0x%02x   - Data-%d", offset + i, b, counter)
		if isMetaTextType(mtype) {
			fmt.Printf(" '%c'", b)
		}
		fmt.Println()
		counter++
	}
}
			
		

// MetaType returns byte indicating the type of this MetaMessage.
//
func (m *MetaMessage) MetaType() MetaType {
	return m.mtype
}


func (m *MetaMessage) Data() []byte {
	return m.data
}

func (m *MetaMessage) ToPortmidiEvent() (portmidi.Event, error) {
	var pme portmidi.Event
	msg := "MetaMessage can not be converted to portmidi.Event. "
	msg += "MetaType is 0x%x '%s'"
	err := fmt.Errorf(msg, byte(m.MetaType()), m.MetaType())
	return pme, err
}

func newMetaMessage(mtype MetaType, data []byte)(*MetaMessage, error) {
	var err error
	var meta *MetaMessage
	if !isMetaType(byte(mtype)) {
		msg := "smf.newMetaMessage, invalid MetaType 0x%x"
		err = fmt.Errorf(msg, byte(mtype))
		return meta, err
	}
	// ISSUE: TODO validate data
	meta = &MetaMessage{mtype, data}
	return meta, err
}
